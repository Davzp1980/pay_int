package handler

import (
	"payint/service"

	"github.com/gorilla/mux"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRouters() *mux.Router {

	router := mux.NewRouter()

	router.Use(h.userIdentity)
	router.HandleFunc("/sign-up-admin", h.singUpAdmin).Methods("POST")
	router.HandleFunc("/sign-up-user", h.singUpUser).Methods("POST")
	router.HandleFunc("/sign-in", h.singIn).Methods("POST")

	router.HandleFunc("/create-account", h.createAccount).Methods("POST")
	router.HandleFunc("/block-account", h.blockAccount).Methods("POST")
	router.HandleFunc("/unblock-account", h.unblockAccount).Methods("POST")
	router.HandleFunc("/get-account-by-id", h.getAccountsById).Methods("GET")
	router.HandleFunc("/get-account-by-iban", h.getAccountsByIban).Methods("GET")
	router.HandleFunc("/get-account-by-balance", h.getAccountsByBalance).Methods("GET")

	router.HandleFunc("/create-payment", h.createPayment).Methods("POST")
	router.HandleFunc("/replenish-account", h.replenishAccount).Methods("POST")
	router.HandleFunc("/get-payment-by-id", h.getPaymentsById).Methods("GET")
	router.HandleFunc("/get-payment-by-date", h.getPaymentsDate).Methods("GET")

	return router

}
