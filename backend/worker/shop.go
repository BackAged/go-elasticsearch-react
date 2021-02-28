package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/BackAged/go-elasticsearch-react/backend/search"
)

// AddShops creates new shop
func (h *handler) AddShops(shops []byte) error {
	log.Printf("worker.AddShops started with event:%s payloadL:%s\n", RoutingKeyShopCreate, shops)

	var shps []*search.Shop
	if err := json.Unmarshal(shops, &shps); err != nil {
		fmt.Println("worker.AddShops couldn't unmarshal msg payload", err)
		return err
	}

	_, err := h.svc.AddShops(context.Background(), shps)
	if err != nil {
		fmt.Println("worker.AddShops service error:", err)
		return err
	}

	return nil
}

// UpdateShops updates shops
func (h *handler) UpdateShops(shops []byte) error {
	log.Printf("worker.UpdateShops started with event:%s payloadL:%s\n", RoutingKeyShopdUpdate, shops)

	var shps []*search.Shop
	if err := json.Unmarshal(shops, &shps); err != nil {
		fmt.Println("worker.UpdateShops couldn't unmarshal msg payload", err)
		return err
	}

	_, err := h.svc.UpdateShops(context.Background(), shps)
	if err != nil {
		fmt.Println("worker.UpdateShops service error:", err)
		return err
	}

	return nil
}

// DeleteShopReq ...
type DeleteShopReq struct {
	IDS []int64 `json:"ids"`
}

// DeleteShops delete shops
func (h *handler) DeleteShops(shopSlugs []byte) error {
	log.Printf("worker.DeleteShops started with event:%s payloadL:%s\n", RoutingKeyShopDelete, shopSlugs)

	var shpSlgs *DeleteShopReq
	if err := json.Unmarshal(shopSlugs, &shpSlgs); err != nil {
		fmt.Println("worker.DeleteShops couldn't unmarshal msg payload", err)
		return err
	}

	err := h.svc.DeleteShops(context.Background(), shpSlgs.IDS)
	if err != nil {
		fmt.Println("worker.DeleteShops service error:", err)
		return err
	}

	return nil
}
