package search

// Shop defines Shop type
type Shop struct {
	ID             int64  `json:"id" mapstructure:"id"`
	Version        int64  `json:"version" mapstructure:"version"`
	Slug           string `json:"slug" mapstructure:"slug"`
	Approval       int32  `json:"approval" mapstructure:"approval"`
	ContatctNumber string `json:"contact_number" mapstructure:"contact_number"`
	OwnerName      string `json:"owner_name" mapstructure:"owner_name"`
	OwnerNumber    string `json:"owner_number" mapstructure:"owner_number"`
	ShopImage      string `json:"shop_image" mapstructure:"shop_image"`
	ShopName       string `json:"shop_name" mapstructure:"shop_name"`
}
