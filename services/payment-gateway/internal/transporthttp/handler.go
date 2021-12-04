package transporthttp

import (
	"github.com/gorilla/mux"
	"net/http"
)


func HandleRoutes(h Handler) *mux.Router{
	r:= mux.NewRouter()
	r.HandleFunc("/authorize",h.AuthorizeHandler)
	r.HandleFunc("/capture", h.CaptureHandler)
	r.HandleFunc("/refund", h.RefundHandler)
	r.HandleFunc("/void", h.VoidHandler)
	return r
}



type Handler struct{

}

func NewHandler() (Handler,error){
	return Handler{},nil
}


func (h Handler) AuthorizeHandler(w http.ResponseWriter, r *http.Request){
	return http.s
}

func (h Handler) CaptureHandler(w http.ResponseWriter, r *http.Request){

}

func (h Handler) RefundHandler(w http.ResponseWriter, r *http.Request){

}

func (h Handler) VoidHandler(w http.ResponseWriter, r *http.Request){

}
