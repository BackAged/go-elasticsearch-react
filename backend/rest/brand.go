package rest

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/BackAged/go-elasticsearch-react/backend/search"
	"github.com/go-chi/chi"
)

// BrandHandler defines brand handler
type BrandHandler struct {
	svc search.Service
}

// NewBrandHandler ...
func NewBrandHandler(svc search.Service) *BrandHandler {
	return &BrandHandler{
		svc: svc,
	}
}

// Router ..
func (h *BrandHandler) Router() http.Handler {
	router := chi.NewRouter()

	router.Get("/", h.SearchAsYouTypeBrand)
	router.Post("/bulk-insert", h.AddBrands)
	router.Post("/bulk-delete", h.DeleteBrands)
	router.Post("/bulk-update", h.UpdateBrands)

	return router
}

// SearchAsYouTypeBrand ...
func (h *BrandHandler) SearchAsYouTypeBrand(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")
	pager := getPager(r)
	skip := (pager.Page - 1) * pager.Limit
	limit := pager.Limit

	brnds, total, err := h.svc.SearchBrandAsType(r.Context(), term, skip, limit)
	if err != nil {
		ServeJSON(w, "", http.StatusInternalServerError, "Something went wrong!", nil, nil, nil)
		return
	}

	ServeJSON(w, "", http.StatusOK, "Successful", brnds, &total, nil)
	return
}

// AddBrands ...
func (h *BrandHandler) AddBrands(w http.ResponseWriter, r *http.Request) {
	brnds := []*search.Brand{}
	err := json.NewDecoder(r.Body).Decode(&brnds)
	if err != nil {
		log.Println("brandHandler.AddBrands =>  invalid request body: ", err)
		ServeJSON(w, "E_INVALID_ARG", http.StatusBadRequest, "invalid request body", nil, nil, nil)
		return
	}
	if len(brnds) == 0 {
		log.Println("brandHandler.AddBrands =>  invalid request body: ", err)
		ServeJSON(w, "E_INVALID_ARG", http.StatusBadRequest, "no brands to insert in request body", nil, nil, nil)
		return
	}

	brnds, err = h.svc.AddBrands(r.Context(), brnds)
	if err != nil {
		ServeJSON(w, "", http.StatusInternalServerError, "Something went wrong!", nil, nil, nil)
		return
	}

	ServeJSON(w, "", http.StatusOK, "Successful", brnds, nil, nil)
	return
}

// UpdateBrands ...
func (h *BrandHandler) UpdateBrands(w http.ResponseWriter, r *http.Request) {
	brnds := []*search.Brand{}
	err := json.NewDecoder(r.Body).Decode(&brnds)
	if err != nil {
		log.Println("brandHandler.UpdateBrands =>  invalid request body: ", err)
		ServeJSON(w, "E_INVALID_ARG", http.StatusBadRequest, "invalid request body", nil, nil, nil)
		return
	}
	if len(brnds) == 0 {
		log.Println("brandHandler.UpdateBrands =>  invalid request body: ", err)
		ServeJSON(w, "E_INVALID_ARG", http.StatusBadRequest, "no brands to insert in request body", nil, nil, nil)
		return
	}

	brnds, err = h.svc.UpdateBrands(r.Context(), brnds)
	if err != nil {
		ServeJSON(w, "", http.StatusInternalServerError, "Something went wrong!", nil, nil, nil)
		return
	}

	ServeJSON(w, "", http.StatusOK, "Successful", brnds, nil, nil)
	return
}

// DeleteBrandsReq ...
type DeleteBrandsReq struct {
	IDS []int64 `json:"ids"`
}

// DeleteBrands ...
func (h *BrandHandler) DeleteBrands(w http.ResponseWriter, r *http.Request) {
	dbrq := &DeleteBrandsReq{}
	err := json.NewDecoder(r.Body).Decode(&dbrq)
	if err != nil {
		log.Println("brandHandler.DeleteBrands =>  invalid request body: ", err)
		ServeJSON(w, "E_INVALID_ARG", http.StatusBadRequest, "invalid request body", nil, nil, nil)
		return
	}
	if len(dbrq.IDS) == 0 {
		log.Println("brandHandler.DeleteBrands =>  invalid request body: ", err)
		ServeJSON(w, "E_INVALID_ARG", http.StatusBadRequest, "no brands slug to delete in request body", nil, nil, nil)
		return
	}

	err = h.svc.DeleteBrands(r.Context(), dbrq.IDS)
	if err != nil {
		ServeJSON(w, "", http.StatusInternalServerError, "Something went wrong!", nil, nil, nil)
		return
	}

	ServeJSON(w, "", http.StatusOK, "Successful", nil, nil, nil)
	return
}
