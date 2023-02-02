package tests

import (
	database "PR_2/databases"
	"PR_2/router"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	model "PR_2/model"
	// localization "PR_2/localise"
	// logs "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)


func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	exitCode := m.Run()		

	os.Exit(exitCode)
}

func makeRequest(method string, url string, body interface{}) *httptest.ResponseRecorder {

	requestBody,_ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))

	writer:= httptest.NewRecorder()
	router.Router().ServeHTTP(writer,request)

	return writer

}

func TestGetAllUser(t *testing.T){
	
	err := database.NewConnection()

	if err != nil {
		fmt.Println("cannot connect")
	}
	
	writer := makeRequest("GET","http://localhost:3000/getAllUser/",nil)

	assert.Equal(t, http.StatusOK, writer.Code)
}

func TestLogin(t *testing.T){
	err := database.NewConnection()

	if err != nil {
		fmt.Println("cannot connect")
	}

	payload := model.Login{
		Email: "vatsal28@gmail.com",
		Password: "vu28",
	}
	writer := makeRequest("POST","http://localhost:3000/users/login",payload)

	assert.Equal(t, http.StatusOK, writer.Code)
}