package model

//Employee is a representation of an employee
type Employee struct {
	ID             uint64     `json:"id"`
	Name           string     `json:"name"`
	Responsibility string     `json:"responsibility"`
	DaysWork       []DaysWork `json:"daysWork" gorm:"foreignkey:userId"`
}
