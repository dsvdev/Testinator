package llm

import "encoding/json"

type ActionType string

const (
	ActionHTTPRequest ActionType = "http_request"
)

type LLMResponse struct {
	Action  ActionType      `json:"action"`
	Details json.RawMessage `json:"details"` // далее дессериализуется вручную в нужный тип
}

type HTTPRequestDetails struct {
	Method         string            `json:"method"`            // e.g. "POST"
	URL            string            `json:"url"`               // e.g. "/users/123"
	Headers        map[string]string `json:"headers,omitempty"` // optional
	Body           json.RawMessage   `json:"body,omitempty"`    // can be map or struct
	ExpectedStatus int               `json:"expected_status"`
	ExpectedBody   json.RawMessage   `json:"expected_body,omitempty"` // optional
}
