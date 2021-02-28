package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/BackAged/go-elasticsearch-react/backend/search"
	"github.com/go-chi/chi"
)

// ProductHandler defines product handler
type ProductHandler struct {
	svc search.Service
}

// NewProductHandler ...
func NewProductHandler(svc search.Service) *ProductHandler {
	return &ProductHandler{
		svc: svc,
	}
}

// Router ..
func (h *ProductHandler) Router() http.Handler {
	router := chi.NewRouter()

	router.Get("/", h.SearchAsYouTypeProduct)
	router.Post("/bulk-insert", h.AddProducts)
	router.Post("/bulk-delete", h.DeleteProducts)
	router.Post("/bulk-update", h.UpdateProducts)
	router.Post("/search", h.SearchFacet)
	router.Post("/update-score", h.UpdateProductScore)

	return router
}

// SearchAsYouTypeProduct ...
func (h *ProductHandler) SearchAsYouTypeProduct(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")
	pager := getPager(r)
	skip := (pager.Page - 1) * pager.Limit
	limit := pager.Limit

	prds, total, err := h.svc.SearchProductAsType(r.Context(), term, skip, limit)
	if err != nil {
		log.Println("productHandler.SearchAsYouTypeProduct =>  service error: ", err)
		ServeJSON(w, "", http.StatusInternalServerError, "Something went wrong!", nil, nil, nil)
		return
	}

	ServeJSON(w, "", http.StatusOK, "Successful", prds, &total, nil)
	return
}

// AddProducts ...
func (h *ProductHandler) AddProducts(w http.ResponseWriter, r *http.Request) {
	prds := []*search.Product{}
	err := json.NewDecoder(r.Body).Decode(&prds)
	if err != nil {
		log.Println("productHandler.AddProducts =>  invalid request body: ", err)
		ServeJSON(w, "E_INVALID_ARG", http.StatusBadRequest, "invalid request body", nil, nil, nil)
		return
	}
	if len(prds) == 0 {
		log.Println("productHandler.AddProducts =>  invalid request body: ", err)
		ServeJSON(w, "E_INVALID_ARG", http.StatusBadRequest, "no brands to insert in request body", nil, nil, nil)
		return
	}

	prds, err = h.svc.AddProducts(r.Context(), prds)
	if err != nil {
		log.Println("productHandler.AddProducts =>  service error: ", err)
		ServeJSON(w, "", http.StatusInternalServerError, "Something went wrong!", nil, nil, nil)
		return
	}

	ServeJSON(w, "", http.StatusOK, "Successful", prds, nil, nil)
	return
}

// DeleteProductReq ...
type DeleteProductReq struct {
	ShopItemIDS []int64 `json:"shop_item_ids"`
}

// DeleteProducts ...
func (h *ProductHandler) DeleteProducts(w http.ResponseWriter, r *http.Request) {
	dbrq := &DeleteProductReq{}
	err := json.NewDecoder(r.Body).Decode(&dbrq)
	if err != nil {
		log.Println("productHandler.DeleteProducts =>  invalid request body: ", err)
		ServeJSON(w, "E_INVALID_ARG", http.StatusBadRequest, "invalid request body", nil, nil, nil)
		return
	}
	if len(dbrq.ShopItemIDS) == 0 {
		log.Println("productHandler.DeleteProducts =>  invalid request body: ", err)
		ServeJSON(w, "E_INVALID_ARG", http.StatusBadRequest, "no product slug to delete in request body", nil, nil, nil)
		return
	}

	err = h.svc.DeleteProducts(r.Context(), dbrq.ShopItemIDS)
	if err != nil {
		log.Println("productHandler.DeleteProducts =>  service error: ", err)
		ServeJSON(w, "", http.StatusInternalServerError, "Something went wrong!", nil, nil, nil)
		return
	}

	ServeJSON(w, "", http.StatusOK, "Successful", nil, nil, nil)
	return
}

// UpdateProducts ...
func (h *ProductHandler) UpdateProducts(w http.ResponseWriter, r *http.Request) {
	prds := []*search.Product{}
	err := json.NewDecoder(r.Body).Decode(&prds)
	if err != nil {
		log.Println("productHandler.UpdateProducts =>  invalid request body: ", err)
		ServeJSON(w, "E_INVALID_ARG", http.StatusBadRequest, "invalid request body", nil, nil, nil)
		return
	}
	if len(prds) == 0 {
		log.Println("productHandler.UpdateProducts =>  invalid request body: ", err)
		ServeJSON(w, "E_INVALID_ARG", http.StatusBadRequest, "no brands to insert in request body", nil, nil, nil)
		return
	}

	prds, err = h.svc.UpdateProducts(r.Context(), prds)
	if err != nil {
		ServeJSON(w, "", http.StatusInternalServerError, "Something went wrong!", nil, nil, nil)
		return
	}

	ServeJSON(w, "", http.StatusOK, "Successful", prds, nil, nil)
	return
}

// UpdateProductScore ...
func (h *ProductHandler) UpdateProductScore(w http.ResponseWriter, r *http.Request) {
	shpItmID := r.URL.Query().Get("shop_item_id")
	shpItmIDINT, err := strconv.ParseInt(shpItmID, 10, 64)
	if err != nil {
		log.Println("productHandler.UpdateProductScore =>  invalid shop item id: ", err)
		ServeJSON(w, "", http.StatusBadRequest, "invalid shop item id", nil, nil, nil)
		return
	}

	err = h.svc.UpdateProductScore(r.Context(), shpItmIDINT)
	if err != nil {
		log.Println("productHandler.UpdateProductScore =>  service error: ", err)
		ServeJSON(w, "", http.StatusInternalServerError, "Something went wrong!", nil, nil, nil)
		return
	}

	ServeJSON(w, "", http.StatusOK, "Successful", nil, nil, nil)
	return
}

type reqFacetSearchProduct struct {
	Term            string         `json:"term"`
	BrandFilters    []string       `json:"brand_filters"`
	ShopFilters     []string       `json:"shop_filters"`
	CategoryFilters []string       `json:"category_filters"`
	ColorFilters    []string       `json:"color_filters"`
	BucketSize      int            `json:"bucket_size"`
	Sort            []*search.Sort `json:"sort"`
}

type searchFacetRes struct {
	Products []*search.Product `json:"products"`
	Facet    *search.FacetRes  `json:"facets"`
}

// SearchFacet ...
func (h *ProductHandler) SearchFacet(w http.ResponseWriter, r *http.Request) {
	var rs reqFacetSearchProduct
	err := json.NewDecoder(r.Body).Decode(&rs)
	if err != nil {
		ServeJSON(w, "E_INVALID_ARG", http.StatusBadRequest, "Invalid input", nil, nil, nil)
		return
	}

	pager := getPager(r)
	skip := (pager.Page - 1) * pager.Limit
	limit := pager.Limit

	req := search.FacetSearchReq{
		Term:            rs.Term,
		BrandFilters:    rs.BrandFilters,
		BucketSize:      int(rs.BucketSize),
		CategoryFilters: rs.CategoryFilters,
		ColorFilters:    rs.ColorFilters,
		From:            skip,
		ShopFilters:     rs.ShopFilters,
		Size:            limit,
		Sort:            rs.Sort,
	}

	prods, fcts, total, err := h.svc.FacetSearchProducts(r.Context(), req)
	if err != nil {
		log.Println("productHandler.SearchFacet =>  service error: ", err)
		ServeJSON(w, "", http.StatusInternalServerError, "Something went wrong!", nil, nil, nil)
		return
	}

	res := &searchFacetRes{
		Products: prods,
		Facet:    fcts,
	}

	ServeJSON(w, "", http.StatusOK, "Successful", res, &total, nil)
	return

}
