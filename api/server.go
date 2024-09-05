package api

import (
	"log"
	"net/http"

	"github.com/anirudhchy/gobank/storage"
	"github.com/anirudhchy/gobank/utils"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      storage.Storage
}

func NewAPIServer(listenAddr string, store storage.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/login", utils.MakeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/account", utils.MakeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", withJWTAuth(utils.MakeHTTPHandleFunc(s.handleGetAccountByID), s.store))
	router.HandleFunc("/transfer", utils.MakeHTTPHandleFunc(s.handleTransfer))

	log.Println("JSON API server running on port", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)

}
