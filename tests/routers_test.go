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
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_NAME", "valhaladb")
	os.Setenv("DB_PASSWORD", "root")
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

	assert.Equal(t, "{\"id\":1,\"name\":\"Schelb\",\"responsibility\":\"barbeiro\",\"daysWork\":[{\"day\":\"Sunday\"},{\"day\":\"Monday\"},{\"day\":\"Tuesday\"},{\"day\":\"Saturday\"}]}", string(b), "The two JSON should be the same.")

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

	assert.Contains(t, string(b), "{\"id\":1,\"name\":\"Schelb\",\"responsibility\":\"barbeiro\",\"daysWork\":[{\"day\":\"Sunday\"},{\"day\":\"Monday\"},{\"day\":\"Tuesday\"},{\"day\":\"Saturday\"}]}", "The two JSON should be the same.")

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

	body, _ := json.Marshal(model.Employee{ID: 999, Name: "employee_test_delete", Responsibility: "barbeiro"})

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
