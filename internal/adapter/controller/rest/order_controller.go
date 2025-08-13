package rest

import (
	"encoding/json"
	"errors"
	"github.com/folivorra/get_order/internal/adapter/mapper"
	"github.com/folivorra/get_order/internal/config"
	"github.com/folivorra/get_order/internal/usecase"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

type Controller struct {
	service *usecase.OrderService
	logger  *slog.Logger
	cfg     config.Config
}

func NewController(service *usecase.OrderService, cfg config.Config, logger *slog.Logger) *Controller {
	return &Controller{
		service: service,
		logger:  logger,
		cfg:     cfg,
	}
}

func (c *Controller) GetOrderToUI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uidV := vars["uid"]

	uid, err := uuid.Parse(uidV)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := c.service.GetOrder(r.Context(), uid)
	if err != nil {
		if errors.Is(err, usecase.OrderDoesNotExists) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	orderDTO := mapper.ConvertFromDomain(order)

	if err = json.NewEncoder(w).Encode(orderDTO); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Controller) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/orders/{uid}", c.GetOrderToUI).Methods("GET")
}
