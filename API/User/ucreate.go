package user

import (
	validation "PR_2/validation"
	database "PR_2/databases"
	model "PR_2/model"
	"time"
	"fmt"
	localization "PR_2/localise"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	logs "github.com/sirupsen/logrus"
)


// @Summary create user in user collection
// @ID create-user
// @Accept json
// @Produce json
// @Param language header string true "languageToken"
// @Param payload body model.User true "Payload for create User API"
// @Success 201 {object} model.User
// @Failure 500 {object} error
// @Router /User/ [post]
func CreateUser(c *gin.Context) {

	languageToken := c.Request.Header.Get("lan")

	userCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)
	user := new(model.User)

	defer cancel()

	if err:= c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, localization.GetMessage(languageToken,"400"))
		logs.Error(err.Error())
		return
	}
	
	err1 := validation.ValidateUmodel(ctx, user.Email, user.Username, user.MobileNo, user.Dob, user.Status)

	if err1 != nil {
		logs.Error(err1.Error())
		c.AbortWithStatusJSON(http.StatusConflict, localization.GetMessage(languageToken,"409"))
		return
	} 
	
	layout := "02/01/2006"
	date := user.Dob
	dateToUnix, err := time.Parse(layout, date) 

	if err != nil {
		logs.Error(err.Error())
	}

	fmt.Println(dateToUnix)
	fmt.Println(dateToUnix.Unix())

	addedUser := model.User {

		ID: primitive.NewObjectID(),
		UserType: user.UserType,
		Firstname : user.Firstname,
		Lastname : user.Lastname,
		Fullname: user.Fullname,
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
		logs.Error(err.Error())
		c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": localization.GetMessage(languageToken,"CreateUser.201"), "Data": res})

}


