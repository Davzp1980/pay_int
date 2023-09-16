package repository

import (
	"payint"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/stretchr/testify/assert"
)

func Test_CreateAdmin(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewAuthPostgres(db)

	tests := []struct {
		name    string
		mock    func()
		input   payint.User
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {

				mock.ExpectQuery("INSERT INTO users").WithArgs("test", "password", true)
			},
			input: payint.User{
				Name:     "test",
				Password: "password",
				IsAdmin:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mock()

			tt.mock()

			err := r.CreateUser(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

			}
			assert.NoError(t, mock.ExpectationsWereMet())

		})
	}
}
