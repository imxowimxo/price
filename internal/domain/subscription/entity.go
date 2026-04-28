package subscription

type Subscription struct {
	ID          int64
	UserID      int64
	ProductID   int64
	TargetPrice float64
	IsTriggered bool
}
