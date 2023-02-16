package tests

import (
	controllers "PR_2/Controller"
	database "PR_2/databases"
	localization "PR_2/localise"
	"PR_2/model"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func Login(email, password string) string {

	database.NewConnection()
	localization.LoadBundel(path)

	payload := model.Login{
		Email:    email,
		Password: password,
	}

	url := "http://localhost:3000/users/login"

	payloadToByte, _ := json.Marshal(payload)
	payloadToReader := bytes.NewReader(payloadToByte)

	request, _ := http.NewRequest("POST", url, payloadToReader)

	r := gin.Default()
	r.POST("/users/login", controllers.Login)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	return response["token"]

}

func CreateUser() model.User {
	
	database.NewConnection()
	localization.LoadBundel(path)

	// payload := map[string]interface{}{

	// 	"user_type" : "User",
    //     "first_name": "Alex2",
    //     "last_name": "Demola2",
    //     "full_name": map[string]interface{}{
    //         "en": "Alex Demola",
    //         "hi": "एलेक्स डेमोला",
    //     },
    //     "email": "alex172@gmail.com",
    //     "mobile_no": "8033101555",
    //     "username": "alex_172",
    //     "status": "Available",
    //     "dob": "17/02/2002",
    //     "login": true,
    //     "total_penalty": 0,
	// }

	payload := model.User{
		UserType:      "User",
		Firstname:     "Alex2",
		Lastname:      "Demola2",
		Fullname:      map[string]interface{}{
		    "en": "Alex Demola",
            "hi": "एलेक्स डेमोला",
		},
		Email:         "alex172@gmail.com",
		MobileNo:      "8033101555",
		Username:      "",
		BooksTaken:    []model.Bookdetails{},
		Status:        "alex_172",
		Dob:           "17/02/2002",
		Login:         false,
		Total_Penalty: 0,
		Address:       model.Address{
			Street:  "Althan",
			City:    "Surat",
			State:   "Gujarat",
			Pincode: 395002,
			Country: "India",
		},
	}

	url := "http://localhost:3000/users/signup"

	payloadToByte, _ := json.Marshal(payload)
	payloadToReader := bytes.NewReader(payloadToByte)

	request, _ := http.NewRequest("POST", url, payloadToReader)

	r := gin.Default()
	r.POST("/users/signup", controllers.SignUp)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	return payload

}

func SetNewPassword(){
	
}

func TestSignup200(t *testing.T) {

	database.NewConnection()
	localization.LoadBundel(path)

	payload := map[string]interface{}{
		"user_type":  "User",
		"first_name": "Vatsal1",
		"last_name":  "Umrigar1",
		"full_name": map[string]interface{}{
			"en": "Vatsal Umrigar 1",
			"hi": "वत्सल उमरीगर 1",
		},
		"email":     "vtu@gmail.com",
		"mobile_no": "7069361585",
		"password":  "vut1",

		"username":      "vatsal_281",
		"status":        "Available",
		"dob":           "28/12/2001",
		"total_penalty": 0,
		"address": map[string]interface{}{
			"street":  "Umra",
			"city":    "Surat",
			"state":   "Gujarat",
			"pincode": 395007,
			"country": "India",
		},
	}

	url := "http://localhost:3000/users/signup"

	payloadToByte, _ := json.Marshal(payload)
	payloadToReader := bytes.NewReader(payloadToByte)

	request, err := http.NewRequest("POST", url, payloadToReader)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.Default()
	r.POST("/users/signup", controllers.SignUp)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
func TestSignup409(t *testing.T) {

	database.NewConnection()
	localization.LoadBundel(path)

	token := Login("neel19@gmail.com", "nt19")

	payload := map[string]interface{}{
		"user_type":  "Librarian",
		"first_name": "Vatsal1",
		"last_name":  "Umrigar1",
		"full_name": map[string]interface{}{
			"en": "Vatsal Umrigar 1",
			"hi": "वत्सल उमरीगर 1",
		},
		"email":     "vtu@gmail.com",
		"mobile_no": "7069361585",
		"password":  "vut1",

		"username":      "vatsal_281",
		"status":        "Available",
		"dob":           "28/12/2001",
		"total_penalty": 0,
		"address": map[string]interface{}{
			"street":  "Umra",
			"city":    "Surat",
			"state":   "Gujarat",
			"pincode": 395007,
			"country": "India",
		},
	}

	url := "http://localhost:3000/users/signup"

	payloadToByte, _ := json.Marshal(payload)
	payloadToReader := bytes.NewReader(payloadToByte)

	request, err := http.NewRequest("POST", url, payloadToReader)
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set("token", token)

	r := gin.Default()
	r.POST("/users/signup", controllers.SignUp)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	if status := w.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}
}

func TestSignup500(t *testing.T) {

	database.NewConnection()
	localization.LoadBundel(path)

	payload := map[string]interface{}{
		"user_type":  "User",
		"first_name": "Vatsal1",
		"last_name":  "Umrigar1",
		"full_name": map[string]interface{}{
			"en": "Vatsal Umrigar 1",
			"hi": "वत्सल उमरीगर 1",
		},
		"email":     "vtu@gmail.com",
		"mobile_no": "7069361584",
		"password":  "vut1",

		"username":      "vatsal_281",
		"status":        "Available",
		"dob":           "28/12/2001",
		"total_penalty": 0,
		"address": map[string]interface{}{
			"street":  "Umra",
			"city":    "Surat",
			"state":   "Gujarat",
			"pincode": 395007,
			"country": "India",
		},
	}

	url := "http://localhost:3000/users/signup"

	payloadToByte, _ := json.Marshal(payload)
	payloadToReader := bytes.NewReader(payloadToByte)

	request, err := http.NewRequest("POST", url, payloadToReader)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.Default()
	r.POST("/users/signup", controllers.SignUp)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	if status := w.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}