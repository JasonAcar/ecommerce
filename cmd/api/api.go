package api

import (
	"database/sql"
	"github.com/JasonAcar/ecommerce/service/cart"
	"github.com/JasonAcar/ecommerce/service/order"
	"github.com/JasonAcar/ecommerce/service/product"
	"github.com/JasonAcar/ecommerce/service/user"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	r := mux.NewRouter()
	sr := r.PathPrefix("/api/v1").Subrouter()
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(sr)

	pStore := product.NewStore(s.db)
	pHandler := product.NewHandler(pStore)
	pHandler.RegisterRoutes(sr)

	orderStore := order.NewStore(s.db)
	cartHandler := cart.NewHandler(orderStore, pStore, userStore)
	cartHandler.RegisterRoutes(sr)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, r)
}
