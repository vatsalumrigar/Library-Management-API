package main

import (
	database "PR_2/databases"
	router "PR_2/router"
	"fmt"
)

func main() {

	err := database.NewConnection()

	if err != nil {
		fmt.Println("cannot connect")
	}
	
	router.Router()
		
}

