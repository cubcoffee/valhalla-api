package router

import (
	"encoding/json"
	"github.com/cubcoffee/valhalla-api/credentials"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

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
		v1.DELETE("/employee/:id", deleteEmployeeByID)
		v1.PUT("/employee", updateEmployee)
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

func deleteEmployeeByID(c *gin.Context) {
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
