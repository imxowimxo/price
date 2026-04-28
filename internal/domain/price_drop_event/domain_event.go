package price_drop_event

type PriceDropEvent struct {
	UserID    int64   `json:"user_id"`
	ProductID int64   `json:"product_id"`
	OldPrice  float64 `json:"old_price"`
	NewPrice  float64 `json:"new_price"`
}
