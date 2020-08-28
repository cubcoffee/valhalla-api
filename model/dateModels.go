package model

//DaysWeek days of week
type DaysWork struct {
	ID       uint64 `json:"-"`
	DayIndex string `json:"day_index"`
	UserID   uint64 `json:"-"`
}
