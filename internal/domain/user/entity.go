package user

import "time"

type User struct {
	ID               int64
	Username         string
	TgID             int64
	PremiumExpiresAt time.Time
	Status           string
	LimitProd        int64
	Reminder         bool
}
