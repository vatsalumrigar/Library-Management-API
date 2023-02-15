package tests

import (
	user "PR_2/API/User"
	database "PR_2/databases"
	localization "PR_2/localise"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestReadOneUser200(t *testing.T) {

	database.NewConnection()
	localization.LoadBundel(path)

	token := Login("vatsal28@gmail.com", "vu28")

	url := "http://localhost:3000/getOneUser/"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set("token", token)

	r := gin.Default()
	r.GET("getOneUser/", user.ReadOneUser)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func TestReadOneUser500(t *testing.T) {

	database.NewConnection()
	localization.LoadBundel(path)

	url := "http://localhost:3000/getOneUser/"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.Default()
	r.GET("getOneUser/", user.ReadOneUser)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	if status := w.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

}