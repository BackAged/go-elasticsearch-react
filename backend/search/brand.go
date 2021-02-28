package search

// Brand defines Brand type
type Brand struct {
	ID       int64  `json:"id" mapstructure:"id"`
	Version  int64  `json:"version" mapstructure:"version"`
	Slug     string `json:"slug" mapstructure:"slug"`
	Name     string `json:"name" mapstructure:"name"`
	ImageURL string `json:"image_url" mapstructure:"image_url"`
}
