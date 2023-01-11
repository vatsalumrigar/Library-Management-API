package user

import (
	validation "PR_2/validation"
	database "PR_2/databases"
	model "PR_2/model"
	"time"
	"fmt"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUser(c *gin.Context) {

	userCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)
	user := new(model.User)

	defer cancel()

	if err:= c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		log.Fatal(err)
		return
	}
	
	err1 := validation.ValidateUmodel(ctx, user.Email, user.Username, user.MobileNo, user.Dob, user.Status)

	if err1 != nil {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err1.Error() })
		return
	} 
	
	layout := "02/01/2006"
	date := user.Dob
	dateToUnix, err := time.Parse(layout, date) 

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dateToUnix)
	fmt.Println(dateToUnix.Unix())

	addedUser := model.User {

		ID: primitive.NewObjectID(),
		UserType: user.UserType,
		Firstname : user.Firstname,
		Lastname : user.Lastname,
		Email : user.Email,
		MobileNo : user.MobileNo,
		Password : user.Password,
		Username : user.Username,
		BooksTaken : user.BooksTaken,
		Status : user.Status,
		Dob : user.Dob,
		Login : user.Login,
		Total_Penalty: user.Total_Penalty,

	}

	result, err := userCollection.InsertOne(ctx, addedUser)
	res := map[string]interface{}{"data": result}

	

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message":err.Error() })
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Posted successfully", "Data": res})

}


