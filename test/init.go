package test

import "net/http"

// RequestHandler is a handler used for testing purposes
var RequestHandler func(*http.Request) *http.Response
