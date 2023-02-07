package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api"
	"github.com/cucumber/messages-go/v16"
)

func bodyToReader(body interface{}) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}
	marshalled, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(marshalled), nil
}

func isAnErrorHandlingScenario(tags []*messages.PickleTag) bool {
	for _, nextTag := range tags {
		if nextTag.Name == "@Error" {
			return true
		}
	}
	return false
}

func handleHttpError(ctx context.Context, response *http.Response) (context.Context, error) {
	if response.StatusCode >= 400 {
		tags, ok := ctx.Value(TagsKey{}).([]*messages.PickleTag)
		if !ok {
			return ctx, fmt.Errorf("unable to get scenario tags")
		}

		if isAnErrorHandlingScenario(tags) {
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return ctx, fmt.Errorf("body read error: %v", err)
			}

			var httpErrorMessage api.HttpErrorMessage
			err = json.Unmarshal(body, &httpErrorMessage)
			if err != nil {
				return ctx, fmt.Errorf("body unmarshal error: %v (body: %v)", err, string(body))
			}
			return context.WithValue(ctx, ErrorKey{}, httpErrorMessage.Error), nil
		}
		return ctx, fmt.Errorf("wrong response code: %v", response.StatusCode)
	}
	return ctx, nil
}

func extractResponse(response *http.Response, responseDto interface{}) error {
	if responseDto == nil {
		return nil
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("wrong response code: %v", response.StatusCode)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("body read error: %v", err)
	}
	err = json.Unmarshal(body, &responseDto)
	if err != nil {
		return fmt.Errorf("body unmarshal error: %v (body: %v)", err, string(body))
	}
	return nil
}

func Call(ctx context.Context, url string, method string, body interface{}, responseDto interface{}) (context.Context, error) {
	bodyReader, err := bodyToReader(body)
	if err != nil {
		return ctx, fmt.Errorf("unable to create a reader from the body: %v", err)
	}
	request, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return ctx, fmt.Errorf("unable to create new http request: %v", err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return ctx, fmt.Errorf("unable to send request: %v", err)
	}
	defer response.Body.Close()

	ctx, err = handleHttpError(ctx, response)
	if err != nil {
		return ctx, fmt.Errorf("unable to handle request error: %v", err)
	}

	err = extractResponse(response, responseDto)
	if err != nil {
		return ctx, fmt.Errorf("unable to extract response: %v", err)
	}

	return ctx, nil
}
