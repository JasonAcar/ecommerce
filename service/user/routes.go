package user

import (
	"fmt"
	"github.com/JasonAcar/ecommerce/common"
	"github.com/JasonAcar/ecommerce/config"
	"github.com/JasonAcar/ecommerce/service/auth"
	"github.com/JasonAcar/ecommerce/types"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods(http.MethodPost)
	router.HandleFunc("/register", h.handleRegister).Methods(http.MethodPost)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.LoginUserPayload
	if err := common.ParseJSON(r, &payload); err != nil {
		common.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// validate payload
	if err := common.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		common.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check if the user exists
	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		common.WriteError(w, http.StatusBadGateway, fmt.Errorf("not found, invalid email or password"))
		return
	}

	if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
		common.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}
	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, err)
		return
	}

	common.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.RegisterUserPayload
	if err := common.ParseJSON(r, &payload); err != nil {
		common.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// validate payload
	if err := common.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		common.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check if the user exists
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		common.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	hashedPasswd, err := auth.HashPassword(payload.Password)
	if err != nil || hashedPasswd == "" {
		common.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// if it doesnt we create the new user
	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPasswd,
	})
	if err != nil {
		common.WriteError(w, http.StatusInternalServerError, err)
	}
	common.WriteJSON(w, http.StatusCreated, nil)
}
