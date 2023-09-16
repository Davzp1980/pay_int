package service

import (
	"payint"
	"payint/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go
type Authorization interface {
	CreateAdmin(admin payint.User) error
	CreateUser(user payint.User) error
	BlockUser(user payint.User) error
	UnBlockUser(user payint.User) error
	GetUser(username, password string) (payint.User, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(signedToken string) (string, error)
}

type Account interface {
	CreateAccount(name string) (string, error)
	BlockAccount(iban string) (string, error)
	UnBlockAccount(iban string) (string, error)
	GetAccountsById() ([]payint.OutputAccounts, error)
	GetAccountsByIban() ([]payint.OutputAccounts, error)
	GetAccountsByBalance() ([]payint.OutputAccounts, error)
}

type Payment interface {
	CreatePayment(payment payint.Payment) (string, error)
	GetPaymentsById() ([]payint.Payment, error)
	GetPaymentsDate() ([]payint.Payment, error)
	ReplenishAccount(name string, deposit int) (string, error)
}

type Service struct {
	Authorization
	Account
	Payment
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Account:       NewAccountsService(repos.Account),
		Payment:       NewPaymentsService(repos),
	}
}
