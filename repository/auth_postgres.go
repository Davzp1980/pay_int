package repository

import (
	"database/sql"
	"errors"
	"log"
	"payint"

	"golang.org/x/crypto/bcrypt"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (a *AuthPostgres) CreateAdmin(admin payint.User) error {

	isAdmin := true
	hashedPassword, _ := HashePassword(admin.Password)

	_, err := a.db.Query("INSERT INTO users (name, password, is_admin) VALUES ($1,$2,$3)",
		admin.Name, hashedPassword, isAdmin)
	if err != nil {

		return errors.New("Invalid input")
	}
	return errors.New("")
}

func (a *AuthPostgres) CreateUser(user payint.User) error {

	isAdmin := false
	hashedPassword, _ := HashePassword(user.Password)

	_, err := a.db.Query("INSERT INTO users (name, password, is_admin) VALUES ($1,$2,$3)",
		user.Name, hashedPassword, isAdmin)
	if err != nil {

		return err
	}
	return nil
}

func (a *AuthPostgres) BlockUser(user payint.User) error {
	var userFromDB payint.User

	err := a.db.QueryRow("SELECT name FROM users WHERE name=$1", user.Name).Scan(&userFromDB.Name)
	if err != nil {
		log.Println(err)
		return err
	}
	expectedName := userFromDB.Name
	inputName := user.Name

	if inputName == expectedName {
		_, err := a.db.Exec("UPDATE users SET blocked=$1 WHERE name=$2", true, user.Name)
		if err != nil {
			log.Println(err)
			return err
		}

	} else {
		log.Println("User", inputName, "Does not exists")
		return errors.New("User does not exists")
	}
	return nil

}

func (a *AuthPostgres) UnBlockUser(user payint.User) error {
	var userFromDB payint.User

	err := a.db.QueryRow("SELECT name FROM users WHERE name=$1", user.Name).Scan(&userFromDB.Name)
	if err != nil {
		log.Println(err)
		return err
	}
	expectedName := userFromDB.Name
	inputName := user.Name

	if inputName == expectedName {
		_, err := a.db.Exec("UPDATE users SET blocked=$1 WHERE name=$2", false, user.Name)
		if err != nil {
			log.Println(err)
			return err

		}

	} else {
		log.Println("User", inputName, "Does not exists")
		return errors.New("User does not exists")
	}
	return nil

}

func (a *AuthPostgres) ChangeUserPassword(user payint.User) error {

	hashedPassword, _ := HashePassword(user.Password)

	_, err := a.db.Exec("UPDATE users SET password=$1 WHERE name=$2;", hashedPassword, user.Name)
	if err != nil {
		log.Println("Change password error")
		return errors.New("Change password error")
	}
	return nil

}

func (a *AuthPostgres) GetUser(username, password string) (payint.User, error) {
	var user payint.User

	err := a.db.QueryRow("SELECT * FROM users WHERE name=$1", username).Scan(
		&user.ID, &user.Name, &user.Password, &user.IsAdmin, &user.Blocked)

	return user, err
}

func HashePassword(password string) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
