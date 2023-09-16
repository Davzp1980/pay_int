package service

import (
	"payint"
	"payint/repository"
)

type PaymentsService struct {
	repo repository.Payment
}

func NewPaymentsService(repo repository.Payment) *PaymentsService {
	return &PaymentsService{repo: repo}
}

func (p *PaymentsService) CreatePayment(payment payint.Payment) (string, error) {
	return p.repo.CreatePayment(payment)
}

func (p *PaymentsService) GetPaymentsById() ([]payint.Payment, error) {
	return p.repo.GetPaymentsById()
}

func (p *PaymentsService) GetPaymentsDate() ([]payint.Payment, error) {
	return p.repo.GetPaymentsDate()
}

func (p *PaymentsService) ReplenishAccount(name string, deposit int) (string, error) {
	return p.repo.ReplenishAccount(name, deposit)
}
