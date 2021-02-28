package search

// Product defines product type
type Product struct {
	ID              int64         `json:"id" mapstructure:"id"`
	Version         int64         `json:"version" mapstructure:"version"`
	Slug            string        `json:"slug" mapstructure:"slug"`
	Name            string        `json:"name" mapstructure:"name"`
	ShopName        string        `json:"shop_name" mapstructure:"shop_name"`
	ShopSlug        string        `json:"shop_slug" mapstructure:"shop_slug"`
	ShopItemID      int64         `json:"shop_item_id,omitempty" mapstructure:"shop_item_id"`
	Price           float64       `json:"price,omitempty" mapstructure:"price"`
	DiscountedPrice float64       `json:"discounted_price,omitempty" mapstructure:"discounted_price"`
	MinPrice        float64       `json:"min_price,omitempty" mapstructure:"min_price"`
	MaxPrice        float64       `json:"max_price,omitempty" mapstructure:"max_price"`
	BrandName       string        `json:"brand_name,omitempty" mapstructure:"brand_name"`
	BrandSlug       string        `json:"brand_slug,omitempty" mapstructure:"brand_slug"`
	CategoryName    string        `json:"category_name,omitempty" mapstructure:"category_name"`
	CategorySlug    string        `json:"category_slug,omitempty" mapstructure:"category_slug"`
	ColorVariants   []string      `json:"color_variants,omitempty" mapstructure:"color_variants"`
	Color           string        `json:"color,omitempty" mapstructure:"color"`
	Ranking         float64       `json:"ranking,omitempty" mapstructure:"ranking"`
	Tags            []string      `json:"tags,omitempty" mapstructure:"tags"`
	ClickStreams    []ClickStream `json:"click_streams,omitempty"`
	Status          bool          `json:"status,omitempty" mapstructure:"status"`
	ProductImage    string        `json:"product_image,omitempty" mapstructure:"product_image"`
}

type ClickStream struct {
	//TODO: figure out later
}
