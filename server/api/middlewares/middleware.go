package middlewares

import "net/http"

type HandlerFuncReturningRequest func(http.ResponseWriter, *http.Request) *http.Request
type Middleware func(next HandlerFuncReturningRequest) HandlerFuncReturningRequest
