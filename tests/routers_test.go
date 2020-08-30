package routers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cubcoffee/valhalla-api/dao"
	"github.com/cubcoffee/valhalla-api/model"
	routers "github.com/cubcoffee/valhalla-api/router"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
)

func TestMain(m *testing.M) {

	identifier := initCompose()
	time.Sleep(40 * time.Second) //tosco.. eu sei :(
	setEnvs()

	retCode := m.Run()

	tearDown(identifier)
	os.Exit(retCode)
}

func setEnvs() {
	os.Setenv("DB_TYPE", "mysql")
	os.Setenv("DB_CONNEC_STRING", "root:root@(localhost:3306)/valhaladb?charset=utf8&parseTime=True&loc=Local")
}

func tearDown(identifier string) {

	composeFilePaths := []string{"./docker-compose.yml"}

	compose := testcontainers.NewLocalDockerCompose(composeFilePaths, identifier)
	execError := compose.Down()
	err := execError.Error
	if err != nil {
		log.Printf("Could not run compose file: %v - %v", composeFilePaths, err)
	}
}

func initCompose() string {

	composeFilePaths := []string{"./docker-compose.yml"}
	identifier := strings.ToLower(uuid.New().String())

	compose := testcontainers.NewLocalDockerCompose(composeFilePaths, identifier)

	execError := compose.
		WithCommand([]string{"up", "-d"}).
		Invoke()

	fmt.Println("EXEC ERROR", execError)

	err := execError.Error
	if err != nil {
		log.Printf("Could not run compose file: %v - %v", composeFilePaths, err)
	}

	fmt.Println("identificer", identifier)
	return identifier
}

func TestHello(t *testing.T) {

	testServer := httptest.NewServer(routers.CreateRouters())
	defer testServer.Close()

	resp, err := http.Get(fmt.Sprintf("%s/v1/hello", testServer.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}
}

func TestPostEmployee(t *testing.T) {

	testServer := httptest.NewServer(routers.CreateRouters())
	defer testServer.Close()

	body, _ := json.Marshal(model.Employee{ID: 99, Name: "employee_test1"})

	resp, err := http.Post(fmt.Sprintf("%s/v1/employee", testServer.URL), "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}
}

func TestGetEmployeeByID(t *testing.T) {

	testServer := httptest.NewServer(routers.CreateRouters())
	defer testServer.Close()

	resp, err := http.Get(fmt.Sprintf("%s/v1/employee/1", testServer.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	b, err := ioutil.ReadAll(resp.Body)

	assert.Equal(t, "{\"id\":1,\"name\":\"Schelb\",\"responsibility\":\"barbeiro\",\"hour_init\":\"08:00:00\",\"hour_end\":\"18:00:00\",\"daysWork\":[{\"day_index\":\"1\"},{\"day_index\":\"1\"},{\"day_index\":\"2\"},{\"day_index\":\"7\"}]}", string(b), "The two JSON should be the same.")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}
}

func TestGetAllEmployee(t *testing.T) {

	testServer := httptest.NewServer(routers.CreateRouters())
	defer testServer.Close()

	resp, err := http.Get(fmt.Sprintf("%s/v1/employees/", testServer.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	b, err := ioutil.ReadAll(resp.Body)

	assert.Contains(t, string(b), "{\"id\":1,\"name\":\"Schelb\",\"responsibility\":\"barbeiro\",\"hour_init\":\"08:00:00\",\"hour_end\":\"18:00:00\",\"daysWork\":[{\"day_index\":\"1\"},{\"day_index\":\"1\"},{\"day_index\":\"2\"},{\"day_index\":\"7\"}]}", "The two JSON should be the same.")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}
}

func TestGetClient(t *testing.T) {

	//Setup
	db, err := dao.InitDb()
	dao.DeleteAllClients(db)
	dao.AddClient(model.Client{

		ID:    1,
		Name:  "Jaspion",
		Email: "jaspion@daileon.com",
		Phone: "55",
	}, db)

	testServer := httptest.NewServer(routers.CreateRouters())
	defer testServer.Close()

	resp, err := http.Get(fmt.Sprintf("%s/v1/client/1", testServer.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	b, err := ioutil.ReadAll(resp.Body)

	assert.Equal(t, "{\"id\":1,\"name\":\"Jaspion\",\"email\":\"jaspion@daileon.com\",\"phone\":\"55\"}", string(b), "The two words should be the same.")

	if err != nil {
		t.Fatalf("Excpected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}

	//Clear
	dao.DeleteClientById(1, db)

	defer db.Close()

}

func TestGetClientWithBadID(t *testing.T) {

	testServer := httptest.NewServer(routers.CreateRouters())
	defer testServer.Close()

	resp, err := http.Get(fmt.Sprintf("%s/v1/client/aaa", testServer.URL))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "{\"message\":\"The ID must be numeric, but was aaa\"}", string(b), "The two words should be the same.")

	if err != nil {
		t.Fatalf("Excpected no error, got %v", err)
	}

	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, got %v", resp.StatusCode)
	}

}

func TestGetNotFoundClient(t *testing.T) {

	db, err := dao.InitDb()
	if err != nil {
		log.Print(err)
	}
	dao.DeleteClientById(5, db)
	defer db.Close()

	testServer := httptest.NewServer(routers.CreateRouters())
	defer testServer.Close()

	resp, err := http.Get(fmt.Sprintf("%s/v1/client/5", testServer.URL))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "{\"message\":\"No resource found with this ID: 5\"}", string(b), "The two words should be the same.")

	if err != nil {
		t.Fatalf("Excpected no error, got %v", err)
	}

	if resp.StatusCode != 404 {
		t.Fatalf("Expected status code 404, got %v", resp.StatusCode)
	}

}

func TestGetAllClients(t *testing.T) {

	//Setup
	db, err := dao.InitDb()
	dao.AddClient(model.Client{
		ID:    1,
		Name:  "Jaspion",
		Email: "jaspion@daileon.com",
		Phone: "55",
	}, db)

	dao.AddClient(model.Client{
		ID:    2,
		Name:  "Jiraya",
		Email: "jiraya@sucessordetodacuri.com",
		Phone: "66",
	}, db)

	dao.AddClient(model.Client{
		ID:    3,
		Name:  "Jiban",
		Email: "jiban@policaldeaco.com",
		Phone: "77",
	}, db)

	dao.AddClient(model.Client{
		ID:    4,
		Name:  "Email Duplicado Júnior",
		Email: "duplicado@ilegal.com",
		Phone: "88",
	}, db)

	testServer := httptest.NewServer(routers.CreateRouters())
	defer testServer.Close()

	resp, err := http.Get(fmt.Sprintf("%s/v1/clients", testServer.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "[{\"id\":1,\"name\":\"Jaspion\",\"email\":\"jaspion@daileon.com\",\"phone\":\"55\"},{\"id\":2,\"name\":\"Jiraya\",\"email\":\"jiraya@sucessordetodacuri.com\",\"phone\":\"66\"},{\"id\":3,\"name\":\"Jiban\",\"email\":\"jiban@policaldeaco.com\",\"phone\":\"77\"},{\"id\":4,\"name\":\"Email Duplicado Júnior\",\"email\":\"duplicado@ilegal.com\",\"phone\":\"88\"}]", string(b), "The two words should be the same.")

	if err != nil {
		t.Fatalf("Excpected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}

	dao.DeleteClientById(1, db)
	dao.DeleteClientById(2, db)
	dao.DeleteClientById(3, db)
	dao.DeleteClientById(4, db)

	defer db.Close()

}

func TestPostClient(t *testing.T) {

	db, err := dao.InitDb()
	dao.DeleteClientById(99, db)

	testServer := httptest.NewServer(routers.CreateRouters())
	defer testServer.Close()

	body, _ := json.Marshal(model.Client{ID: 99, Name: "employee_test1", Email: "duds@23cm.com", Phone: "55"})

	resp, err := http.Post(fmt.Sprintf("%s/v1/client", testServer.URL), "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}

	resp, err = http.Get(fmt.Sprintf("%s/v1/client/99", testServer.URL))
	b, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "{\"id\":99,\"name\":\"employee_test1\",\"email\":\"duds@23cm.com\",\"phone\":\"55\"}", string(b), "The two words should be the same.")

	dao.DeleteClientById(99, db)

	defer db.Close()

}

func TestPostBadClientWithNoName(t *testing.T) {

	testServer := httptest.NewServer(routers.CreateRouters())
	defer testServer.Close()

	body, _ := json.Marshal(model.Client{ID: 100, Name: "", Email: "duds@23cm.com", Phone: "55"})

	resp, err := http.Post(fmt.Sprintf("%s/v1/client", testServer.URL), "application/json", bytes.NewBuffer(body))
	b, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "{\"message\":\"The Client is invalid\"}", string(b), "The two words should be the same.")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, got %v", resp.StatusCode)
	}

	resp, err = http.Get(fmt.Sprintf("%s/v1/client/100", testServer.URL))
	assert.Equal(t, resp.StatusCode, 404, "The two words should be the same.")

}

func TestPostBadClientWithNoEmail(t *testing.T) {

	testServer := httptest.NewServer(routers.CreateRouters())
	defer testServer.Close()

	body, _ := json.Marshal(model.Client{ID: 101, Name: "Duduzão the bala", Email: "", Phone: "55"})

	resp, err := http.Post(fmt.Sprintf("%s/v1/client", testServer.URL), "application/json", bytes.NewBuffer(body))
	b, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "{\"message\":\"The Client is invalid\"}", string(b), "The two words should be the same.")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, got %v", resp.StatusCode)
	}

	resp, err = http.Get(fmt.Sprintf("%s/v1/client/101", testServer.URL))
	assert.Equal(t, resp.StatusCode, 404, "The two words should be the same.")

}

func TestPostBadClientWithSameEmail(t *testing.T) {

	//Setup
	db, err := dao.InitDb()
	dao.AddClient(model.Client{
		ID:    6,
		Name:  "Duduzão  já existente and the bala",
		Email: "duplicado@ilegal.com",
		Phone: "66",
	}, db)

	testServer := httptest.NewServer(routers.CreateRouters())
	defer testServer.Close()

	body, _ := json.Marshal(model.Client{ID: 5, Name: "Duduzão the bala tentando", Email: "duplicado@ilegal.com", Phone: "55"})

	resp, err := http.Post(fmt.Sprintf("%s/v1/client", testServer.URL), "application/json", bytes.NewBuffer(body))
	b, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "{\"message\":\"The email duplicado@ilegal.com already exists\"}", string(b), "The two words should be the same.")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, got %v", resp.StatusCode)
	}

	resp, err = http.Get(fmt.Sprintf("%s/v1/client/5", testServer.URL))
	b, err = ioutil.ReadAll(resp.Body)
	assert.Equal(t, 404, resp.StatusCode, "The two words should be the same.")

	dao.DeleteClientById(6, db)

	defer db.Close()

}

func TestDeleteClient(t *testing.T) {

	db, err := dao.InitDb()
	if err != nil {
		t.Fatal("Error in initializing DB")
	}

	clientInserted := model.Client{
		ID:    999,
		Name:  "You'll die mother fu....",
		Email: "bad@bad.com",
		Phone: "55",
	}

	dao.AddClient(clientInserted, db)

	testServer := httptest.NewServer(routers.CreateRouters())
	defer testServer.Close()

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/client/999", testServer.URL), nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Equal(t, 204, resp.StatusCode, "The two words should be the same.")

	resp, err = http.Get(fmt.Sprintf("%s/v1/client/999", testServer.URL))

	assert.Equal(t, 404, resp.StatusCode, "The two words should be the same.")

	if err != nil {
		t.Fatalf("Excpected no error, got %v", err)
	}

	defer db.Close()
}

func TestDeleteClientWithBadID(t *testing.T) {

	testServer := httptest.NewServer(routers.CreateRouters())
	defer testServer.Close()

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/client/aaa", testServer.URL), nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "{\"message\":\"The ID must be numeric, but was aaa\"}", string(b), "The two words should be the same.")

	if err != nil {
		t.Fatalf("Excpected no error, got %v", err)
	}

	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, got %v", resp.StatusCode)
	}

}

func TestUpdateClient(t *testing.T) {

	db, err := dao.InitDb()
	if err != nil {
		t.Fatal("Error in initializing DB")
	}

	clientInserted := model.Client{
		ID:    999,
		Name:  "You'll die mother fu....",
		Email: "bad@bad.com",
		Phone: "55",
	}

	dao.AddClient(clientInserted, db)

	clientUpdated := model.Client{
		Name:  "I'm new baby!",
		Email: "new@new.com",
		Phone: "66",
	}
	body, _ := json.Marshal(clientUpdated)
	testServer := httptest.NewServer(routers.CreateRouters())
	defer testServer.Close()

	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/v1/client/999", testServer.URL), bytes.NewBuffer(body))
	client := &http.Client{}
	resp, err := client.Do(req)

	assert.Equal(t, 200, resp.StatusCode, "The two words should be the same.")

	resp, err = http.Get(fmt.Sprintf("%s/v1/client/999", testServer.URL))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "{\"id\":999,\"name\":\"I'm new baby!\",\"email\":\"new@new.com\",\"phone\":\"66\"}", string(b), "The two words should be the same.")

	dao.DeleteClientById(999, db)
	defer db.Close()

}
