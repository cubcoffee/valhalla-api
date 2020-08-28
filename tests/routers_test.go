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

	err := execError.Error
	if err != nil {
		log.Printf("Could not run compose file: %v - %v", composeFilePaths, err)
	}

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

	body, _ := json.Marshal(dao.Employee{ID: 99, Name: "employee_test1"})

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

func TestDeleteEmployee(t *testing.T) {

	// Create client
	client := &http.Client{}

	testServer := httptest.NewServer(routers.CreateRouters())
	defer testServer.Close()

	body, _ := json.Marshal(dao.Employee{ID: 999, Name: "employee_test_delete", Responsibility: "barbeiro"})

	//Add employee to delete after
	resp, err := http.Post(fmt.Sprintf("%s/v1/employee", testServer.URL), "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Expected no error in create employee_test_delete, got %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200 in create employee_test_delete, got %v", resp.StatusCode)
	}

	//Deleting
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/employee/999", testServer.URL), nil)

	respDel, err := client.Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
		return
	}

	if respDel.StatusCode != 204 {
		t.Fatalf("Expected status code 204, got %v", resp.StatusCode)
	}
}

func TestUpgradeEmployee(t *testing.T) {
	// Create client
	client := &http.Client{}

	testServer := httptest.NewServer(routers.CreateRouters())
	defer testServer.Close()

	body, _ := json.Marshal(dao.Employee{
		ID:             998,
		Name:           "employee_test_update",
		Responsibility: "barbeiro",
		HourInit:       "08:00:00",
		HourEnd:        "18:00:00"})

	//Add employee to update after
	resp, err := http.Post(fmt.Sprintf("%s/v1/employee", testServer.URL), "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Expected no error in create employee_test_update, got %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200 in create employee_test_update, got %v", resp.StatusCode)
	}

	body, _ = json.Marshal(dao.Employee{
		ID:             998,
		Name:           "employee_test_update",
		Responsibility: "atendente",
		HourInit:       "08:00:00",
		HourEnd:        "18:00:00"})

	//Update
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/employee/", testServer.URL), bytes.NewBuffer(body))

	respUp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
		return
	}
	if respUp.StatusCode != 200 {
		t.Fatalf("Expected status code 204, got %v", resp.StatusCode)
	}

	//verify
	b, err := ioutil.ReadAll(respUp.Body)
	assert.Contains(t, string(b), "\"id\":998,\"name\":\"employee_test_update\",\"responsibility\":\"atendente\",", "The two JSON should be the same.")

	resp, err = http.Get(fmt.Sprintf("%s/v1/employee/998", testServer.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	b, err = ioutil.ReadAll(resp.Body)

	assert.Equal(t, "{\"id\":998,\"name\":\"employee_test_update\",\"responsibility\":\"atendente\",\"hour_init\":\"08:00:00\",\"hour_end\":\"18:00:00\",\"daysWork\":[]}", string(b), "The two JSON should be the same.")
}
