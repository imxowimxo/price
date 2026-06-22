package outbox

import "time"

type PriceDropEvent struct {
	UserID    int64   `json:"user_id"`
	ProductID int64   `json:"product_id"`
	OldPrice  float64 `json:"old_price"`
	NewPrice  float64 `json:"new_price"`
}

type SubscriptionReminderEvent struct {
	UserID    int64     `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

const (
	TopicPriceDrops    = "price_drops"
	TopicSubscriptions = "subscription_reminders"
)
