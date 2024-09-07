package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/anirudhchy/gobank/types"
	"github.com/gorilla/mux"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

func MakeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, types.ApiError{Error: err.Error()})
		}
	}
}

func GetAccountNumber(r *http.Request) (int, error) {
	accNumberStr := mux.Vars(r)["number"]

	accountNumber, err := strconv.Atoi(accNumberStr)

	if err != nil {
		return accountNumber, fmt.Errorf("invalid account number given %s", accNumberStr)
	}

	return accountNumber, nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
