package handler

import "net/http"

type Router struct {
	ah Auth
}

func NewRouter(ah Auth) Router {
	return Router{ah}
}

func (r *Router) HandleRequest(mux *http.ServeMux) {
	mux.Handle("POST /passkey/register/start", http.HandlerFunc(r.ah.BeginRegistration))
	mux.Handle("POST /passkey/register/finish", http.HandlerFunc(r.ah.FinishRegistration))
}
