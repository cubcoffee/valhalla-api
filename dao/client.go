package dao

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
)

type Client struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

/*GetAllClients retrieve all clients*/
func GetAllClients(db *gorm.DB) []Client {

	clients := []Client{}
	db.Find(&clients)
	return clients
}

/*GetClientById return a persisted Client from Database*/
func GetClientById(id int, db *gorm.DB) Client {

	client := Client{}
	db.Where("id = " + fmt.Sprint(id)).First(&client)
	return client
}

/*DeleteClientById remove a client by ID*/
func DeleteClientById(id int, db *gorm.DB) {
	client := Client{ID: id}
	db.Delete(&client)
}

/*DeleteAllClients remove a client by ID*/
func DeleteAllClients(db *gorm.DB) {
	client := Client{}
	db.Where(&client).Delete(Client{})

}

/*AddClient in database*/
func AddClient(cient Client, db *gorm.DB) Client {

	row := new(Client)

	d := db.Create(&cient).Scan(row)
	if d.Error != nil {
		log.Print(d.Error)
	}

	return *row

}

/*GetClientByEmail get client by email*/
func GetClientByEmail(email string, db *gorm.DB) Client {

	cli := Client{}
	db.Where("email = ?", email).First(&cli)
	return cli
}

/*UpdateClient update client by email*/
func UpdateClient(id int, client Client, db *gorm.DB) {

	clientUpdated := Client{}
	db.Where("id = " + fmt.Sprint(id)).First(&clientUpdated)

	clientUpdated.Email = client.Email
	clientUpdated.Name = client.Name
	clientUpdated.Phone = client.Phone

	db.Save(&clientUpdated)

}