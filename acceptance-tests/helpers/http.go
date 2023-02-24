package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/api/middlewares"
	"github.com/ThomasFerro/l-edition-libre/application"
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

func handleHttpError(ctx context.Context, response *http.Response) (context.Context, bool, error) {
	if response.StatusCode >= 400 {
		tags, ok := ctx.Value(TagsKey{}).([]*messages.PickleTag)
		if !ok {
			return ctx, false, fmt.Errorf("unable to get scenario tags")
		}

		if isAnErrorHandlingScenario(tags) {
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return ctx, false, fmt.Errorf("body read error: %v", err)
			}

			var httpErrorMessage helpers.HttpErrorMessage
			err = json.Unmarshal(body, &httpErrorMessage)
			if err != nil {
				return ctx, false, fmt.Errorf("body unmarshal error: %v (body: %v)", err, string(body))
			}
			return context.WithValue(ctx, ErrorKey{}, httpErrorMessage.Error), true, nil
		}
		return ctx, false, fmt.Errorf("wrong response code: %v", response.StatusCode)
	}
	return ctx, false, nil
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

// TODO: Ne pas passer par un header mais par un JWT (OAuth2 ?)
func addUserHeader(ctx context.Context, request *http.Request) {
	currentUser, ok := ctx.Value(AuthentifiedUser{}).(application.UserID)
	if !ok {
		return
	}
	request.Header.Add(middlewares.UserIDHeader, currentUser.String())
}

func addCustomHeaders(ctx context.Context, request *http.Request, headers map[string]string) {
	for headerKey, headerValue := range headers {
		request.Header.Add(headerKey, headerValue)
	}
}

type HttpRequest struct {
	Url         string
	Method      string
	Headers     map[string]string
	Body        interface{}
	ResponseDto interface{}
}

func doCall(ctx context.Context, request *http.Request, responseDto interface{}) (context.Context, error) {
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return ctx, fmt.Errorf("unable to send request: %v", err)
	}
	defer response.Body.Close()

	ctx, handled, err := handleHttpError(ctx, response)
	if err != nil {
		return ctx, fmt.Errorf("unable to handle request error: %v", err)
	}
	if handled {
		return ctx, nil
	}

	err = extractResponse(response, responseDto)
	if err != nil {
		return ctx, fmt.Errorf("unable to extract response(%v): %v", response.StatusCode, err)
	}

	return ctx, nil
}

func Call(ctx context.Context, httpRequest HttpRequest) (context.Context, error) {
	bodyReader, err := bodyToReader(httpRequest.Body)
	if err != nil {
		return ctx, fmt.Errorf("unable to create a reader from the body: %v", err)
	}
	request, err := http.NewRequest(httpRequest.Method, httpRequest.Url, bodyReader)
	if err != nil {
		return ctx, fmt.Errorf("unable to create new http request: %v", err)
	}
	addUserHeader(ctx, request)
	addCustomHeaders(ctx, request, httpRequest.Headers)
	return doCall(ctx, request, httpRequest.ResponseDto)
}

func requestFromFormData(url string, filePath string, otherData map[string]string) (io.Reader, string, error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, "", err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return nil, "", err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, "", err
	}

	for dataKey, dataValue := range otherData {
		formField, err := writer.CreateFormField(dataKey)
		if err != nil {
			return nil, "", err
		}
		_, err = io.Copy(formField, strings.NewReader(dataValue))
		if err != nil {
			return nil, "", err
		}
	}

	return body, writer.FormDataContentType(), nil
}

func PostFile(ctx context.Context, url string, filePath string, otherData map[string]string, responseDto interface{}) (context.Context, error) {
	body, formDataContentType, err := requestFromFormData(url, filePath, otherData)
	if err != nil {
		return ctx, err
	}

	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", formDataContentType)

	addUserHeader(ctx, request)
	return doCall(ctx, request, responseDto)
}
