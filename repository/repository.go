package repository

import (
	"database/sql"
	"payint"
)

type Authorization interface {
	CreateAdmin(admin payint.User) error
	CreateUser(user payint.User) error
	BlockUser(user payint.User) error
	UnBlockUser(user payint.User) error
	GetUser(username, password string) (payint.User, error)
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

type Repository struct {
	Authorization
	Account
	Payment
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Account:       NewAccountsPostgres(db),
		Payment:       NewPaymentPostgres(db),
	}
}
