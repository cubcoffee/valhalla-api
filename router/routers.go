package routers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/cubcoffee/valhalla-api/dao"
	"github.com/cubcoffee/valhalla-api/model"
	"github.com/gin-gonic/gin"
)

func CreateRouters() {

	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.GET("/hello", helloHandler)
		v1.GET("/employee/:id", getEmployeeById)
		v1.POST("/employee", addEmployee)
	}

	r.Run()

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
	}

	emp := model.Employee{}
	reqBody, _ := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(reqBody, &emp)
	dao.AddEmployee(emp, db)
	db.Close()
}

func getEmployeeById(c *gin.Context) {

	i := c.Param("id")
	id, err := strconv.Atoi(i)

	db, err := dao.InitDb()
	if err != nil {
		log.Print(err)
	}

	emp := dao.GetEmployeeById(id, db)
	log.Print(emp)
	c.JSON(http.StatusOK, emp)

	db.Close()

}
