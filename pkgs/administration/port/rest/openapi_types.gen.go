// Package rest provides primitives to interact with the openapi HTTP API.
//
// Code generated by unknown module path version unknown version DO NOT EDIT.
package rest

const (
	BasicAuthScopes = "basicAuth.Scopes"
)

// AddInputReqBody defines model for AddInputReqBody.
type AddInputReqBody struct {
	ChatId      int64  `json:"chatId"`
	Description string `json:"description"`
	InputType   string `json:"inputType"`
	Name        string `json:"name"`
	UserId      int64  `json:"userId"`
}

// ErrorResp defines model for ErrorResp.
type ErrorResp struct {
	HttpStatus int     `json:"http_status"`
	Message    *string `json:"message,omitempty"`
	Name       string  `json:"name"`
}

// SuccessResp defines model for SuccessResp.
type SuccessResp struct {
	Success bool `json:"success"`
}

// AddInputJSONBody defines parameters for AddInput.
type AddInputJSONBody AddInputReqBody

// AddInputJSONRequestBody defines body for AddInput for application/json ContentType.
type AddInputJSONRequestBody AddInputJSONBody