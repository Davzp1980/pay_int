package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"payint"

	"payint/service"
	mock_service "payint/service/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Test_SignUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user payint.User)

	testTable := []struct {
		name                string
		inpitBody           string
		inpitUser           payint.User
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inpitBody: `{"name":"Test", "password":"test"}`,
			inpitUser: payint.User{
				Name:     "Test",
				Password: "test",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user payint.User) {
				s.EXPECT().CreateAdmin(user).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "Admin Test created",
		},
		{
			name:      "Empty fields",
			inpitBody: `{"name":"", "password":"test"}`,
			inpitUser: payint.User{
				Name:     "",
				Password: "test",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user payint.User) {
				s.EXPECT().CreateAdmin(user).Return(errors.New("Invalid input"))
			},
			expectedStatusCode: 400,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inpitUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			//Test server
			r := mux.NewRouter()
			r.HandleFunc("/sign-up", handler.singUpAdmin)

			//Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up",
				bytes.NewBufferString(testCase.inpitBody))

			//Perform request

			r.ServeHTTP(w, req)

			//Assert

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}

}
