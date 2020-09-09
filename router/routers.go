package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/cubcoffee/valhalla-api/credentials"

	"github.com/cubcoffee/valhalla-api/dao"
	"github.com/cubcoffee/valhalla-api/model"
	"github.com/gin-gonic/gin"
)

func CreateRouters() *gin.Engine {

	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.GET("/hello", helloHandler)

		v1.GET("/employees", getAllEmployees)
		v1.GET("/employee/:id", getEmployeeByID)
		v1.POST("/employee", addEmployee)
		v1.DELETE("/employee/:id", deleteEmployeeById)
		v1.PUT("/employee", updateEmployee)

		v1.GET("/clients", getAllClients)
		v1.GET("/client/:id", getClientById)
		v1.POST("/client", addClient)
		v1.DELETE("/client/:id", deleteClientById)
		v1.PUT("/client/:id", updateClient)

	}

	return r
}

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, Valhalla",
	})
}

func addEmployee(c *gin.Context) {
	db, err := dao.InitDb()
	if err != nil {
		log.Print(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	defer db.Close()
	emp := dao.Employee{}
	reqBody, _ := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(reqBody, &emp)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	credential := dao.Credential{}
	// Fixed password until we don't implement the sent by email password service
	credential.Hash, credential.Salt, err = credentials.GenerateHash("12345678")
	credential, err = dao.AddCredential(credential)
	if err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
	}
	emp.Credential = credential
	emp.CredentialID = credential.ID
	dao.AddEmployee(emp, db)
}

func getEmployeeByID(c *gin.Context) {

	i := c.Param("id")
	id, err := strconv.ParseUint(i, 10, 64)
	if err != nil {
		log.Print(err)
	}

	db, err := dao.InitDb()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()
	emp := dao.GetEmployeeById(id, db)
	c.JSON(http.StatusOK, emp)
}

func deleteEmployeeById(c *gin.Context) {
	i := c.Param("id")
	id, err := strconv.ParseUint(i, 10, 64)
	if err != nil {
		log.Print(err)
	}

	db, err := dao.InitDb()
	if err != nil {
		log.Print(err)
	}

	defer db.Close()

	dao.DeleteEmployeeById(id, db)
	c.Status(http.StatusNoContent)

}

func updateEmployee(c *gin.Context) {
	db, err := dao.InitDb()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	emp := dao.Employee{}
	reqBody, _ := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(reqBody, &emp)

	if emp.ID != 0 {
		row := dao.UpdateEmployee(emp, db)
		c.JSON(http.StatusOK, row)
	} else {
		c.JSON(http.StatusBadRequest, model.ErrorM{Cod: "000", Description: "The id must be entered in the request body"})
	}
}

func getAllEmployees(c *gin.Context) {
	db, err := dao.InitDb()
	if err != nil {
		log.Print(err)
	}
	emps := dao.GetAllEmployee(db)
	c.JSON(http.StatusOK, emps)

	db.Close()
}

func getClientById(c *gin.Context) {

	i := c.Param("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		e := model.Error{
			Message: fmt.Sprintf("The ID must be numeric, but was %v", i),
		}
		c.JSON(http.StatusBadRequest, e)
		return
	}

	db, err := dao.InitDb()
	if err != nil {
		fmt.Println("ERROR", err)

		log.Print(err)
	}
	emp := dao.GetClientById(id, db)

	if emp.ID == 0 {
		err := model.Error{
			Message: fmt.Sprintf("No resource found with this ID: %v", id),
		}
		c.JSON(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, emp)

	defer db.Close()

}

func getAllClients(c *gin.Context) {
	db, err := dao.InitDb()
	if err != nil {
		log.Print(err)
	}
	clients := dao.GetAllClients(db)
	c.JSON(http.StatusOK, clients)

	db.Close()
}

func addClient(c *gin.Context) {
	db, err := dao.InitDb()
	if err != nil {
		log.Print(err)
	}

	client := dao.Client{}
	reqBody, _ := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(reqBody, &client)

	var isInvalid bool
	if client.Email == "" {
		isInvalid = true
	}
	if client.Name == "" {
		isInvalid = true
	}

	if isInvalid {
		err := model.Error{
			Message: "The Client is invalid",
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}

	cli := dao.GetClientByEmail(client.Email, db)
	if cli.ID != 0 {
		err := model.Error{
			Message: fmt.Sprintf("The email %v already exists", client.Email),
		}
		c.JSON(http.StatusBadRequest, err)
		return

	}

	dao.AddClient(client, db)
	defer db.Close()
}

func deleteClientById(c *gin.Context) {
	i := c.Param("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		e := model.Error{
			Message: fmt.Sprintf("The ID must be numeric, but was %v", i),
		}
		c.JSON(http.StatusBadRequest, e)
		return
	}

	db, err := dao.InitDb()
	if err != nil {
		log.Print(err)
	}

	defer db.Close()

	dao.DeleteClientById(id, db)
	c.Status(http.StatusNoContent)

}

func updateClient(c *gin.Context) {

	i := c.Param("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		e := model.Error{
			Message: fmt.Sprintf("The ID must be numeric, but was %v", i),
		}
		c.JSON(http.StatusBadRequest, e)
		return
	}

	db, err := dao.InitDb()
	client := dao.Client{}
	reqBody, _ := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(reqBody, &client)

	var isInvalid bool
	if client.Email == "" {
		isInvalid = true
	}
	if client.Name == "" {
		isInvalid = true
	}

	if isInvalid {
		err := model.Error{
			Message: "The Client is invalid",
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}

	cli := dao.GetClientByEmail(client.Email, db)
	if cli.ID != 0 {
		err := model.Error{
			Message: fmt.Sprintf("The email %v already exists", client.Email),
		}
		c.JSON(http.StatusBadRequest, err)
		return

	}
	dao.UpdateClient(id, client, db)

	defer db.Close()

}
