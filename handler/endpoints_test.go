package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/asrul10/UserService/generated"
	"github.com/asrul10/UserService/helper"
	"github.com/asrul10/UserService/repository"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

func TestRegisterUser(t *testing.T) {
	// Mocking the repository
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := repository.NewMockRepositoryInterface(ctrl)
	h := helper.NewHelper(helper.NewHelperOptions{
		JwtPrivateKeyPath: "../storage/key.pem",
		JwtPublicKeyPath:  "../storage/key.pem.pub",
	})

	// Test cases
	tests := []struct {
		caseName     string
		payload      string
		mockFunc     func()
		expectedCode int
	}{
		{
			caseName:     "Empty payload",
			payload:      "",
			mockFunc:     func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			caseName: "Positive case",
			payload:  `{"phoneNumber":"+62123456789","fullName":"test","password":"Test123/"}`,
			mockFunc: func() {
				m.
					EXPECT().
					GetUserByPhoneNumber(gomock.Any(), gomock.Any()).
					Return(repository.GetUserByPhoneNumberOutput{}, errors.New("not found"))
				m.
					EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(repository.CreateUserOutput{UserId: 1}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			caseName: "Duplicate phone number",
			payload:  `{"phoneNumber":"+62123456789","fullName":"test","password":"Test123/"}`,
			mockFunc: func() {
				m.
					EXPECT().
					GetUserByPhoneNumber(gomock.Any(), gomock.Any()).
					Return(repository.GetUserByPhoneNumberOutput{UserId: 1}, nil)
			},
			expectedCode: http.StatusConflict,
		},
		{
			caseName:     "Invalid phone number",
			payload:      `{"phoneNumber":"123","fullName":"test","password":"Test123/"}`,
			mockFunc:     func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			caseName:     "Invalid full name",
			payload:      `{"phoneNumber":"+62123456789","fullName":"t","password":"Test123/"}`,
			mockFunc:     func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			caseName:     "Invalid password",
			payload:      `{"phoneNumber":"+62123456789","fullName":"test","password":"test"}`,
			mockFunc:     func() {},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.caseName, func(t *testing.T) {
			// Creating the server
			e := echo.New()
			server := NewServer(NewServerOptions{
				Repository: m,
				Helper:     h,
				Echo:       e,
			})
			generated.RegisterHandlers(e, server)

			req := httptest.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(test.payload),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			test.mockFunc()

			if err := server.RegisterUser(c); err != nil {
				t.Errorf("Error: %v", err)
			}
			if rec.Code != test.expectedCode {
				t.Errorf("Expected %d, got %d", test.expectedCode, rec.Code)
			}
		})
	}
}

func TestLoginUser(t *testing.T) {
	// Mocking the repository
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := repository.NewMockRepositoryInterface(ctrl)
	h := helper.NewHelper(helper.NewHelperOptions{
		JwtPrivateKeyPath: "../storage/key.pem",
		JwtPublicKeyPath:  "../storage/key.pem.pub",
	})

	// Test cases
	tests := []struct {
		caseName     string
		payload      string
		mockFunc     func()
		expectedCode int
	}{
		{
			caseName:     "Empty payload",
			payload:      "",
			mockFunc:     func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			caseName: "Positive case",
			payload:  `{"phoneNumber":"+62123456789","password":"Test123/"}`,
			mockFunc: func() {
				hashPassword, _ := h.HashPassword("Test123/")
				m.
					EXPECT().
					GetUserByPhoneNumber(gomock.Any(), gomock.Any()).
					Return(repository.GetUserByPhoneNumberOutput{
						UserId:   1,
						Password: hashPassword,
					}, nil)
				m.
					EXPECT().
					SuccessLoginCount(gomock.Any(), gomock.Any()).
					Return(repository.SuccessLoginCountOutput{
						UserId: 1,
					}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			caseName: "User not found",
			payload:  `{"phoneNumber":"+62123456789","password":"Test123/"}`,
			mockFunc: func() {
				m.
					EXPECT().
					GetUserByPhoneNumber(gomock.Any(), gomock.Any()).
					Return(repository.GetUserByPhoneNumberOutput{}, errors.New("not found"))
			},
			expectedCode: http.StatusNotFound,
		},
		{
			caseName: "Invalid password",
			payload:  `{"phoneNumber":"+62123456789","password":"WrongPassword12/"}`,
			mockFunc: func() {
				hashPassword, _ := h.HashPassword("Test123/")
				m.
					EXPECT().
					GetUserByPhoneNumber(gomock.Any(), gomock.Any()).
					Return(repository.GetUserByPhoneNumberOutput{
						UserId:   1,
						Password: hashPassword,
					}, nil)
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.caseName, func(t *testing.T) {
			// Creating the server
			e := echo.New()
			server := NewServer(NewServerOptions{
				Repository: m,
				Helper:     h,
				Echo:       e,
			})
			generated.RegisterHandlers(e, server)

			req := httptest.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(test.payload),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			test.mockFunc()

			if err := server.LoginUser(c); err != nil {
				t.Errorf("Error: %v", err)
			}
			if rec.Code != test.expectedCode {
				t.Errorf("Expected %d, got %d", test.expectedCode, rec.Code)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	// Mocking the repository
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := repository.NewMockRepositoryInterface(ctrl)
	h := helper.NewHelper(helper.NewHelperOptions{
		JwtPrivateKeyPath: "../storage/key.pem",
		JwtPublicKeyPath:  "../storage/key.pem.pub",
	})

	// Test cases
	tests := []struct {
		caseName     string
		token        func() string
		mockFunc     func()
		expectedCode int
	}{
		{
			caseName: "Positive case",
			token: func() string {
				token := ""
				h.GenerateAccessToken(&token, 1)
				return token
			},
			mockFunc: func() {
				m.
					EXPECT().
					GetUserById(gomock.Any(), gomock.Any()).
					Return(repository.GetUserByIdOutput{
						UserId:      1,
						PhoneNumber: "+62123456789",
						FullName:    "test",
					}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			caseName: "Unauthorized",
			token: func() string {
				return ""
			},
			mockFunc:     func() {},
			expectedCode: http.StatusForbidden,
		},
	}

	for _, test := range tests {
		t.Run(test.caseName, func(t *testing.T) {
			// Creating the server
			e := echo.New()
			server := NewServer(NewServerOptions{
				Repository: m,
				Helper:     h,
				Echo:       e,
			})
			generated.RegisterHandlers(e, server)

			req := httptest.NewRequest(
				http.MethodGet,
				"/",
				nil,
			)
			// set bearer token
			req.Header.Set("Authorization", "Bearer "+test.token())
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			test.mockFunc()

			if err := server.GetUser(c); err != nil {
				t.Errorf("Error: %v", err)
			}
			if rec.Code != test.expectedCode {
				t.Errorf("Expected %d, got %d", test.expectedCode, rec.Code)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	// Mocking the repository
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := repository.NewMockRepositoryInterface(ctrl)
	h := helper.NewHelper(helper.NewHelperOptions{
		JwtPrivateKeyPath: "../storage/key.pem",
		JwtPublicKeyPath:  "../storage/key.pem.pub",
	})

	// Test cases
	tests := []struct {
		caseName     string
		payload      string
		token        func() string
		mockFunc     func()
		expectedCode int
	}{
		{
			caseName: "Unauthorized",
			payload:  `{"phoneNumber":"+62123456789","fullName":"test"}`,
			token: func() string {
				return ""
			},
			mockFunc:     func() {},
			expectedCode: http.StatusForbidden,
		},
		{
			caseName: "Positive case",
			payload:  `{"phoneNumber":"+62123456789","fullName":"test"}`,
			token: func() string {
				token := ""
				h.GenerateAccessToken(&token, 1)
				return token
			},
			mockFunc: func() {
				m.
					EXPECT().
					IsPhoneNumberChanged(gomock.Any(), gomock.Any()).
					Return(repository.IsPhoneNumberChangedOutput{IsChanged: false}, nil)
				m.
					EXPECT().
					UpdateUserById(gomock.Any(), gomock.Any()).
					Return(repository.UpdateUserByIdOutput{
						UserId:      1,
						PhoneNumber: "+62123456789",
						FullName:    "test",
					}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			caseName: "Empty payload",
			payload:  "",
			token: func() string {
				token := ""
				h.GenerateAccessToken(&token, 1)
				return token
			},
			mockFunc:     func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			caseName: "Invalid full name",
			payload:  `{"phoneNumber":"+62123456789","fullName":"t"}`,
			token: func() string {
				token := ""
				h.GenerateAccessToken(&token, 1)
				return token
			},
			mockFunc:     func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			caseName: "Update phone number",
			payload:  `{"phoneNumber":"+62123456789","fullName":"test"}`,
			token: func() string {
				token := ""
				h.GenerateAccessToken(&token, 1)
				return token
			},
			mockFunc: func() {
				m.
					EXPECT().
					IsPhoneNumberChanged(gomock.Any(), gomock.Any()).
					Return(repository.IsPhoneNumberChangedOutput{IsChanged: true}, nil)
			},
			expectedCode: http.StatusConflict,
		},
	}

	for _, test := range tests {
		t.Run(test.caseName, func(t *testing.T) {
			// Creating the server
			e := echo.New()
			server := NewServer(NewServerOptions{
				Repository: m,
				Helper:     h,
				Echo:       e,
			})
			generated.RegisterHandlers(e, server)

			req := httptest.NewRequest(
				http.MethodPut,
				"/",
				strings.NewReader(test.payload),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set("Authorization", "Bearer "+test.token())
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			test.mockFunc()

			if err := server.UpdateUser(c); err != nil {
				t.Errorf("Error: %v", err)
			}
			if rec.Code != test.expectedCode {
				t.Errorf("Expected %d, got %d", test.expectedCode, rec.Code)
			}
		})
	}
}
