package router

import (
	"net/http"

	"github.com/kobayashiyabako16g/passkey-auth-example/internal/ui/handler"
)

type Router struct {
	ah handler.Auth
}

func NewRouter(ah handler.Auth) Router {
	return Router{ah}
}

func (r *Router) HandleRequest(mux *http.ServeMux) {
	mux.Handle("POST /passkey/register/start", http.HandlerFunc(r.ah.BeginRegistration))
	mux.Handle("POST /passkey/register/finish", http.HandlerFunc(r.ah.FinishRegistration))
	mux.Handle("POST /passkey/login/start", http.HandlerFunc(r.ah.BeginLogin))
	mux.Handle("POST /passkey/login/finish", http.HandlerFunc(r.ah.FinishLogin))
}
