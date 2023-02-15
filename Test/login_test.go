package tests

import (
	controllers "PR_2/Controller"
	database "PR_2/databases"
	localization "PR_2/localise"
	model "PR_2/model"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

const path string = "../"

func TestLogin200(t *testing.T) {

	database.NewConnection()
	localization.LoadBundel(path)

	payload := model.Login{
		Email:    "neel19@gmail.com",
		Password: "nt19",
	}

	url := "http://localhost:3000/users/login"

	payloadToByte, _ := json.Marshal(payload)
	payloadToReader := bytes.NewReader(payloadToByte)

	request, err := http.NewRequest("POST", url, payloadToReader)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.Default()
	r.POST("/users/login", controllers.Login)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func TestLogin400(t *testing.T) {

	database.NewConnection()
	localization.LoadBundel(path)

	payload := map[string]interface{}{
		"emaila":    "vatsal28@gmail.com",
		"passwords": 768,
	}

	url := "http://localhost:3000/users/login"

	payloadToByte, _ := json.Marshal(payload)
	payloadToReader := bytes.NewReader(payloadToByte)

	request, err := http.NewRequest("POST", url, payloadToReader)
	if err != nil {
		t.Fatal(err)
	}
	r := gin.Default()
	r.POST("/users/login", controllers.Login)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	if status := w.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestLogin500(t *testing.T) {

	database.NewConnection()
	localization.LoadBundel(path)

	payload := model.Login{
		Email:    "vatsal28@gmail.com",
		Password: "v8",
	}

	url := "http://localhost:3000/users/login"

	payloadToByte, _ := json.Marshal(payload)
	payloadToReader := bytes.NewReader(payloadToByte)

	request, err := http.NewRequest("POST", url, payloadToReader)
	if err != nil {
		t.Fatal(err)
	}
	r := gin.Default()
	r.POST("/users/login", controllers.Login)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	if status := w.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

}
