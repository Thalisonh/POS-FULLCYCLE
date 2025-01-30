package web

import (
	"encoding/json"
	"net/http"

	"github.com/thalisonh/20-CleanArch/internal/entity"
	"github.com/thalisonh/20-CleanArch/internal/usecase"
)

type WebOrderHandler struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewWebOrderHandler(
	OrderRepository entity.OrderRepositoryInterface,
) *WebOrderHandler {
	return &WebOrderHandler{
		OrderRepository: OrderRepository,
	}
}

func (h *WebOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createOrder := usecase.NewCreateOrderUseCase(h.OrderRepository)
	output, err := createOrder.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebOrderHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	createOrder := usecase.NewListOrderUseCase(h.OrderRepository)
	output, err := createOrder.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
