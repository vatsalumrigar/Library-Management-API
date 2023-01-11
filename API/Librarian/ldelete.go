package librarians

import (
	database "PR_2/databases"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func DeleteLibrarian(c *gin.Context) {

	librarianCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)
	librarianId := c.Param("librarianId")

	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(librarianId) 
	result, err := librarianCollection.DeleteOne(ctx, bson.M{"_Id": objId}) 
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