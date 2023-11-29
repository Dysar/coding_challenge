package model

// OrderRequest represents the JSON request payload
type OrderRequest struct {
	OrderQuantity int `json:"order_quantity"`
}

// PacksNeeded represents the JSON response payload
type PacksNeeded struct {
	OrderQuantity int           `json:"order_quantity"`
	Packs         []PackDetails `json:"packs_needed"`
}

// PackDetails represents details about each pack size and the number of packs needed
type PackDetails struct {
	PackSize   int `json:"pack_size"`
	PacksCount int `json:"packs_count"`
}

type PackSizes struct {
	PackSizes []int `json:"pack_sizes"`
}
