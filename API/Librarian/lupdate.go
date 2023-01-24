package librarians

import (
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	logs "github.com/sirupsen/logrus"
	localization "PR_2/localise"
)

// @Summary update librarian
// @ID update-librarian
// @Accept json
// @Produce json
// @Param language header string true "languageToken"
// @Param librarianId path string true "LibrarianID" 
// @Param payload body model.User true "Payload for update Librarian API"
// @Success 201 {object} model.User
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /updateLibrarian/{librarianId} [put]
func UpdateLibrarian(c *gin.Context) {

	languageToken := c.Request.Header.Get("lan")
	
	librarianCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)
	
	
	librarianId := c.Param("librarianId")
	var librarian model.User
	
	defer cancel()
	
	objId, _ := primitive.ObjectIDFromHex(librarianId)
	
	if err := c.BindJSON(&librarian); err != nil {
		logs.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, localization.GetMessage(languageToken,"400"))
		return
	}

	edited := bson.M {

		"UserType" : librarian.UserType,
		"Firstname" : librarian.Firstname,
		"Lastname" : librarian.Lastname,
		"Email" : librarian.Email ,
		"MobileNo" : librarian.MobileNo ,
		"Password" : librarian.Password ,
		"Username" : librarian.Username ,
		"Status" : librarian.Status,
		"Dob" : librarian.Dob,
		"Login": librarian.Login,
		"Address": librarian.Address,

	}

	result, err := librarianCollection.UpdateOne(ctx, bson.M{"_Id": objId}, bson.M{"$set": edited})

	res := map[string]interface{}{"data": result}

	if err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
		return
		}

	if result.MatchedCount < 1 {
		logs.Error("Data doesn't exist")
		c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"UpdateLibrarian.500"))
		return
	}
			
	c.JSON(http.StatusCreated, gin.H{"message": localization.GetMessage(languageToken,"UpdateLibrarian.201"), "Data": res})	
}

