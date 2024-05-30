package handler

import "net/http"

func (h *Handler) Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I'm Alive"))
}
