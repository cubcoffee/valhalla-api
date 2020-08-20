package model

//DaysWeek days of week
type DaysWork struct {
	ID     uint64 `json:"-"`
	Day    string `json:"day"`
	UserID uint   `json:"-"`
}
