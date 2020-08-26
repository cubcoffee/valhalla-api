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

func TestGetEmployee(t *testing.T) {

	testServer := httptest.NewServer(routers.CreateRouters())
	defer testServer.Close()

	resp, err := http.Get(fmt.Sprintf("%s/v1/employee/1", testServer.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	b, err := ioutil.ReadAll(resp.Body)

	assert.Equal(t, "{\"id\":1,\"name\":\"Schelb\"}", string(b), "The two words should be the same.")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}
}

func TestGetClient(t *testing.T) {

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

	testServer := httptest.NewServer(routers.CreateRouters())
	defer testServer.Close()

	resp, err := http.Get(fmt.Sprintf("%s/v1/clients", testServer.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "[{\"id\":1,\"name\":\"Jaspion\",\"email\":\"jaspion@daileon.com\",\"phone\":\"55\"},{\"id\":2,\"name\":\"Jiraya\",\"email\":\"jiraya@sucessordetodacuri.com\",\"phone\":\"66\"},{\"id\":3,\"name\":\"Jiban\",\"email\":\"jiban@policaldeaco.com\",\"phone\":\"77\"},{\"id\":4,\"name\":\"Email Duplicado Júnior\",\"email\":\"duplicado@ilegal.com\",\"phone\":\"77\"}]", string(b), "The two words should be the same.")

	if err != nil {
		t.Fatalf("Excpected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}

}

func TestPostClient(t *testing.T) {

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

}
