package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/BackAged/go-elasticsearch-react/backend/search"
)

// AddProducts creates new product
func (h *handler) AddProducts(products []byte) error {
	log.Printf("worker.AddProducts started with event:%s payloadL:%s\n", RoutingKeyProductCreate, products)

	var prds []*search.Product
	if err := json.Unmarshal(products, &prds); err != nil {
		fmt.Println("worker.AddProducts couldn't unmarshal msg payload", err)
		return err
	}

	_, err := h.svc.AddProducts(context.Background(), prds)
	if err != nil {
		fmt.Println("worker.AddProducts service error:", err)
		return err
	}

	return nil
}

// UpdateProducts updates products
func (h *handler) UpdateProducts(products []byte) error {
	log.Printf("worker.UpdateProducts started with event:%s payloadL:%s\n", RoutingKeyProductUpdate, products)

	var prds []*search.Product
	if err := json.Unmarshal(products, &prds); err != nil {
		fmt.Println("worker.UpdateProducts couldn't unmarshal msg payload", err)
		return err
	}

	_, err := h.svc.UpdateProducts(context.Background(), prds)
	if err != nil {
		fmt.Println("worker.UpdateProducts service error:", err)
		return err
	}

	return nil
}

// DeleteProductReq ...
type DeleteProductReq struct {
	ShopItemIDS []int64 `json:"shop_item_ids"`
}

// DeleteProducts delete shops
func (h *handler) DeleteProducts(productSlugs []byte) error {
	log.Printf("worker.DeleteProducts started with event:%s payloadL:%s\n", RoutingKeyProductCreate, productSlugs)

	var prdSlgs *DeleteProductReq
	if err := json.Unmarshal(productSlugs, &prdSlgs); err != nil {
		fmt.Println("worker.DeleteShops couldn't unmarshal msg payload", err)
		return err
	}

	err := h.svc.DeleteProducts(context.Background(), prdSlgs.ShopItemIDS)
	if err != nil {
		fmt.Println("worker.DeleteShops service error:", err)
		return err
	}

	return nil
}
