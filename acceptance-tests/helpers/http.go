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

func handleHttpError(ctx context.Context, response *http.Response) (context.Context, bool, error) {
	if response.StatusCode >= 400 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return ctx, false, fmt.Errorf("body read error: %v", err)
		}
		if isAnErrorHandlingScenario(ctx) {
			var httpErrorMessage helpers.HttpErrorMessage
			err = json.Unmarshal(body, &httpErrorMessage)
			if err != nil {
				return ctx, false, fmt.Errorf("body unmarshal error: %v (body: %v)", err, string(body))
			}
			return context.WithValue(ctx, ErrorKey{}, httpErrorMessage.Error), true, nil
		}
		return ctx, false, fmt.Errorf("wrong response code: %v (%v)", response.StatusCode, string(body))
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

// func addUserHeader(ctx context.Context, request *http.Request) {
// 	currentUser, ok := ctx.Value(AuthentifiedUser{}).(application.UserID)
// 	if !ok {
// 		return
// 	}
// 	request.Header.Add(middlewares.UserIDHeader, currentUser.String())
// 	request.Header.Add("Authorization", "Bearer oui")
// }

func addCustomHeaders(ctx context.Context, request *http.Request, headers map[string]string) {
	for headerKey, headerValue := range headers {
		request.Header.Add(headerKey, headerValue)
	}
}

func addToken(ctx context.Context, request *http.Request, token string) {
	if token == "" {
		return
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))
}

type HttpRequest struct {
	Url         string
	Method      string
	Headers     map[string]string
	Body        interface{}
	ResponseDto interface{}
	Token       string
}

func DoCall(ctx context.Context, request *http.Request, responseDto interface{}) (context.Context, error) {
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
	addCustomHeaders(ctx, request, httpRequest.Headers)
	addToken(ctx, request, httpRequest.Token)
	return DoCall(ctx, request, httpRequest.ResponseDto)
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

	// TODO: AddToken ?
	// addUserHeader(ctx, request)
	return DoCall(ctx, request, responseDto)
}
