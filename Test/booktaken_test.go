package tests

import (
	user "PR_2/API/User"
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

func TestBookTaken201(t *testing.T) {

	database.NewConnection()
	localization.LoadBundel(path)

	token := Login("vatsal28@gmail.com", "vu28")

	payload := map[string]interface{}{
		"title" : "Life Of Pie",
	}

	url := "http://localhost:3000/UserBookTaken/"

	payloadToByte, _ := json.Marshal(payload)
	payloadToReader := bytes.NewReader(payloadToByte)

	request, err := http.NewRequest("PATCH", url, payloadToReader)
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set("token", token)

	r := gin.Default()
	r.PATCH("UserBookTaken/", user.UserBooksTaken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	if status := w.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

}
func TestBookTaken400(t *testing.T) {

	database.NewConnection()
	localization.LoadBundel(path)

	token := Login("neel19@gmail.com", "nt19")

	payload := map[string]interface{}{
		"title" : 20,
	}

	url := "http://localhost:3000/UserBookTaken/"

	payloadToByte, _ := json.Marshal(payload)
	payloadToReader := bytes.NewReader(payloadToByte)

	request, err := http.NewRequest("PATCH", url, payloadToReader)
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set("token", token)

	r := gin.Default()
	r.PATCH("UserBookTaken/", user.UserBooksTaken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	if status := w.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

}
func TestBookTaken500(t *testing.T) {

	database.NewConnection()
	localization.LoadBundel(path)

	payload := model.UserBook{
		Title: "Life Of Pie",
	}

	url := "http://localhost:3000/UserBookTaken/"

	payloadToByte, _ := json.Marshal(payload)
	payloadToReader := bytes.NewReader(payloadToByte)

	request, err := http.NewRequest("PATCH", url, payloadToReader)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.Default()
	r.PATCH("UserBookTaken/", user.UserBooksTaken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	if status := w.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

}