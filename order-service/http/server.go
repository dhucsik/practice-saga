package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"order-service/models"
	"order-service/service"
)

type Server struct {
	srv service.Service
}

func NewServer(srv service.Service) *Server {
	return &Server{srv: srv}
}

func (s *Server) Start() error {
	http.HandleFunc("/orders", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost {
			var order models.Order

			data, err := io.ReadAll(request.Body)
			if err != nil {
				http.Error(writer, "error reading body", http.StatusBadRequest)
				return
			}
			err = json.Unmarshal(data, &order)
			if err != nil {
				http.Error(writer, "unmarshal error", http.StatusBadRequest)
				return
			}

			err = s.srv.CreateOrder(request.Context(), &order)
			if err != nil {
				http.Error(writer, "error: "+err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Fprint(writer, "success")
		}
	})

	return http.ListenAndServe(":8080", nil)
}
