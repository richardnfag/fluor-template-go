//go:generate go tool wit-bindgen-go generate --world function --out internal ./fluor-fun.wasm

package main

import (
	"bytes"
	"encoding/json"
	"net/url"
	"strings"

	"example.com/{{project-name}}/internal/fluor/fun/handler"
	"go.bytecodealliance.org/cm"
)

func init() {
	handler.Exports.Handle = handle
}

func handle(request handler.Request) cm.Result[handler.ResponseShape, handler.Response, handler.Error] {
	requestID := headerValue(request.Headers, "x-request-id")
	switch request.Method {
	case "GET":
		return listItems(request, requestID)
	case "POST":
		return createItem(request, requestID)
	case "PUT":
		return updateItem(request, requestID)
	case "DELETE":
		return deleteItem(request, requestID)
	default:
		return failure(handler.ErrorCodeBadRequest, "method "+request.Method+" is not supported")
	}
}

func listItems(request handler.Request, requestID string) cm.Result[handler.ResponseShape, handler.Response, handler.Error] {
	items := []map[string]string{
		{"name": "notebook"},
		{"name": "pencil"},
	}
	if name, ok := queryParam(request.Query, "name"); ok {
		if name != "notebook" && name != "pencil" {
			return failure(handler.ErrorCodeNotFound, "item "+name+" was not found")
		}
		items = []map[string]string{
			{"name": name},
		}
	}
	return respond(200, requestID, envelope(request, requestID, "items", items))
}

func createItem(request handler.Request, requestID string) cm.Result[handler.ResponseShape, handler.Response, handler.Error] {
	name, ok := jsonStringField(request.Body, "name")
	if !ok {
		return failure(handler.ErrorCodeBadRequest, "body must be a JSON object with a string name")
	}
	return respond(201, requestID, envelope(request, requestID, "created", map[string]string{"name": name}))
}

func updateItem(request handler.Request, requestID string) cm.Result[handler.ResponseShape, handler.Response, handler.Error] {
	id, ok := queryParam(request.Query, "id")
	if !ok {
		return failure(handler.ErrorCodeBadRequest, "query parameter id is required")
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(request.Body), &body); err != nil || body == nil {
		return failure(handler.ErrorCodeBadRequest, "body must be a JSON object")
	}
	return respond(200, requestID, envelope(request, requestID, "updated", map[string]any{"id": id, "body": body}))
}

func deleteItem(request handler.Request, requestID string) cm.Result[handler.ResponseShape, handler.Response, handler.Error] {
	id, ok := queryParam(request.Query, "id")
	if !ok {
		return failure(handler.ErrorCodeBadRequest, "query parameter id is required")
	}
	return respond(200, requestID, envelope(request, requestID, "deleted", id))
}

func envelope(request handler.Request, requestID, field string, value any) map[string]any {
	var query any
	if text := request.Query.Value(); request.Query.Some() != nil {
		query = text
	}
	var id any
	if requestID != "" {
		id = requestID
	}
	payload := map[string]any{
		"function":  "{{project-name}}",
		"language":  "go",
		"method":    request.Method,
		"path":      request.Path,
		"query":     query,
		"requestId": id,
	}
	payload[field] = value
	return payload
}

func respond(status uint16, requestID string, payload map[string]any) cm.Result[handler.ResponseShape, handler.Response, handler.Error] {
	var encoded bytes.Buffer
	encoder := json.NewEncoder(&encoded)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(payload)
	body := bytes.TrimRight(encoded.Bytes(), "\n")
	headers := []handler.Header{
		{Name: "content-type", Value: "application/json"},
	}
	if requestID != "" {
		headers = append(headers, handler.Header{Name: "x-request-id", Value: requestID})
	}
	return cm.OK[cm.Result[handler.ResponseShape, handler.Response, handler.Error]](handler.Response{
		Status:  status,
		Headers: cm.ToList(headers),
		Body:    string(body),
	})
}

func failure(code handler.ErrorCode, message string) cm.Result[handler.ResponseShape, handler.Response, handler.Error] {
	return cm.Err[cm.Result[handler.ResponseShape, handler.Response, handler.Error]](handler.Error{
		Code:    code,
		Message: message,
	})
}

func headerValue(headers cm.List[handler.Header], name string) string {
	for _, header := range headers.Slice() {
		if strings.EqualFold(header.Name, name) {
			return header.Value
		}
	}
	return ""
}

func queryParam(query cm.Option[string], key string) (string, bool) {
	text := query.Some()
	if text == nil {
		return "", false
	}
	for _, part := range strings.Split(*text, "&") {
		name, value, found := strings.Cut(part, "=")
		if !found || formDecode(name) != key {
			continue
		}
		return formDecode(value), true
	}
	return "", false
}

func formDecode(value string) string {
	decoded, err := url.QueryUnescape(value)
	if err != nil {
		return value
	}
	return decoded
}

func jsonStringField(body, field string) (string, bool) {
	var payload map[string]any
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		return "", false
	}
	value, ok := payload[field].(string)
	if !ok || value == "" {
		return "", false
	}
	return value, true
}

func main() {}
