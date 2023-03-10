package user

import (
	database "PR_2/databases"
	model "PR_2/model"
	"fmt"
	"net/http"
	"strconv"
	"time"
	logs "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	localization "PR_2/localise"
)

// @Summary check penalty of users
// @ID users-penalty-check
// @Accept json
// @Produce json
// @Param payload body model.PenaltyUsers true "Payload for Penalty Users API"
// @Param language header string true "languageToken"
// @Success 202 {object} model.IsPenalty
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /PenaltyUser/ [patch]
func UserPenaltyCheck(c *gin.Context) {

	languageToken := c.Request.Header.Get("lan")

	bookCollection := database.GetCollection("Books")
	userCollection := database.GetCollection("User")
	ctx, cancel := database.DbContext(10)

	penalty_user := new(model.PenaltyUsers)
	
	defer cancel()
	
	if err:= c.BindJSON(&penalty_user); err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusBadRequest, localization.GetMessage(languageToken,"400"))
		return
	}

	var book model.Books

	var IsPenalty []model.IsPenalty

	for _, user := range penalty_user.User_id {
		
		var users model.User

		objId, _ := primitive.ObjectIDFromHex(user)
		err1 := userCollection.FindOne(ctx, bson.M{"_Id": objId}).Decode(&users)

		if err1 != nil {
			logs.Error(err1.Error())
			c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
			return
		}

		total_penalty := users.Total_Penalty

		var ispenalty model.IsPenalty

		for ind , booktaken := range users.BooksTaken {

			err2 := bookCollection.FindOne(ctx,bson.M{"Title": booktaken.Title}).Decode(&book)

			if err2 != nil {
				logs.Error(err2.Error())
				c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
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
				logs.Info("no pending penalty of user")
				c.JSON(http.StatusAccepted, localization.GetMessage(languageToken,"UserPenaltyCheck.202"))
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
					
					logs.Error(err.Error())
					c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"UserPenaltyCheck.500"))
					return

				}
				
			} 

		}

		IsPenalty = append(IsPenalty, ispenalty)
	}

	c.JSON(http.StatusAccepted, gin.H{"message": IsPenalty})	

}

