package dao

import (
	"fmt"
	"log"
	"os"

	"github.com/cubcoffee/valhalla-api/model"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func InitDb() (*gorm.DB, error) {

	db, err := gorm.Open(os.Getenv("DB_TYPE"),
		os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@("+os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT")+")/"+os.Getenv("DB_NAME")+"?charset=utf8&parseTime=True&loc=Local")

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

func GetEmployeeById(id int, db *gorm.DB) model.Employee {

	emp := model.Employee{}
	db.Where("id = " + fmt.Sprint(id)).First(&emp)
	return emp
}

func DeleteEmployeeById(id int, db *gorm.DB) {
	emp := model.Employee{ID: id}
	db.Delete(&emp)
}

func GetAllEmployee(db *gorm.DB) []model.Employee {

	emps := []model.Employee{}
	db.Find(&emps)
	return emps
}

/*GetAllClients retrieve all clients*/
func GetAllClients(db *gorm.DB) []model.Client {

	clients := []model.Client{}
	db.Find(&clients)
	return clients
}

/*GetClientByID return a persisted Client from Database*/
func GetClientByID(id int, db *gorm.DB) model.Client {

	client := model.Client{}
	db.Where("id = " + fmt.Sprint(id)).First(&client)
	return client
}

/*DeleteClientById remove a client by ID*/
func DeleteClientById(id int, db *gorm.DB) {
	client := model.Client{ID: id}
	db.Delete(&client)
}

/*AddClient in database*/
func AddClient(cient model.Client, db *gorm.DB) model.Client {

	row := new(model.Client)

	d := db.Create(&cient).Scan(row)
	if d.Error != nil {
		log.Print(d.Error)
	}

	return *row

}

/*GetClientByEmail get client by email*/
func GetClientByEmail(email string, db *gorm.DB) model.Client {

	cli := model.Client{}
	db.Where("email = ?", email).First(&cli)
	return cli
}
