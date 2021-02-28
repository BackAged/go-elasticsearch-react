package worker

// routing key const
const (
	RoutingKeyBrandCreate   = "catalog.brand.create"
	RoutingKeyBrandDelete   = "catalog.brand.delete"
	RoutingKeyBrandUpdate   = "catalog.brand.update"
	RoutingKeyShopCreate    = "catalog.shop.create"
	RoutingKeyShopDelete    = "catalog.shop.delete"
	RoutingKeyShopdUpdate   = "catalog.shop.update"
	RoutingKeyProductCreate = "catalog.product.create"
	RoutingKeyProductDelete = "catalog.product.delete"
	RoutingKeyProductUpdate = "catalog.product.update"
)

// RoutingKeys ...
func RoutingKeys() []string {
	return []string{
		RoutingKeyBrandCreate, RoutingKeyBrandDelete, RoutingKeyBrandUpdate,
		RoutingKeyProductCreate, RoutingKeyProductDelete, RoutingKeyProductUpdate,
		RoutingKeyShopCreate, RoutingKeyShopDelete, RoutingKeyShopdUpdate,
	}
}

// SetUpRouter ...
func (w *Worker) SetUpRouter(h Handler) error {
	if err := w.RegisterTask(RoutingKeyBrandCreate, h.AddBrands); err != nil {
		return err
	}
	if err := w.RegisterTask(RoutingKeyBrandDelete, h.DeleteBrands); err != nil {
		return err
	}
	if err := w.RegisterTask(RoutingKeyBrandUpdate, h.UpdateBrands); err != nil {
		return err
	}
	if err := w.RegisterTask(RoutingKeyShopCreate, h.AddShops); err != nil {
		return err
	}
	if err := w.RegisterTask(RoutingKeyShopDelete, h.DeleteShops); err != nil {
		return err
	}
	if err := w.RegisterTask(RoutingKeyShopdUpdate, h.UpdateShops); err != nil {
		return err
	}
	if err := w.RegisterTask(RoutingKeyProductCreate, h.AddProducts); err != nil {
		return err
	}
	if err := w.RegisterTask(RoutingKeyProductDelete, h.DeleteProducts); err != nil {
		return err
	}
	if err := w.RegisterTask(RoutingKeyProductUpdate, h.UpdateProducts); err != nil {
		return err
	}

	return nil
}
