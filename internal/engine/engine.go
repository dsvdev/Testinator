package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dsvdev/Testinator/internal/model"
	llm_model "github.com/dsvdev/Testinator/internal/model/llm"
	"io"
	"log"
	"net/http"
	"reflect"
)

func NewEngine(baseURL string) *TestEngine {
	return &TestEngine{
		baseURL:    baseURL,
		httpClient: http.Client{},
	}
}

type TestEngine struct {
	httpClient http.Client
	baseURL    string
}

func (e *TestEngine) ExecuteLlmResponse(llmAction llm_model.LLMResponse, context model.TestExecutionContext) (model.TestExecutionContext, error) {
	log.Printf("Llm action: %v", llmAction)
	switch llmAction.Action {
	case llm_model.ActionHTTPRequest:
		var details llm_model.HTTPRequestDetails
		err := json.Unmarshal(llmAction.Details, &details)
		if err != nil {
			return nil, err
		}

		resp, err := e.executeHttpRequest(details)
		if err != nil {
			return nil, err
		}

		context = append(context, "Result - "+resp)
		return context, err
	default:
		return context, fmt.Errorf("unknown action: %s", llmAction.Action)
	}
}

func (e *TestEngine) executeHttpRequest(details llm_model.HTTPRequestDetails) (string, error) {
	// Создаем тело запроса, если есть
	var bodyReader io.Reader
	if len(details.Body) > 0 {
		bodyReader = bytes.NewReader(details.Body)
	}

	// Собираем сам запрос
	req, err := http.NewRequest(details.Method, e.baseURL+details.URL, bodyReader)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Устанавливаем заголовки, если есть
	for k, v := range details.Headers {
		req.Header.Set(k, v)
	}

	// Выполняем запрос
	log.Printf("HTTP request: %v", req)

	resp, err := e.httpClient.Do(req)
	if details.ExpectedStatus != 0 {
		if resp.StatusCode != details.ExpectedStatus {
			return "", fmt.Errorf("HTTP request failed, expected status %d, was %d", details.ExpectedStatus, resp.StatusCode)
		}
	}

	log.Printf("HTTP response: %v", resp)

	body, err := getBody(resp)
	if err != nil {
		return "", fmt.Errorf("Decode response failed: %w", err)
	}
	expectedBody := details.Body
	if err != nil {
		return "", fmt.Errorf("failed to marshal expected response body: %w", err)
	}

	if len(expectedBody) != 0 {
		if !compareBodies(body, expectedBody) {
			return "", fmt.Errorf("expected response body did not match actual response body, expected %s, got %s", expectedBody, body)
		}
	}

	return fmt.Sprintf("HTTP/1.1 %d %s", resp.StatusCode, resp.Status), nil
}

func getBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return bodyBytes, nil
}

func compareBodies(actualBody, expectedBody []byte) bool {
	var actual interface{}
	var expected interface{}

	json.Unmarshal(actualBody, &actual)
	json.Unmarshal(expectedBody, &expected)

	return reflect.DeepEqual(actual, expected)
}
