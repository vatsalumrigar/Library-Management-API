package librarians

import (
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateLibrarian(c *gin.Context){

	librarianCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)

	librarian := new(model.Librarian)

	defer cancel()

	if err := c.BindJSON(&librarian); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	
	/*err1 := validation.ValidateUmodel(ctx, librarian.Email, librarian.Username, librarian.MobileNo, librarian.Dob, librarian.Status)

	if err1 != nil {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err1.Error() })
		return
	} */


	addedLibrarian := model.Librarian {
			
		ID: primitive.NewObjectID(),
		UserType: librarian.UserType,
		Firstname : librarian.Firstname,
		Lastname : librarian.Lastname,
		Email : librarian.Email,
		MobileNo : librarian.MobileNo,
		Password : librarian.Password,
		Username : librarian.Username,
		Status : librarian.Status,
		Dob : librarian.Dob,
		Login : librarian.Login,
		Address : librarian.Address,

	}

	result, err := librarianCollection.InsertOne(ctx, addedLibrarian)
	res := map[string]interface{}{"data": result}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message":err.Error() })
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Posted successfully", "Data": res})

}