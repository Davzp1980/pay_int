package repository

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"payint"
	"strconv"
)

type AccountsPostgres struct {
	db *sql.DB
}

func NewAccountsPostgres(db *sql.DB) *AccountsPostgres {
	return &AccountsPostgres{db: db}
}

func (a *AccountsPostgres) CreateAccount(name string) (string, error) {

	var user payint.User
	var account payint.Account

	err := a.db.QueryRow("SELECT id FROM users WHERE name=$1", name).Scan(&user.ID)
	if err != nil {
		log.Println("User does not exists")

		return "", err
	}
	i := strconv.Itoa(rand.Intn(1000000000))
	iban := i + name
	fmt.Println(iban)

	err = a.db.QueryRow("INSERT INTO accounts (user_id, iban) VALUES ($1,$2) RETURNING id", user.ID, iban).Scan(
		&account.ID)
	if err != nil {
		log.Println("Create account error")

		return "", err
	}
	return iban, nil
}

func (a *AccountsPostgres) BlockAccount(iban string) (string, error) {

	var account payint.Account

	err := a.db.QueryRow("SELECT iban, blocked FROM accounts WHERE iban=$1", iban).Scan(&account.Iban, &account.Blocked)
	if err != nil {
		log.Println(err)
	}
	expectedIban := account.Iban
	inputIban := iban

	if inputIban == expectedIban {
		_, err := a.db.Exec("UPDATE accounts SET blocked=$1 WHERE iban=$2", true, iban)
		if err != nil {
			log.Println(err)
			return "", err
		}
	} else {
		log.Println("Account", iban, "Does not exists")
		return "", err
	}
	return iban, nil
}

func (a *AccountsPostgres) UnBlockAccount(iban string) (string, error) {

	var account payint.Account

	err := a.db.QueryRow("SELECT iban, blocked FROM accounts WHERE iban=$1", iban).Scan(&account.Iban, &account.Blocked)
	if err != nil {
		log.Println(err)
		return "", err
	}
	expectedIban := account.Iban
	inputIban := iban

	if inputIban == expectedIban {
		_, err := a.db.Exec("UPDATE accounts SET blocked=$1 WHERE iban=$2", false, iban)
		if err != nil {
			log.Println(err)
			return "", err
		}
	} else {
		log.Println("Account", inputIban, "Does not exists")
		return "", err
	}
	return iban, nil
}

func (a *AccountsPostgres) GetAccountsById() ([]payint.OutputAccounts, error) {

	sortedAccounts := []payint.OutputAccounts{}

	rows, err := a.db.Query("SELECT id, user_id, iban, balance  FROM accounts ORDER BY id")
	if err != nil {
		log.Println("Account selection error")
		return sortedAccounts, err

	}

	for rows.Next() {
		var a payint.OutputAccounts

		if err = rows.Scan(&a.Id, &a.User_id, &a.Iban, &a.Balance); err != nil {
			log.Println(err)
			return sortedAccounts, err
		}
		sortedAccounts = append(sortedAccounts, a)
	}
	return sortedAccounts, nil
}

func (a *AccountsPostgres) GetAccountsByIban() ([]payint.OutputAccounts, error) {

	sortedAccounts := []payint.OutputAccounts{}

	rows, err := a.db.Query("SELECT * FROM accounts ORDER BY iban")
	if err != nil {
		log.Println("Account selection error")
		return sortedAccounts, err

	}

	for rows.Next() {
		var a payint.OutputAccounts

		if err = rows.Scan(&a.Id, &a.User_id, &a.Iban, &a.Balance); err != nil {
			log.Println(err)
			return sortedAccounts, err
		}
		sortedAccounts = append(sortedAccounts, a)
	}
	return sortedAccounts, nil

}

func (a *AccountsPostgres) GetAccountsByBalance() ([]payint.OutputAccounts, error) {

	sortedAccounts := []payint.OutputAccounts{}

	rows, err := a.db.Query("SELECT * FROM accounts ORDER BY balance")
	if err != nil {
		log.Println("Account selection error")
		return sortedAccounts, err

	}

	for rows.Next() {
		var a payint.OutputAccounts

		if err = rows.Scan(&a.Id, &a.User_id, &a.Iban, &a.Balance); err != nil {
			log.Println(err)
			return sortedAccounts, err
		}
		sortedAccounts = append(sortedAccounts, a)
	}
	return sortedAccounts, nil

}
