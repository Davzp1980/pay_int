package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"payint"
	"time"
)

type PaymentPostgres struct {
	db *sql.DB
}

func NewPaymentPostgres(db *sql.DB) *PaymentPostgres {
	return &PaymentPostgres{db: db}
}

func (p *PaymentPostgres) CreatePayment(payment payint.Payment) (string, error) {

	var payer payint.User
	var receiverAccount payint.Account
	var payerAccount payint.Account

	//по имени отправителя получаем id
	err := p.db.QueryRow("SELECT id, name FROM users WHERE name=$1", payment.Payer).Scan(&payer.ID, &payer.Name)
	if err != nil {
		log.Println("User does not exists")
		return "Payment did not make", errors.New("User does not exists")
	}
	// по номеу счета (iban) получаем id, Iban, Balance получателя
	err = p.db.QueryRow("SELECT id, user_id, iban, balance, blocked  FROM accounts WHERE iban=$1", payment.RecieverIban).Scan(
		&receiverAccount.ID, &receiverAccount.UserId, &receiverAccount.Iban, &receiverAccount.Balance, &receiverAccount.Blocked)
	if err != nil {
		log.Println("Account does not exists")
		return "Payment did not make", errors.New("Account does not exists")
	}
	if receiverAccount.Blocked {
		log.Println("The recipient’s account is blocked")
		return "Payment did not make", errors.New("The recipient’s account is blocked")
	}
	// создаем платеж
	err = p.db.QueryRow("INSERT INTO payments (user_id, reciever, reciever_iban, payer, payer_iban, amount_payment, date) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id",
		receiverAccount.UserId, payment.Reciever, payment.RecieverIban, payment.Payer, payment.PayerIban, payment.AmountPayment, time.Now()).Scan(
		&payment.ID)
	if err != nil {
		log.Println("Payment creation error")
		return "Payment did not make", errors.New("Payment creation error")
	}
	// проверяем достаточно ли денег на счете отправителя и заблокирован ли он, снимаем сумму платежа со счета
	err = p.db.QueryRow("SELECT balance, blocked FROM accounts WHERE iban=$1", payment.PayerIban).Scan(&payerAccount.Balance, &payerAccount.Blocked)
	if err != nil {
		log.Println("Payer's balance is incorrect")
		return "Payment did not make", errors.New("Payer's balance is incorrect")
	}
	if payerAccount.Blocked {
		log.Println("Payer's account is blocked")
		return "Payment did not make", errors.New("Payer's account is blocked")
	}

	if payerAccount.Balance < payment.AmountPayment {
		log.Println("Not enough money in the account")
		return "Payment did not make", errors.New("Not enough money in the account")
	}
	payerBalance := payerAccount.Balance - payment.AmountPayment

	_, err = p.db.Exec("UPDATE accounts SET balance=$2 WHERE iban=$1", payment.PayerIban, payerBalance)
	if err != nil {
		log.Println("Payer's balance update error")
		return "Payment did not make", errors.New("Payer's balance update error")
	}
	// изменяем баланс получателя в соответствии с указынным номером счета (iban) и суммой платежа
	balance := receiverAccount.Balance + payment.AmountPayment

	_, err = p.db.Exec("UPDATE accounts SET balance=$2 WHERE iban=$1", payment.RecieverIban, balance)
	if err != nil {
		log.Println("Recipient's balance update error")
		return "Payment did not make", errors.New("Recipient's balance update error")
	}

	return fmt.Sprintf("You are sent %v to %s", payment.AmountPayment, payment.Reciever), nil

}

func (p *PaymentPostgres) GetPaymentsById() ([]payint.Payment, error) {
	sortedPayments := []payint.Payment{}
	rows, err := p.db.Query("SELECT * FROM payments ORDER BY id")
	if err != nil {
		log.Println("Account selection error")
		return sortedPayments, err

	}

	for rows.Next() {
		var p payint.Payment

		if err = rows.Scan(&p.ID, &p.UserId, &p.Reciever, &p.RecieverIban, &p.Payer, &p.PayerIban, &p.AmountPayment, &p.Date); err != nil {
			log.Println(err)
			return sortedPayments, err
		}
		sortedPayments = append(sortedPayments, p)
	}

	return sortedPayments, nil

}

func (p *PaymentPostgres) GetPaymentsDate() ([]payint.Payment, error) {
	sortedPayments := []payint.Payment{}

	rows, err := p.db.Query("SELECT * FROM payments ORDER BY date DESC")
	if err != nil {
		log.Println("Account selection error")

	}

	for rows.Next() {
		var p payint.Payment

		if err = rows.Scan(&p.ID, &p.UserId, &p.Reciever, &p.RecieverIban, &p.Payer, &p.PayerIban, &p.AmountPayment, &p.Date); err != nil {
			log.Println(err)
			return sortedPayments, err
		}
		sortedPayments = append(sortedPayments, p)
	}
	return sortedPayments, nil

}

func (p *PaymentPostgres) ReplenishAccount(name string, deposit int) (string, error) {

	var user payint.User
	var amountAccount int

	//balance:=input.AmountReplenish+
	err := p.db.QueryRow("SELECT id FROM users WHERE name=$1", name).Scan(&user.ID)
	if err != nil {

		log.Println("User does not exists")
		return "", errors.New("User does not exists")
	}

	err = p.db.QueryRow("SELECT balance FROM accounts WHERE user_id=$1", user.ID).Scan(&amountAccount)
	if err != nil {
		log.Println("Account does not exists")
		return "", errors.New("Account does not exists")
	}
	balance := deposit + amountAccount

	_, err = p.db.Exec("UPDATE accounts SET balance=$1", balance)
	if err != nil {
		log.Println("UPDATE account error")
		return "", errors.New("UPDATE account error")
	}

	return fmt.Sprintf("Account was replenished for %v. Amount in the account %v", deposit, balance), nil

}
