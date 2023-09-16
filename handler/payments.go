package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"payint"
)

func (h *Handler) createPayment(w http.ResponseWriter, r *http.Request) {
	var payment payint.Payment

	json.NewDecoder(r.Body).Decode(&payment)

	responseString, err := h.services.CreatePayment(payment)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(responseString)

}

func (h *Handler) getPaymentsById(w http.ResponseWriter, r *http.Request) {

	payments, err := h.services.GetPaymentsById()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(payments)
}

func (h *Handler) getPaymentsDate(w http.ResponseWriter, r *http.Request) {

	payments, err := h.services.GetPaymentsDate()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(payments)
}

func (h *Handler) replenishAccount(w http.ResponseWriter, r *http.Request) {
	var input payint.InputPayment

	json.NewDecoder(r.Body).Decode(&input)

	responseString, err := h.services.ReplenishAccount(input.PayerName, input.AmountPayment)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(responseString)
}
