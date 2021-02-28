package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/BackAged/go-elasticsearch-react/backend/search"
)

// CreateOrder creates new order
func (h *handler) AddBrands(brands []byte) error {
	log.Printf("worker.AddBrands started with event:%s payloadL:%s\n", RoutingKeyBrandCreate, brands)

	var brnds []*search.Brand
	if err := json.Unmarshal(brands, &brnds); err != nil {
		fmt.Println("worker.AddBrands couldn't unmarshal msg payload", err)
		return err
	}

	_, err := h.svc.AddBrands(context.Background(), brnds)
	if err != nil {
		fmt.Println("worker.AddBrands service error:", err)
		return err
	}

	return nil
}

// UpdateBrands updates brands
func (h *handler) UpdateBrands(brands []byte) error {
	log.Printf("worker.UpdateBrands started with event:%s payloadL:%s\n", RoutingKeyBrandUpdate, brands)

	var brnds []*search.Brand
	if err := json.Unmarshal(brands, &brnds); err != nil {
		fmt.Println("worker.UpdateBrands couldn't unmarshal msg payload", err)
		return err
	}

	_, err := h.svc.UpdateBrands(context.Background(), brnds)
	if err != nil {
		fmt.Println("worker.AddBrands service error:", err)
		return err
	}

	return nil
}

// DeleteBrandsReq ...
type DeleteBrandsReq struct {
	IDS []int64 `json:"slugs"`
}

// DeleteBrands creates new order
func (h *handler) DeleteBrands(brandSlugs []byte) error {
	log.Printf("worker.DeleteBrands started with event:%s payloadL:%s\n", RoutingKeyBrandDelete, brandSlugs)

	var brndIDS *DeleteBrandsReq
	if err := json.Unmarshal(brandSlugs, &brndIDS); err != nil {
		fmt.Println("worker.DeleteBrands couldn't unmarshal msg payload", err)
		return err
	}

	err := h.svc.DeleteBrands(context.Background(), brndIDS.IDS)
	if err != nil {
		fmt.Println("worker.DeleteBrands service error:", err)
		return err
	}

	return nil
}
