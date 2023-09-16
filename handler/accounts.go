package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"payint"
)

func (h *Handler) createAccount(w http.ResponseWriter, r *http.Request) {
	var input payint.Input

	json.NewDecoder(r.Body).Decode(&input)

	iban, err := h.services.CreateAccount(input.Name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Write([]byte(fmt.Sprintf("Account %s created", iban)))

}

func (h *Handler) blockAccount(w http.ResponseWriter, r *http.Request) {
	var input payint.Input

	json.NewDecoder(r.Body).Decode(&input)

	iban, err := h.services.BlockAccount(input.Iban)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Write([]byte(fmt.Sprintf("Account %s blocked", iban)))

}

func (h *Handler) unblockAccount(w http.ResponseWriter, r *http.Request) {
	var input payint.Input

	json.NewDecoder(r.Body).Decode(&input)

	iban, err := h.services.BlockAccount(input.Iban)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Write([]byte(fmt.Sprintf("Account %s blocked", iban)))

}

func (h *Handler) getAccountsById(w http.ResponseWriter, r *http.Request) {

	accounts, err := h.services.GetAccountsById()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(accounts)

}

func (h *Handler) getAccountsByIban(w http.ResponseWriter, r *http.Request) {

	accounts, err := h.services.GetAccountsByIban()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(accounts)

}

func (h *Handler) getAccountsByBalance(w http.ResponseWriter, r *http.Request) {

	accounts, err := h.services.GetAccountsByBalance()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(accounts)

}
