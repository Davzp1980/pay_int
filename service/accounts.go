package service

import (
	"payint"
	"payint/repository"
)

type AccountsService struct {
	repo repository.Account
}

func NewAccountsService(repo repository.Account) *AccountsService {
	return &AccountsService{repo: repo}
}

func (a *AccountsService) CreateAccount(name string) (string, error) {
	return a.repo.CreateAccount(name)
}

func (a *AccountsService) BlockAccount(iban string) (string, error) {
	return a.repo.BlockAccount(iban)
}

func (a *AccountsService) UnBlockAccount(iban string) (string, error) {
	return a.repo.UnBlockAccount(iban)
}

func (a *AccountsService) GetAccountsById() ([]payint.OutputAccounts, error) {
	return a.repo.GetAccountsById()
}

func (a *AccountsService) GetAccountsByIban() ([]payint.OutputAccounts, error) {
	return a.repo.GetAccountsByIban()
}

func (a *AccountsService) GetAccountsByBalance() ([]payint.OutputAccounts, error) {
	return a.repo.GetAccountsByBalance()
}
