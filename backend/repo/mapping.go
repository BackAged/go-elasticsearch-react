package repo

// ProductMapping ...
const ProductMapping = `{
	"properties" : {
	  "id" : {
		"type" : "integer"
	  },
	  "version" : {
		"type" : "integer"
	  },
	  "slug" : {
		"type" : "keyword"
	  },
	  "name" : {
		"type" : "text",
		"fields" : {
		  "keyword" : {
			"type" : "keyword"
		  }
		}
	  },
	  "shop_name" : {
		"type" : "text",
		"fields" : {
		  "keyword" : {
			"type" : "keyword",
			"ignore_above" : 256
		  }
		}
	  },
	  "shop_slug" : {
		"type" : "keyword"
	  },
	  "shop_item_id" : {
		"type" : "integer"
	  },
	  "brand_name" : {
		"type" : "text",
		"fields" : {
		  "keyword" : {
			"type" : "keyword",
			"ignore_above" : 256
		  }
		}
	  },
	  "brand_slug" : {
		"type" : "keyword"
	  },
	  "category_name" : {
		"type" : "text",
		"fields" : {
		  "keyword" : {
			"type" : "keyword",
			"ignore_above" : 256
		  }
		}
	  },
	  "category_slug" : {
		"type" : "keyword"
	  },
	  "color" : {
		  "type" : "text",
		  "fields" : {
			"keyword" : {
			  "type" : "keyword",
			  "ignore_above" : 256
			}
		  }
	  },
	  "color_variants" : {
		"type" : "keyword"
	  },
	  "discounted_price" : {
		"type" : "double"
	  },
	  "max_price" : {
		"type" : "double"
	  },
	  "min_price" : {
		"type" : "double"
	  },
	  "price" : {
		"type" : "double"
	  },
	  "product_image" : {
		"type" : "keyword"
	  },
	  "ranking" : {
		"type" : "rank_features"
	  },
	  "pop_score": {
		  "type" : "double",
		  "null_value" : "0.0"
	  },
	  "tags" : {
		"type" : "keyword"
	  },
	  "click_streams" : {
		"type" : "keyword"
	  },
	  "status" : {
		"type" : "keyword"
	  }
	}
  }`

// ShopMapping ...
const ShopMapping = `{
	"properties" : {
	  "id" : {
		"type" : "integer"
	  },
	  "slug" : {
		"type" : "keyword"
	  },
	  "shop_name" : {
		"type" : "text",
		"fields" : {
		  "keyword" : {
			"type" : "keyword"
		  },
		  "search_as_type": {
			  "type": "search_as_you_type"
		  }
		}
	  },
	  "approval" : {
		"type" : "integer"
	  },
	  "shop_image" : {
		"type" : "keyword"
	  },
	  "owner_name" : {
		"type" : "keyword"
	  },
	  "owner_number" : {
		"type" : "keyword"
	  },
	  "contact_number" : {
		"type" : "keyword"
	  },
	  "ranking" : {
		"type" : "rank_features"
	  },
	  "tags" : {
		"type" : "keyword"
	  },
	  "status" : {
		"type" : "keyword"
	  }
	}
  }`

// BrandMapping ...
const BrandMapping = `{
	"properties" : {
	  "id" : {
		"type" : "integer"
	  },
	  "slug" : {
		"type" : "keyword"
	  },
	  "name" : {
		"type" : "text",
		"fields" : {
		  "keyword" : {
			"type" : "keyword"
		  },
		  "search_as_type": {
			  "type": "search_as_you_type"
		  }
		}
	  },
	  "image_url" : {
		"type" : "keyword"
	  },
	  "ranking" : {
		"type" : "rank_features"
	  },
	  "tags" : {
		"type" : "keyword"
	  },
	  "status" : {
		"type" : "keyword"
	  }
	}
  }`
