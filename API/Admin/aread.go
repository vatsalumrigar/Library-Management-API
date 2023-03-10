package admin

import (
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"
	"github.com/gin-gonic/gin"
	logs "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	localization "PR_2/localise"
)

// @Summary read admin
// @ID read-admin
// @Produce json
// @Success 201 {object} model.User
// @Param language header string true "languageToken"
// @Param adminId path string true "AdminID"
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /getOneAdmin/{adminId} [get]
func ReadOneAdmin(c *gin.Context)  {

	languageToken := c.Request.Header.Get("lan")

	adminCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)

	defer cancel()

	adminId := c.Param("adminId")
	var result model.Admin


	objId, _ := primitive.ObjectIDFromHex(adminId)

	err := adminCollection.FindOne(ctx, bson.M{"_Id": objId}).Decode(&result)
	

	if err != nil {
		c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
		logs.Error(err.Error())
		return
	}
	//res := map[string]interface{}{"data":result}
	c.JSON(http.StatusOK, gin.H{"message": localization.GetMessage(languageToken,"ReadOneAdmin.200"), "Data": result})

}
