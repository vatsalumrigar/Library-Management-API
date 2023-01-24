package admin

import (
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	logs "github.com/sirupsen/logrus"
	localization "PR_2/localise"
)

// @Summary create admin
// @ID create-admin
// @Accept json
// @Produce json
// @Param language header string true "languageToken"
// @Param payload body model.Admin true "Query Payload for create Admin API"
// @Success 201 {object} model.User
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /Admin/ [post]
func CreateAdmin(c *gin.Context){

	languageToken := c.Request.Header.Get("lan")

	adminCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)

	admin := new(model.Admin)
	//address := new(model.AdminAddress)

	defer cancel()

	if err := c.BindJSON(&admin); err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusBadRequest, localization.GetMessage(languageToken,"400"))
		return
	}

	addedAdmin := model.Admin {
			
		ID: primitive.NewObjectID(),
		UserType: admin.UserType,
		Firstname : admin.Firstname,
		Lastname : admin.Lastname,
		Email : admin.Email,
		MobileNo : admin.MobileNo,
		Password : admin.Password,
		Username : admin.Username,
		Status : admin.Status,
		Dob : admin.Dob,
		Login : admin.Login,
		Address: admin.Address,

	}

	result, err := adminCollection.InsertOne(ctx, addedAdmin)
	res := map[string]interface{}{"data": result}

	if err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": localization.GetMessage(languageToken,"200"), "Data": res})

}