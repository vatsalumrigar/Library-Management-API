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

func TestUserParam200(t *testing.T) {

	database.NewConnection()
	localization.LoadBundel(path)

	// var payload interface{}

	// payload = nil

	url := "http://localhost:3000/UserParam/"

	// payloadToByte, _ := json.Marshal(payload)
	// payloadToReader := bytes.NewReader(payloadToByte)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	q := request.URL.Query()
    q.Add("user_id", "63b2ac989361b49c6443dccf")
    request.URL.RawQuery = q.Encode()


	r := gin.Default()
	r.GET("UserParam/", user.UserParam)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",status, http.StatusOK)
	}

}
