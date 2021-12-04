package transporthttp

import "net/http"

type Handler struct{

}

func NewHandler() (Handler,error){
	return Handler{},nil
}


func (h Handler) AuthorizeHandler(w http.ResponseWriter, r *http.Request){

}

func (h Handler) CaptureHandler(w http.ResponseWriter, r *http.Request){

}

func (h Handler) RefundHandler(w http.ResponseWriter, r *http.Request){

}

func (h Handler) VoidHandler(w http.ResponseWriter, r *http.Request){

}
