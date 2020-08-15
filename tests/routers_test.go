package routers

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	"github.com/testcontainers/testcontainers-go"
)

func TestMain(m *testing.M) {

	identifier := initCompose()
	time.Sleep(40 * time.Second) //tosco.. eu sei :(

	os.Setenv("DB_TYPE", "mysql")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_NAME", "valhaladb")
	os.Setenv("DB_PASSWORD", "root")

	retCode := m.Run()
	tearDown(identifier)
	os.Exit(retCode)

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
