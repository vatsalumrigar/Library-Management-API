// @title Library Management API
// @version 1.0
// @description This is a  Library Management API server.
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @host localhost:3000
// @BasePath /
// @query.collection.format multi
package main

import (
	database "PR_2/databases"
	router "PR_2/router"
	"fmt"
	"os"

	logs "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logs.SetFormatter(&logs.JSONFormatter{})
  
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logs.SetOutput(os.Stdout)
  
	// Only log the warning severity or above.
	logs.SetLevel(logs.InfoLevel)
  }

func main() {

	err := database.NewConnection()

	if err != nil {
		fmt.Println("cannot connect")
	}
	
	router.Router()
		
}

