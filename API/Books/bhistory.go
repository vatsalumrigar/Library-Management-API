package book

import (
	middleware "PR_2/Middleware"
	database "PR_2/databases"
	model "PR_2/model"
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HistoryBook(c *gin.Context) {

	accountingCollection := database.GetCollection("Accounting")
	ctx, cancel := database.DbContext(10)

	defer cancel()

	if middleware.Authentication(c){

		booktitle := new(model.HistoryPayload)
	
		if err := c.BindJSON(&booktitle); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		} 


		uId,err := c.Get("uid")

		if !err {
			c.JSON(http.StatusNotFound, gin.H{"message": err})
			return
		}

		userId := uId.(string)
		objId, _ := primitive.ObjectIDFromHex(userId)

		var user model.Accounting

		err1 := accountingCollection.FindOne(ctx, bson.M{"UserId": objId.Hex()}).Decode(&user)
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err1.Error()})
			return
		}
		
		operationsCond :=  []interface{}{}

		unwindCond := map[string]interface{}{
			"$unwind": "$PenaltyDetails",
		}

		matchCond := map[string]interface{}{
			"$match": map[string]interface{}{
				"UserId" : objId.Hex(),
				"PenaltyDetails.BookTitle" : booktitle.BookTitle,	
			},
		}
		
		sortcond := map[string]interface{}{
			"$sort": map[string]interface{}{
				"PenaltyDetails.TimePenaltyCheck" : -1,
			},
		}
		// projectStage := map[string]interface{}{
		// 	"$project": map[string]interface{}{
		// 		"PenaltyDetails": map[string]interface{}{
		// 			"$filter" : map[string]interface{}{
		// 				"input": "$PenaltyDetails",
		// 				"as": "user",
		// 				"cond": map[string]interface{}{
		// 					"$eq": []string{"$$user.BookTitle",booktitle.BookTitle},
		// 				},	
		// 			},
		// 		},
		// 	},
		// }

		// projectStage1 := bson.D{{"$project", bson.D{{"PenaltyDetails", bson.D{{"$filter", bson.D{{"input", "$PenaltyDetails"}, {"as", "user"}, {"cond", bson.D{{"$eq", bson.A{"$$user.BookTitle", booktitle.BookTitle}}}}}}}}}}}

		
 		operationsCond = append(operationsCond,unwindCond,matchCond,sortcond)


		cursor, err3 := accountingCollection.Aggregate(ctx, operationsCond)

		if err3 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err3.Error()})
			return
		}
		
		var results []model.Accounting2
		if err3 = cursor.All(ctx, &results); err3 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err3.Error()})
			return
		}

		phistory:= []interface{}{}
		operationsCondcount := []interface{}{}
		count := 0

		for _ ,p :=  range  results {
			phistory = append(phistory, p.PenaltyDetail)
			
			count += 1

			fmt.Printf("count: %v\n", count)

		}

		operationsCondcount = append(operationsCondcount, count)
		fmt.Printf("operationsCondcount: %v\n", operationsCondcount)

		respModel := map[string]interface{}{

			"UserId": user.UserId,
			"Name" : user.Firstname + user.LastName,
			"Email" : user.Email,
			"PenaltyHistory": phistory,
			"TimePenaltyPay": user.TimePenaltyPay,
			"TotalPenalty": user.TotalPenalty,
						
		}	
			
		// c.JSON(http.StatusOK, gin.H{"Data": results})
	
		c.JSON(http.StatusOK, gin.H{"Data": respModel})

	}


}
