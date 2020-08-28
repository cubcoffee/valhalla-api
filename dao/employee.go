package dao

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
)

//DaysWork days of week
type DaysWork struct {
	ID       uint64 `json:"-"`
	DayIndex string `json:"day_index"`
	UserID   uint64 `json:"-"`
}

//Employee is a representation of an employee
type Employee struct {
	ID             uint64     `json:"id"`
	Name           string     `json:"name"`
	Responsibility string     `json:"responsibility"`
	HourInit       string     `json:"hour_init"`
	HourEnd        string     `json:"hour_end"`
	CredentialID   int64      `json:"-"`
	Credential     Credential `json:"-"`
	DaysWork       []DaysWork `json:"daysWork" gorm:"foreignkey:userId"`
}

func AddEmployee(emp Employee, db *gorm.DB) Employee {

	row := new(Employee)
	d := db.Create(&emp).Scan(row)
	if d.Error != nil {
		log.Print(d.Error)
	}
	return *row
}

func GetEmployeeById(id uint64, db *gorm.DB) Employee {

	emp := Employee{}
	db.Preload("DaysWork").Where("id = " + fmt.Sprint(id)).First(&emp)
	return emp
}

func DeleteEmployeeById(id uint64, db *gorm.DB) {
	emp := Employee{ID: id}
	db.Where("user_id = " + fmt.Sprint(id)).Delete(DaysWork{})
	db.Delete(&emp)
}

func UpdateEmployee(emp Employee, db *gorm.DB) Employee {

	db.Where("user_id = " + fmt.Sprint(emp.ID)).Delete(DaysWork{})
	db.Model(&emp).Updates(emp)

	return emp
}

func GetAllEmployee(db *gorm.DB) []Employee {

	emps := []Employee{}
	db.Preload("DaysWork").Find(&emps)
	return emps
}
