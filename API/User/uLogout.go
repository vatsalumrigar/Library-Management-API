package user

import(
	database "PR_2/databases"
	model "PR_2/model"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func LogoutUser(c *gin.Context) {

	userCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)
	logout_user := new(model.Logout) 

	defer cancel()

	if err:= c.BindJSON(&logout_user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		log.Fatal(err)
		return
	}

	var user model.User 

	err := userCollection.FindOne(ctx, bson.M{"Email": logout_user.Email}).Decode(&user)

	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{"message1": "email not found"})
		return

	}

	count_email, _ := userCollection.CountDocuments(ctx, bson.M{"Email": logout_user.Email,"Login": true})

	if count_email > 0 {

		update := bson.M{
						
			"$set": bson.M{
				"Login": false,
			},
			
		}

		_, err := userCollection.UpdateOne(ctx, bson.M{"Email": logout_user.Email },update)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return

		}

		c.JSON(http.StatusAccepted, gin.H{"logged out": user.Firstname} )
		return

	} else {

		c.JSON(http.StatusNotFound, gin.H{"message": "user alreday logged out"} )
		return

	}

}