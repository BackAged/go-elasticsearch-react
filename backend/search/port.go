package search

import "context"

// ProductRepo defines interface for infra
type ProductRepo interface {
	BulkInsert(context.Context, []*Product) ([]*Product, error)
	Add(context.Context, *Product) (*Product, error)
	Search(ctx context.Context, term string, skip int64, limit int64) ([]*Product, int64, error)
	SearchFacet(context.Context, FacetSearchReq) ([]*Product, *FacetRes, int64, error)
	UpdateProductScore(context.Context, int64) error
	DeleteMany(ctx context.Context, shopItemIDS []int64) error
	UpdateMany(ctx context.Context, products []*Product) ([]*Product, error)
}

// BrandRepo defines interface for infra
type BrandRepo interface {
	SearchAsType(ctx context.Context, term string, skip int64, limit int64) ([]*Brand, int64, error)
	BulkInsert(ctx context.Context, brands []*Brand) ([]*Brand, error)
	DeleteMany(ctx context.Context, brandIDS []int64) error
	UpdateMany(ctx context.Context, brands []*Brand) ([]*Brand, error)
}

// ShopRepo defines interface for infra
type ShopRepo interface {
	SearchAsType(ctx context.Context, term string, skip int64, limit int64) ([]*Shop, int64, error)
	BulkInsert(ctx context.Context, shops []*Shop) ([]*Shop, error)
	DeleteMany(ctx context.Context, brandIDS []int64) error
	UpdateMany(ctx context.Context, shops []*Shop) ([]*Shop, error)
}

// Service provides port for application adapter.
type Service interface {
	AddProduct(context.Context, *Product) (*Product, error)
	AddProducts(context.Context, []*Product) ([]*Product, error)
	SearchProductAsType(ctx context.Context, term string, skip int64, limit int64) ([]*Product, int64, error)
	DeleteProducts(ctx context.Context, shopItemIDS []int64) error
	FacetSearchProducts(context.Context, FacetSearchReq) ([]*Product, *FacetRes, int64, error)
	UpdateProductScore(context.Context, int64) error
	UpdateProducts(ctx context.Context, products []*Product) ([]*Product, error)

	SearchShopAsType(ctx context.Context, term string, skip int64, limit int64) ([]*Shop, int64, error)
	AddShops(ctx context.Context, shops []*Shop) ([]*Shop, error)
	DeleteShops(ctx context.Context, shopIDS []int64) error
	UpdateShops(ctx context.Context, shops []*Shop) ([]*Shop, error)

	SearchBrandAsType(ctx context.Context, term string, skip int64, limit int64) ([]*Brand, int64, error)
	AddBrands(ctx context.Context, brands []*Brand) ([]*Brand, error)
	DeleteBrands(ctx context.Context, brandIDS []int64) error
	UpdateBrands(ctx context.Context, brands []*Brand) ([]*Brand, error)
}
