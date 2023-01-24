package admin

import (
	//validation "PR_2/validation"
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	logs "github.com/sirupsen/logrus"
	localization "PR_2/localise"
)

// @Summary update admin
// @ID update-admin
// @Accept json
// @Produce json
// @Success 201 {object} model.User
// @Param language header string true "languageToken"
// @Param adminId path string true "AdminID" 
// @Param payload body model.Admin true "Query Payload for update Admin API"
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /updateAdmin/{adminId} [put]
func UpdateAdmin(c *gin.Context) {
	
	languageToken := c.Request.Header.Get("lan")

	adminCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)
	
	
	adminId := c.Param("adminId")
	var admin model.Admin
	
	defer cancel()
	
	objId, _ := primitive.ObjectIDFromHex(adminId)
	
	if err := c.BindJSON(&admin); err != nil {
		logs.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, localization.GetMessage(languageToken,"400"))
		return
	}

	edited := bson.M {

		"UserType" : admin.UserType,
		"Firstname" : admin.Firstname,
		"Lastname" : admin.Lastname,
		"Email" : admin.Email ,
		"MobileNo" : admin.MobileNo ,
		"Password" : admin.Password ,
		"Username" : admin.Username ,
		"Status" : admin.Status,
		"Dob" : admin.Dob,
		"Login": admin.Login,
		"Address": admin.Address,

	}

	result, err := adminCollection.UpdateOne(ctx, bson.M{"_Id": objId}, bson.M{"$set": edited})

	res := map[string]interface{}{"data": result}

	if err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
		return
		}

	if result.MatchedCount < 1 {
		logs.Error("Data doesn't exist")
		c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"UpdateAdmin.500"))
		return
	}
			
	c.JSON(http.StatusCreated, gin.H{"message": localization.GetMessage(languageToken,"201"), "Data": res})	
}

