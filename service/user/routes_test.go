package user

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/JasonAcar/ecommerce/types"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockUserStore struct {
	store *sql.DB
}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, fmt.Errorf("user not found")
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) CreateUser(types.User) error {
	return nil
}

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}

	t.Run("should fail if the user payload is invalid", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "user",
			LastName:  "123",
			Email:     "invalid",
			Password:  "asdf",
		}

		m, _ := json.Marshal(payload)

		handler := NewHandler(userStore)
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(m))
		//req, err := http.NewRequest(http.MethodPost, "/register", nil) // This will PASS the test because the request fails due to empty request body
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("unexpected status code %d got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should correctly register the user", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "user",
			LastName:  "123",
			Email:     "valid@mail.com",
			Password:  "asdf",
		}

		m, _ := json.Marshal(payload)

		handler := NewHandler(userStore)
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(m))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated {
			t.Errorf("unexpected status code %d got %d", http.StatusCreated, rr.Code)
		}
	})
}
