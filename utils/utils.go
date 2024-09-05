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

func GetID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)

	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}

	return id, nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
