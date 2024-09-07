package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anirudhchy/gobank/types"
	"github.com/anirudhchy/gobank/utils"
)

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	var req types.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	acc, err := s.store.GetAccountByNumber(int(req.Number))

	if err != nil {
		return err
	}

	// fmt.Printf("%+v\n", acc)

	if !acc.ValidatePassword(req.Password) {
		return fmt.Errorf("not authenticated")
	}

	token, err := createJWT(acc)

	if err != nil {
		return err
	}

	resp := types.LoginResponse{
		Number:    acc.Number,
		FirstName: acc.FirstName,
		Lastname:  acc.LastName,
		Token:     token,
	}

	return utils.WriteJSON(w, http.StatusOK, resp)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {

	if r.Method == "GET" {
		return s.handleGetAccount(w)
	}

	if r.Method == "POST" {
		return s.handleRegisterAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter) error {

	accounts, err := s.store.GetAccounts()

	if err != nil {
		return err
	}

	responseAccounts := make([]*types.GetAccountsResponse, 0, len(accounts))

	for _, account := range accounts {

		responseAccount := &types.GetAccountsResponse{
			FirstName: account.FirstName,
			LastName:  account.LastName,
			Number:    account.Number,
		}

		responseAccounts = append(responseAccounts, responseAccount)
	}

	return utils.WriteJSON(w, http.StatusOK, responseAccounts)
}

func (s *APIServer) handleRegisterAccount(w http.ResponseWriter, r *http.Request) error {

	req := new(types.RegisterAccountRequest)
	// createAccountReq := CreateAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	account, err := types.NewAccount(req.FirstName, req.LastName, req.Password)
	if err != nil {
		return err
	}

	if err := s.store.CreateAccount(account); err != nil {
		return nil
	}

	tokenString, err := createJWT(account)

	if err != nil {
		return err
	}

	resp := types.RegisterAccountResponse{
		Number:    account.Number,
		FirstName: account.FirstName,
		Lastname:  account.LastName,
		Token:     tokenString,
	}

	return utils.WriteJSON(w, http.StatusOK, resp)
}

func (s *APIServer) handleAccountByNumber(w http.ResponseWriter, r *http.Request) error {

	if r.Method == "GET" {

		return s.handleGetAccountByNumber(w, r)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccountByNumber(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAccountByNumber(w http.ResponseWriter, r *http.Request) error {

	accountNumber, err := utils.GetAccountNumber(r)

	if err != nil {
		return err
	}

	account, err := s.store.GetAccountByNumber(accountNumber)

	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccountByNumber(w http.ResponseWriter, r *http.Request) error {

	accountNumber, err := utils.GetAccountNumber(r)

	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(accountNumber); err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, map[string]int{"deleted": accountNumber})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {

	transferReq := new(types.TransferRequest)

	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}

	// fmt.Println(transferReq.ToAccountNumber,reflect.TypeOf(transferReq.ToAccountNumber))

	accountNumber, err := utils.GetAccountNumber(r)

	if err != nil {
		return err
	}

	fromAccount, err := s.store.GetAccountByNumber(accountNumber)

	if err != nil {
		return err
	}

	if fromAccount.Number == transferReq.ToAccountNumber {
		return fmt.Errorf("to account number not valid")
	}

	if fromAccount.Balance < transferReq.Amount {
		return fmt.Errorf("you dont have enough balance to complete the transaction")
	}

	toAccount, err := s.store.GetAccountByNumber(int(transferReq.ToAccountNumber))

	if err != nil {
		return err
	}

	// todo: refactor following operation so that these two queries are completed in a transaction
	err = s.store.UpdateAccountBalanceByNumber(fromAccount.Balance-transferReq.Amount, fromAccount.Number)

	if err != nil {
		return err
	}

	err = s.store.UpdateAccountBalanceByNumber(toAccount.Balance+transferReq.Amount, toAccount.Number)

	if err != nil {
		return err
	}

	defer r.Body.Close()

	return utils.WriteJSON(w, http.StatusOK, "amount transferred")
}
