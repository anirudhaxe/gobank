package api

import (
	"fmt"
	"net/http"

	"github.com/anirudhchy/gobank/storage"
	"github.com/anirudhchy/gobank/types"
	"github.com/anirudhchy/gobank/utils"
	jwt "github.com/golang-jwt/jwt/v5"
)

func permissionDenied(w http.ResponseWriter) {

	utils.WriteJSON(w, http.StatusForbidden, types.ApiError{Error: "permission denied"})
}

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50TnVtYmVyIjo5MTQ2MzksImV4cGlyZXNBdCI6MTUwMDB9.Sn6T1Dk6xHogu8CNSO6mUrj0GgmNqpd_RDNUL2R7j-Q

func withJWTAuth(handlerFunc http.HandlerFunc, s storage.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling JWT auth middleware")

		tokenString := r.Header.Get("x-jwt-token")

		token, err := validateJWT(tokenString)

		if err != nil {
			permissionDenied(w)
			return
		}

		if !token.Valid {
			permissionDenied(w)
			return
		}

		userID, err := utils.GetID(r)

		if err != nil {
			permissionDenied(w)
			return
		}

		account, err := s.GetAccountByID(userID)

		if err != nil {
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		// panic(reflect.TypeOf(claims["accountNumber"]))

		if account.Number != int64(claims["accountNumber"].(float64)) {
			permissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}
