package search

// Sort defines sort criteria
type Sort struct {
	FieldName string `json:"field_name"`
	Order     string `json:"order"`
}

// FacetSearchReq defines dto of facet search
type FacetSearchReq struct {
	Term            string
	From            int64
	Size            int64
	BucketSize      int
	CategoryFilters []string
	BrandFilters    []string
	ShopFilters     []string
	ColorFilters    []string
	Sort            []*Sort
}

// Bucket ...
type Bucket struct {
	Key      string `json:"key" mapstructure:"key"`
	DocCount int32  `json:"doc_count" mapstructure:"doc_count"`
}

// FacetRes facet response
type FacetRes struct {
	Brands     []Bucket `json:"brands"`
	Categories []Bucket `json:"categories"`
	Shops      []Bucket `json:"shops"`
	Colors     []Bucket `json:"colors"`
}
