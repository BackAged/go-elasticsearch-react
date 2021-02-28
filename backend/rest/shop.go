package rest

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/BackAged/go-elasticsearch-react/backend/search"
	"github.com/go-chi/chi"
)

// ShopHandler defines shop handler
type ShopHandler struct {
	svc search.Service
}

// NewShopHandler ...
func NewShopHandler(svc search.Service) *ShopHandler {
	return &ShopHandler{
		svc: svc,
	}
}

// Router ..
func (h *ShopHandler) Router() http.Handler {
	router := chi.NewRouter()

	router.Get("/", h.SearchAsYouTypeShop)
	router.Post("/bulk-insert", h.AddShops)
	router.Post("/bulk-delete", h.DeleteShops)
	router.Post("/bulk-update", h.UpdateShops)

	return router
}

// SearchAsYouTypeShop ...
func (h *ShopHandler) SearchAsYouTypeShop(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")
	pager := getPager(r)
	skip := (pager.Page - 1) * pager.Limit
	limit := pager.Limit

	brnds, total, err := h.svc.SearchShopAsType(r.Context(), term, skip, limit)
	if err != nil {
		ServeJSON(w, "", http.StatusOK, "Successful", nil, nil, nil)
		return
	}

	ServeJSON(w, "", http.StatusOK, "Successful", brnds, &total, nil)
	return
}

// AddShops ...
func (h *ShopHandler) AddShops(w http.ResponseWriter, r *http.Request) {
	brnds := []*search.Shop{}
	err := json.NewDecoder(r.Body).Decode(&brnds)
	if err != nil {
		log.Println("shopHandler.AddShops =>  invalid request body: ", err)
		ServeJSON(w, "E_INVALID_ARG", http.StatusBadRequest, "invalid request body", nil, nil, nil)
		return
	}
	if len(brnds) == 0 {
		log.Println("shopHandler.AddShops =>  invalid request body: ", err)
		ServeJSON(w, "E_INVALID_ARG", http.StatusBadRequest, "no brands to insert in request body", nil, nil, nil)
		return
	}

	brnds, err = h.svc.AddShops(r.Context(), brnds)
	if err != nil {
		ServeJSON(w, "", http.StatusInternalServerError, "Something went wrong!", nil, nil, nil)
		return
	}

	ServeJSON(w, "", http.StatusOK, "Successful", brnds, nil, nil)
	return
}

// DeleteShopReq ...
type DeleteShopReq struct {
	IDS []int64 `json:"ids"`
}

// DeleteShops ...
func (h *ShopHandler) DeleteShops(w http.ResponseWriter, r *http.Request) {
	dbrq := &DeleteShopReq{}
	err := json.NewDecoder(r.Body).Decode(&dbrq)
	if err != nil {
		log.Println("shopHandler.DeleteShops =>  invalid request body: ", err)
		ServeJSON(w, "E_INVALID_ARG", http.StatusBadRequest, "invalid request body", nil, nil, nil)
		return
	}
	if len(dbrq.IDS) == 0 {
		log.Println("shopHandler.DeleteShops =>  invalid request body: ", err)
		ServeJSON(w, "E_INVALID_ARG", http.StatusBadRequest, "no brands slug to delete in request body", nil, nil, nil)
		return
	}

	err = h.svc.DeleteShops(r.Context(), dbrq.IDS)
	if err != nil {
		ServeJSON(w, "", http.StatusInternalServerError, "Something went wrong!", nil, nil, nil)
		return
	}

	ServeJSON(w, "", http.StatusOK, "Successful", nil, nil, nil)
	return
}

// UpdateShops ...
func (h *ShopHandler) UpdateShops(w http.ResponseWriter, r *http.Request) {
	shps := []*search.Shop{}
	err := json.NewDecoder(r.Body).Decode(&shps)
	if err != nil {
		log.Println("shopHandler.UpdateShops =>  invalid request body: ", err)
		ServeJSON(w, "E_INVALID_ARG", http.StatusBadRequest, "invalid request body", nil, nil, nil)
		return
	}
	if len(shps) == 0 {
		log.Println("shopHandler.UpdateShops =>  invalid request body: ", err)
		ServeJSON(w, "E_INVALID_ARG", http.StatusBadRequest, "no brands to insert in request body", nil, nil, nil)
		return
	}

	shps, err = h.svc.UpdateShops(r.Context(), shps)
	if err != nil {
		ServeJSON(w, "", http.StatusInternalServerError, "Something went wrong!", nil, nil, nil)
		return
	}

	ServeJSON(w, "", http.StatusOK, "Successful", shps, nil, nil)
	return
}
