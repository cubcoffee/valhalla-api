package model

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

/*Client Entity*/
type Client struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func initDB() {

	fmt.Println("DATABASE", &db)
	if db != nil {
		fmt.Println("STATUS", db)
		return
	}

	var err error
	db, err = gorm.Open(os.Getenv("DB_TYPE"), os.Getenv("DB_CONNEC_STRING"))

	if err != nil {
		fmt.Printf("ERROR connection type: %s, string %s. Error: %s", os.Getenv("DB_TYPE"), os.Getenv("DB_CONNEC_STRING"), err)
	}

}

/*GetAll retrieve all clients*/
func (c Client) GetAll() []Client {
	initDB()
	clients := []Client{}
	db.Find(&clients)
	return clients
}

/*GetByID return a persisted Client from Database*/
func (c Client) GetByID() Client {
	initDB()
	client := Client{}
	db.Where("id = " + fmt.Sprint(c.ID)).First(&client)
	return client
}

/*Delete remove a client by ID*/
func (c Client) Delete() {
	initDB()
	client := Client{ID: c.ID}
	db.Delete(&client)
}

/*Add in database*/
func (c Client) Add() Client {
	initDB()
	row := new(Client)

	d := db.Create(&c).Scan(row)
	if d.Error != nil {
		log.Print(d.Error)
	}

	return *row

}

/*GetByEmail get client by email*/
func (c Client) GetByEmail() Client {

	initDB()
	cli := Client{}
	db.Where("email = ?", c.Email).First(&cli)
	return cli
}

/*UpdateClient update client by email*/
func (c Client) UpdateClient(newClient Client) {
	initDB()
	clientUpdated := Client{}
	db.Where("id = " + fmt.Sprint(c.ID)).First(&clientUpdated)

	clientUpdated.Email = newClient.Email
	clientUpdated.Name = newClient.Name
	clientUpdated.Phone = newClient.Phone

	db.Save(&clientUpdated)

}

func (c *Client) SetID(id int) {
	c.ID = id
}

/*DeleteAllClients remove a client by ID*/
func (c Client) DeleteAllClients() {

	initDB()
	client := Client{}
	db.Where(&client).Delete(Client{})

}

/*CloseDB closes DB connection*/
func (c Client) CloseDB() {
	db.Close()
	db = nil
}
