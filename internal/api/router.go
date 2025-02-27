package api

import (
	"github.com/SaveljevRoman/go-layout-project/internal/service"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(userService *service.UserService, productService *service.ProductService, purchaseService *service.PurchaseService) http.Handler {
	router := mux.NewRouter()

	// Инициализация хендлеров
	userHandlers := NewUserHandlers(userService)
	productHandlers := NewProductHandlers(productService)
	purchaseHandlers := NewPurchaseHandlers(purchaseService)

	// Определение маршрутов

	// Группа маршрутов для пользователей
	userRouter := router.PathPrefix("/api/users").Subrouter()
	userRouter.HandleFunc("", userHandlers.GetAllUsers).Methods("GET")
	userRouter.HandleFunc("", userHandlers.CreateUser).Methods("POST")
	userRouter.HandleFunc("/{id:[0-9]+}", userHandlers.GetUser).Methods("GET")
	userRouter.HandleFunc("/{id:[0-9]+}", userHandlers.UpdateUser).Methods("PUT")
	userRouter.HandleFunc("/{id:[0-9]+}", userHandlers.DeleteUser).Methods("DELETE")
	userRouter.HandleFunc("/{user_id:[0-9]+}/purchases", purchaseHandlers.GetUserPurchases).Methods("GET")

	// Группа маршрутов для продуктов
	productRouter := router.PathPrefix("/api/products").Subrouter()
	productRouter.HandleFunc("", productHandlers.GetAllProducts).Methods("GET")
	productRouter.HandleFunc("", productHandlers.CreateProduct).Methods("POST")
	productRouter.HandleFunc("/{id:[0-9]+}", productHandlers.GetProduct).Methods("GET")
	productRouter.HandleFunc("/{id:[0-9]+}", productHandlers.UpdateProduct).Methods("PUT")
	productRouter.HandleFunc("/{id:[0-9]+}", productHandlers.DeleteProduct).Methods("DELETE")

	// Группа маршрутов для покупок
	purchaseRouter := router.PathPrefix("/api/purchases").Subrouter()
	purchaseRouter.HandleFunc("", purchaseHandlers.GetAllPurchases).Methods("GET")
	purchaseRouter.HandleFunc("", purchaseHandlers.CreatePurchase).Methods("POST")
	purchaseRouter.HandleFunc("/{id:[0-9]+}", purchaseHandlers.GetPurchase).Methods("GET")
	purchaseRouter.HandleFunc("/{id:[0-9]+}/status", purchaseHandlers.UpdatePurchaseStatus).Methods("PUT")

	// Промежуточное ПО
	router.Use(LoggingMiddleware)

	return router
}
