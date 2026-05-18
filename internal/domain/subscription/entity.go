package subscription

type Subscription struct {
	ID          int64
	UserID      int64
	ProductID   int64
	TargetPrice float64
	IsTriggered bool
}

type ProductWithSubDTO struct {
	ID           int64
	URL          string
	CurrentPrice float64
	Name         string
	TargetPrice  float64
}
