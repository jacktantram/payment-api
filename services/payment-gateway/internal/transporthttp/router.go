package transporthttp

import (
	"github.com/gorilla/mux"
)

func HandleRoutes(h Handler) *mux.Router{
	r:= mux.NewRouter()
	r.HandleFunc("/authorize",h.AuthorizeHandler)
	r.HandleFunc("/capture", h.CaptureHandler)
	r.HandleFunc("/refund", h.RefundHandler)
	r.HandleFunc("/void", h.VoidHandler)
	return r
}

