// @title Library Management API
// @version 1.0
// @description This is a  Library Management API server.
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /
// @query.collection.format multi
package admin

import (
	database "PR_2/databases"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary delete admin
// @ID delete-admin
// @Produce json
// @Param adminId path string true "AdminID" 
// @Success 201 {object} model.User
// @Failure 500 {object} error
// @Router /deleteAdmin/{adminId} [delete]
func DeleteAdmin(c *gin.Context) {

	adminCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)
	adminId := c.Param("adminId")

	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(adminId) 
	result, err := adminCollection.DeleteOne(ctx, bson.M{"_Id": objId}) 
	res := map[string]interface{}{"data": result}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	
	if result.DeletedCount < 1 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No data to delete"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"message": "Admin deleted successfully", "Data": res})

}