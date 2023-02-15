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

func TestPenaltyCheck202(t *testing.T) {

	database.NewConnection()
	localization.LoadBundel(path)

	payload := model.PenaltyUsers{
		User_id: []string{
			"63b2ac989361b49c6443dccf",
		},
	}

	url := "http://localhost:3000/PenaltyUser/"

	payloadToByte, _ := json.Marshal(payload)
	payloadToReader := bytes.NewReader(payloadToByte)

	request, err := http.NewRequest("PATCH", url, payloadToReader)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.Default()
	r.PATCH("PenaltyUser/", user.UserPenaltyCheck)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	if status := w.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",status, http.StatusAccepted)
	}

}

func TestPenaltyCheck400(t *testing.T) {

	database.NewConnection()
	localization.LoadBundel(path)

	// payload := model.PenaltyUsers{
	// 	User_id: []string{
	// 		"63b2ac989361b49c6443dcc",
	// 	},
	// }

	payload := map[string]interface{}{
		"userid" : "63b2ac989361b49c6443dccf",
	}

	url := "http://localhost:3000/PenaltyUser/"

	payloadToByte, _ := json.Marshal(payload)
	payloadToReader := bytes.NewReader(payloadToByte)

	request, err := http.NewRequest("PATCH", url, payloadToReader)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.Default()
	r.PATCH("PenaltyUser/", user.UserPenaltyCheck)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	if status := w.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",status, http.StatusBadRequest)
	}

}

func TestPenaltyCheck500(t *testing.T) {

	database.NewConnection()
	localization.LoadBundel(path)

	payload := model.PenaltyUsers{
		User_id: []string{
			"63b2ac989361b49c6443dcc",
		},
	}

	url := "http://localhost:3000/PenaltyUser/"

	payloadToByte, _ := json.Marshal(payload)
	payloadToReader := bytes.NewReader(payloadToByte)

	request, err := http.NewRequest("PATCH", url, payloadToReader)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.Default()
	r.PATCH("PenaltyUser/", user.UserPenaltyCheck)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	if status := w.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",status, http.StatusInternalServerError)
	}

}
