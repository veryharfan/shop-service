package productrepo

type CreateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Category    string `json:"category"`
	ImageURL    string `json:"image_url"`
	ShopID      int64  `json:"shop_id"`
}

type CreateProductResponse struct {
	ID int64 `json:"id"`
}
