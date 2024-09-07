package types

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

type LoginRequest struct {
	Number   int64  `json:"number"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Number    int64  `json:"number"`
	FirstName string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Token     string `json:"token"`
}

type GetAccountsResponse struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Number    int64  `json:"number"`
}

type TransferRequest struct {
	ToAccountNumber int64 `json:"toAccountNumber"`
	Amount          int64 `json:"amount"`
}

type RegisterAccountRequest struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Password  string `json:"password"`
}

type RegisterAccountResponse struct {
	Number    int64  `json:"number"`
	FirstName string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Token     string `json:"token"`
}

type Account struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	Number            int64     `json:"number"`
	EncryptedPassword string    `json:"-"`
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"createdAt"`
}

type ApiError struct {
	Error string `json:"error"`
}

func (a *Account) ValidatePassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(pw)) == nil
}

func NewAccount(firstName, lastName, password string) (*Account, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Account{
		// ID:        rand.Intn(10000),
		FirstName:         firstName,
		LastName:          lastName,
		EncryptedPassword: string(encpw),
		Number:            int64(rand.Intn(1000000)),
		// Balance will be initialised to zero value of int64 which is 0
		CreatedAt: time.Now().UTC(),
	}, nil
}
