package user

import (
	database "PR_2/databases"
	model "PR_2/model"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func LoginUser(c *gin.Context) {

	userCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)
	login_user := new(model.Login)
	
	defer cancel()

	if err:= c.BindJSON(&login_user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		log.Fatal(err)
		return
	}

	var user1 model.User
	var user2 model.User

	//err1 := userCollection.FindOne(ctx, bson.M{"Email": login_user.Login}).Decode(&user1)

	slice := bson.M{
		
			"$or" : [] interface{} {
				bson.M {"Email": login_user.Login},
				bson.M {"Username": login_user.Login},
				bson.M {"Mobile_No": login_user.Login},
			},
		
	}

	err1 := userCollection.FindOne(ctx, slice).Decode(&user1)

	if err1 != nil {

		c.JSON(http.StatusNotFound, gin.H{"message": "login credentials invalid!"})
		return

	} else {	

		err2 := userCollection.FindOne(ctx, bson.M{"Password": login_user.Password}).Decode(&user2)

		if user1.Email != user2.Email || err2 != nil {

			c.JSON(http.StatusNotFound, gin.H{"message1": "password incorrect"})

		}
	
		if user1.Email == user2.Email || user1.MobileNo == user2.MobileNo || user1.Username == user2.Username {

			count_, _ := userCollection.CountDocuments(ctx, slice)

			if count_ > 0 {

				if (login_user.Login == user1.Email || login_user.Login == user1.Username || login_user.Login == user1.MobileNo) && login_user.Password == user2.Password {
				
					update := bson.M{
						
						"$set": bson.M{
							"Login": true ,
						},
						
					}

					_, err := userCollection.UpdateOne(ctx,slice,update)

					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"message": err})
						return

					}

				}

				c.JSON(http.StatusAccepted, gin.H{"hello": user1.Firstname} )

			}

		} 
	}
}
