package api

import (
	"encoding/json"
	"github.com/SaveljevRoman/go-layout-project/internal/models"
	"github.com/SaveljevRoman/go-layout-project/internal/service"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type PurchaseHandlers struct {
	purchaseService *service.PurchaseService
}

func NewPurchaseHandlers(purchaseService *service.PurchaseService) *PurchaseHandlers {
	return &PurchaseHandlers{
		purchaseService: purchaseService,
	}
}

func (h *PurchaseHandlers) CreatePurchase(w http.ResponseWriter, r *http.Request) {
	var request models.PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	purchase, err := h.purchaseService.CreatePurchase(ctx, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(purchase)
}

func (h *PurchaseHandlers) GetPurchase(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID покупки", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	purchase, err := h.purchaseService.GetPurchase(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if purchase == nil {
		http.Error(w, "Покупка не найдена", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(purchase)
}

func (h *PurchaseHandlers) GetUserPurchases(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["user_id"], 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	purchases, err := h.purchaseService.GetUserPurchases(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(purchases)
}

func (h *PurchaseHandlers) UpdatePurchaseStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID покупки", http.StatusBadRequest)
		return
	}

	var statusRequest struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&statusRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.purchaseService.UpdatePurchaseStatus(ctx, id, statusRequest.Status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	purchase, err := h.purchaseService.GetPurchase(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(purchase)
}

func (h *PurchaseHandlers) GetAllPurchases(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	purchases, err := h.purchaseService.GetAllPurchases(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(purchases)
}
