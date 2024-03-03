package webhook

import "net/http"

type HttpRequest struct {
	ConversationID string      `json:"conversation_id"`
	Headers        http.Header `json:"headers"`
	Body           []byte      `json:"body"`
	Path           string      `json:"path"`
	ContentLength  int64       `json:"content_length"`
	Method         string      `json:"method"`
}

type HttpResponse struct {
	MessageType    int
	StatusCode     int         `json:"status_code"`
	Headers        http.Header `json:"headers"`
	Body           []byte      `json:"body"`
	Error          error       `json:"error"`
	ContentLength  int64       `json:"content_length"`
	ConversationID string      `json:"conversation_id"`
}
