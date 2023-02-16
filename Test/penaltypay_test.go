package tests

import (
	user "PR_2/API/User"
	database "PR_2/databases"
	localization "PR_2/localise"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPenaltyPay201(t *testing.T) {

	database.NewConnection()
	localization.LoadBundel(path)

	payload := map[string]interface{}{
		"username": "neel_19",
		"pay_amount": 810 ,
	}

	url := "http://localhost:3000/PenaltyPay/"

	payloadToByte, _ := json.Marshal(payload)
	payloadToReader := bytes.NewReader(payloadToByte)

	request, err := http.NewRequest("PATCH", url, payloadToReader)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.Default()
	r.PATCH("PenaltyPay/", user.UserPenaltyPay)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	if status := w.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",status, http.StatusCreated)
	}

}

func TestPenaltyPay400(t *testing.T) {

	database.NewConnection()
	localization.LoadBundel(path)

	payload := map[string]interface{}{
		"username": "vatsal_28",
		"pay_amount": "320",
	}

	url := "http://localhost:3000/PenaltyPay/"

	payloadToByte, _ := json.Marshal(payload)
	payloadToReader := bytes.NewReader(payloadToByte)

	request, err := http.NewRequest("PATCH", url, payloadToReader)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.Default()
	r.PATCH("PenaltyPay/", user.UserPenaltyPay)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	if status := w.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",status, http.StatusBadRequest)
	}

}

func TestPenaltyPay500(t *testing.T) {

	database.NewConnection()
	localization.LoadBundel(path)

	payload := map[string]interface{}{
		"username": "vatsal_28",
		"pay_amount": 320,
	}

	url := "http://localhost:3000/PenaltyPay/"

	payloadToByte, _ := json.Marshal(payload)
	payloadToReader := bytes.NewReader(payloadToByte)

	request, err := http.NewRequest("PATCH", url, payloadToReader)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.Default()
	r.PATCH("PenaltyPay/", user.UserPenaltyPay)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	if status := w.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",status, http.StatusInternalServerError)
	}

}
