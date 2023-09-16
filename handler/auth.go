package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"payint"
	"time"
)

func (h *Handler) singUpAdmin(w http.ResponseWriter, r *http.Request) {
	var input payint.User

	json.NewDecoder(r.Body).Decode(&input)

	err := h.services.Authorization.CreateAdmin(input)
	if err != nil {
		log.Println("User creation error")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte(fmt.Sprintf("Admin %s created", input.Name)))
}

func (h *Handler) singUpUser(w http.ResponseWriter, r *http.Request) {
	var input payint.User

	json.NewDecoder(r.Body).Decode(&input)

	err := h.services.Authorization.CreateUser(input)
	if err != nil {
		log.Println("User creation error")
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Write([]byte(fmt.Sprintf("User %s created", input.Name)))
}

func (h *Handler) singIn(w http.ResponseWriter, r *http.Request) {
	var input payint.Input

	json.NewDecoder(r.Body).Decode(&input)

	tokenString, err := h.services.Authorization.GenerateToken(input.Name, input.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: time.Now().Add(5 * time.Minute),
	})
	w.Write([]byte(fmt.Sprintf("Welcome %s", input.Name)))

}
