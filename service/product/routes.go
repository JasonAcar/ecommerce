package product

import (
	"fmt"
	"github.com/JasonAcar/ecommerce/common"
	"github.com/JasonAcar/ecommerce/types"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.handleProducts).Methods(http.MethodGet)
	router.HandleFunc("/products", h.handleProducts).Methods(http.MethodPost)
}

func (h *Handler) handleProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		products, err := h.store.GetProducts()
		if err != nil {
			common.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		common.WriteJSON(w, http.StatusOK, products)
		return
	}
	if r.Method == http.MethodPost {
		var payload types.CreateProductPayload
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
		newId, err := h.store.AddProduct(types.Product{
			Name:        payload.Name,
			Description: payload.Description,
			Image:       payload.Image,
			Price:       payload.Price,
			Quantity:    payload.Quantity,
		})
		if err != nil {
			common.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		if newId == 0 {
			common.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error creating new product"))
			return
		}

		common.WriteJSON(w, http.StatusCreated, map[string]int64{"productId": newId})
		return
	}
	common.WriteError(w, http.StatusUnprocessableEntity, fmt.Errorf("method not allowed"))
}
