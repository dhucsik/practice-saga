package http

import (
	"fmt"
	"net/http"
	"payment-service/service"
	"strconv"
	"strings"
)

type Server struct {
	srv service.Service
}

func NewServer(srv service.Service) *Server {
	return &Server{srv: srv}
}

func (s *Server) Start() error {
	http.HandleFunc("/payments/", func(writer http.ResponseWriter, request *http.Request) {
		idStr := strings.TrimPrefix(request.URL.Path, "/payments/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(writer, "error incorrect id: "+err.Error(), http.StatusBadRequest)
			return
		}

		err = s.srv.PaidPayment(request.Context(), id)
		if err != nil {
			http.Error(writer, "error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(writer, "success")
	})

	return http.ListenAndServe(":8081", nil)
}
