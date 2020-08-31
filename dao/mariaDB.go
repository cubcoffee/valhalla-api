package dao

import (
	"fmt"
	"log"
	"os"

	"github.com/cubcoffee/valhalla-api/model"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func InitDb() (*gorm.DB, error) {

	if db != nil {
		return db, nil

	}
	db, err := gorm.Open(os.Getenv("DB_TYPE"), os.Getenv("DB_CONNEC_STRING"))

	if err != nil {
		return nil, err
	}

	return db, nil
}

func AddEmployee(emp model.Employee, db *gorm.DB) model.Employee {

	row := new(model.Employee)
	d := db.Create(&emp).Scan(row)
	if d.Error != nil {
		log.Print(d.Error)
	}
	return *row
}

func GetEmployeeById(id uint64, db *gorm.DB) model.Employee {

	emp := model.Employee{}
	db.Preload("DaysWork").Where("id = " + fmt.Sprint(id)).First(&emp)
	return emp
}

func DeleteEmployeeById(id uint64, db *gorm.DB) {
	emp := model.Employee{ID: id}
	db.Where("user_id = " + fmt.Sprint(id)).Delete(model.DaysWork{})
	db.Delete(&emp)
}

func UpdateEmployee(emp model.Employee, db *gorm.DB) model.Employee {

	db.Where("user_id = " + fmt.Sprint(emp.ID)).Delete(model.DaysWork{})
	db.Model(&emp).Updates(emp)

	return emp
}

func GetAllEmployee(db *gorm.DB) []model.Employee {

	emps := []model.Employee{}
	db.Preload("DaysWork").Find(&emps)
	return emps
}
