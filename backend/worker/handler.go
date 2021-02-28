package worker

import (
	"github.com/BackAged/go-elasticsearch-react/backend/search"
)

// Handler defines handler interface
type Handler interface {
	AddBrands(brands []byte) error
	DeleteBrands(brandSlugs []byte) error
	UpdateBrands(brands []byte) error
	AddShops(shops []byte) error
	DeleteShops(shopSlugs []byte) error
	UpdateShops(shops []byte) error
	AddProducts(products []byte) error
	DeleteProducts(productSlugs []byte) error
	UpdateProducts(products []byte) error
}

type handler struct {
	svc search.Service
}

// NewHandler instantiate a new handler
func NewHandler(svc search.Service) Handler {
	return &handler{
		svc: svc,
	}
}
