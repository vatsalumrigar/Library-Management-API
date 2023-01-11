package user

import (
	database "PR_2/databases"
	model "PR_2/model"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UserPenaltyCheck(c *gin.Context) {


	bookCollection := database.GetCollection("Books")
	userCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)

	penalty_user := new(model.PenaltyUsers)
	
	defer cancel()
	
	if err:= c.BindJSON(&penalty_user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	var book model.Books

	var IsPenalty []model.IsPenalty

	for _, user := range penalty_user.User_id {
		
		var users model.User

		objId, _ := primitive.ObjectIDFromHex(user)
		err1 := userCollection.FindOne(ctx, bson.M{"_Id": objId}).Decode(&users)

		if err1 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err1})
			return
		}

		total_penalty := users.Total_Penalty

		var ispenalty model.IsPenalty

		for ind , booktaken := range users.BooksTaken {

			err2 := bookCollection.FindOne(ctx,bson.M{"Title": booktaken.Title}).Decode(&book)

			if err2 != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err2})
				return
			}

			timenow := time.Now().Unix() 
			//var timenow int64
			//timenow = 1672396843	
			// timebooktaken :=  booktaken.TimeTaken
			timebooktaken :=  booktaken.TimePenaltyCalc
			// timepenaltypay:=  booktaken.TimePenaltyCalc

			tn := time.Unix(timenow, 0)
			tbt := time.Unix(timebooktaken, 0)
			//tpp := time.Unix(timepenaltypay, 0)

			difference := tn.Sub(tbt)
			daydiff := int(difference.Hours()/24) 
			fmt.Println(daydiff)
			
			if daydiff == 0{
				c.JSON(http.StatusAccepted, gin.H{"message": "no pending penalty"})
				return
			}

			penaltycheckdays := 15
			
			if booktaken.TimePenaltyCalc != booktaken.TimeTaken {
				penaltycheckdays = 0
			}

			if  daydiff > penaltycheckdays {

				total_penalty = total_penalty + (book.Penalty * (daydiff-penaltycheckdays))

				ispenalty.Username = users.Firstname
				ispenalty.Bookname = append(ispenalty.Bookname, booktaken.Title)
				ispenalty.Penalty =  total_penalty

				match := bson.M{"_Id" : objId}
				indtostr := strconv.Itoa(ind)
				update := bson.M{
						
					"$set": bson.M{
						"Status": "Unavailable",
						"Total_Penalty": total_penalty ,
						"Books_Taken."+ indtostr + ".TimePenaltyCalc" : time.Now().Unix() ,
					},
					
				}

				_, err := userCollection.UpdateOne(ctx,match,update)


				if err != nil {
					
					c.JSON(http.StatusInternalServerError, gin.H{"message": "could not updated user penalty"})
					return

				}
				
			} 

		}

		IsPenalty = append(IsPenalty, ispenalty)
	}

	c.JSON(http.StatusAccepted, gin.H{"message": IsPenalty})	

}

