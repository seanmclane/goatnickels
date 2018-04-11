package rpc

import (
	"encoding/json"
)

var BroadcastChannel = make(chan JsonRpcMessage)

type JsonRpcMessage struct {
	Version string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Error   *jsonRpcError   `json:"error,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
}

func BuildNotification(Method string, Params json.RawMessage) (msg JsonRpcMessage) {
	msg = JsonRpcMessage{
		Version: "2.0",
		Method:  Method,
		Params:  Params,
	}
	return msg
}

func BuildRequest(Id json.RawMessage, Method string, Params json.RawMessage) (msg JsonRpcMessage) {
	msg = JsonRpcMessage{
		Version: "2.0",
		ID:      Id,
		Method:  Method,
		Params:  Params,
	}
	return msg
}

func BuildResponse(Id json.RawMessage, Result json.RawMessage, Error *jsonRpcError) (msg JsonRpcMessage) {
	msg = JsonRpcMessage{
		Version: "2.0",
		ID:      Id,
		Result:  Result,
		Error:   Error,
	}
	return msg
}

//TODO add json rpc error codes

type jsonRpcError struct {
	Code    string
	Message string
	Data    string
}
