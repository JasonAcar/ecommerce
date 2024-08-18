package cart

import (
	"fmt"
	"github.com/JasonAcar/ecommerce/common"
	"github.com/JasonAcar/ecommerce/service/auth"
	"github.com/JasonAcar/ecommerce/types"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	store        types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
}

func NewHandler(store types.OrderStore, productStore types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, productStore: productStore, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	var cart types.CartCheckoutPayload
	if err := common.ParseJSON(r, &cart); err != nil {
		common.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := common.Validate.Struct(cart); err != nil {
		errors := err.(validator.ValidationErrors)
		common.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}
	//get products
	productIDs, err := getCartItemsIDs(cart.Items)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, err)
	}
	ps, err := h.productStore.GetProductByIDs(productIDs)
	if err != nil {
		common.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	orderId, totalPrice, err := h.createOrder(ps, cart.Items, userID)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, err)
		return
	}
	common.WriteJSON(w, http.StatusOK, map[string]any{
		"total_price": totalPrice,
		"order_id":    orderId,
	})
}
