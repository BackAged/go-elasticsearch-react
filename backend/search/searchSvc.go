package search

import (
	"context"
	"fmt"
)

type service struct {
	prdRepo  ProductRepo
	brndRepo BrandRepo
	shpRepo  ShopRepo
}

// NewService creates a service with the necessary dependencies.
func NewService(prdRepo ProductRepo, brndRepo BrandRepo, shpRepo ShopRepo) Service {
	return &service{
		prdRepo:  prdRepo,
		brndRepo: brndRepo,
		shpRepo:  shpRepo,
	}
}

/////////////////// Product //////////////////
func (s *service) AddProduct(ctx context.Context, product *Product) (*Product, error) {
	prd, err := s.prdRepo.Add(context.Background(), product)
	fmt.Print(prd)
	return prd, err
}

func (s *service) AddProducts(ctx context.Context, products []*Product) ([]*Product, error) {
	prds, err := s.prdRepo.BulkInsert(context.Background(), products)
	fmt.Print(prds)
	return prds, err
}

func (s *service) SearchProductAsType(ctx context.Context, term string, skip int64, limit int64) ([]*Product, int64, error) {
	prds, total, err := s.prdRepo.Search(context.Background(), term, skip, limit)
	if err != nil {
		return nil, 0, err
	}

	return prds, total, err
}

func (s *service) DeleteProducts(ctx context.Context, shopItemIDS []int64) error {
	return s.prdRepo.DeleteMany(context.Background(), shopItemIDS)
}

func (s *service) FacetSearchProducts(ctx context.Context, req FacetSearchReq) ([]*Product, *FacetRes, int64, error) {
	return s.prdRepo.SearchFacet(context.Background(), req)
}

func (s *service) UpdateProductScore(ctx context.Context, shopItemID int64) error {
	return s.prdRepo.UpdateProductScore(context.Background(), shopItemID)
}

func (s *service) UpdateProducts(ctx context.Context, products []*Product) ([]*Product, error) {
	prds, err := s.prdRepo.UpdateMany(context.Background(), products)
	if err != nil {
		return nil, err
	}

	return prds, err
}

/////////////////// Brand //////////////////
func (s *service) SearchBrandAsType(ctx context.Context, term string, skip int64, limit int64) ([]*Brand, int64, error) {
	brnds, total, err := s.brndRepo.SearchAsType(context.Background(), term, skip, limit)
	if err != nil {
		return nil, 0, err
	}

	return brnds, total, err
}

func (s *service) AddBrands(ctx context.Context, brands []*Brand) ([]*Brand, error) {
	brnds, err := s.brndRepo.BulkInsert(context.Background(), brands)
	if err != nil {
		return nil, err
	}

	return brnds, err
}

func (s *service) UpdateBrands(ctx context.Context, brands []*Brand) ([]*Brand, error) {
	brnds, err := s.brndRepo.UpdateMany(context.Background(), brands)
	if err != nil {
		return nil, err
	}

	return brnds, err
}

func (s *service) DeleteBrands(ctx context.Context, brandIDS []int64) error {
	return s.brndRepo.DeleteMany(context.Background(), brandIDS)
}

/////////////////// Shop //////////////////
func (s *service) SearchShopAsType(ctx context.Context, term string, skip int64, limit int64) ([]*Shop, int64, error) {
	shps, total, err := s.shpRepo.SearchAsType(context.Background(), term, skip, limit)
	if err != nil {
		return nil, 0, err
	}

	return shps, total, err
}

func (s *service) AddShops(ctx context.Context, shops []*Shop) ([]*Shop, error) {
	brnds, err := s.shpRepo.BulkInsert(context.Background(), shops)
	if err != nil {
		return nil, err
	}

	return brnds, err
}

func (s *service) DeleteShops(ctx context.Context, shopIDS []int64) error {
	return s.shpRepo.DeleteMany(context.Background(), shopIDS)
}

func (s *service) UpdateShops(ctx context.Context, shops []*Shop) ([]*Shop, error) {
	shps, err := s.shpRepo.UpdateMany(context.Background(), shops)
	if err != nil {
		return nil, err
	}

	return shps, err
}
