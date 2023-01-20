package book

import (
	middleware "PR_2/Middleware"
	database "PR_2/databases"
	model "PR_2/model"
	"strconv"
	"net/http"
	logs "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary get book history
// @ID book-history
// @Accept json
// @Produce json
// @Param strOffset path string true "offset"
// @Param strPageNumber path string true "pagenumber"
// @Param uId header string true "UserID"
// @Param payload body model.HistoryPayload true "Query Payload for Book History API"
// @Success 200 {object} string
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /getHistoryBook/ [get]
func HistoryBook(c *gin.Context) {

	accountingCollection := database.GetCollection("Accounting")
	ctx, cancel := database.DbContext(10)

	defer cancel()

	if middleware.Authentication(c){

		booktitle := new(model.HistoryPayload)
	
		if err := c.BindJSON(&booktitle); err != nil {
			logs.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		} 


		uId,err := c.Get("uid")

		if !err {
			logs.Error(err)
			c.JSON(http.StatusNotFound, gin.H{"message": err})
			return
		}

		userId := uId.(string)
		objId, _ := primitive.ObjectIDFromHex(userId)

		strOffset := c.Request.URL.Query().Get("offset")
		offset, _ := strconv.Atoi(strOffset)
		strPageNumber := c.Request.URL.Query().Get("pageNumber")
		pageNumber, _ := strconv.Atoi(strPageNumber)

		limit := offset
		skip := 0


		if pageNumber != 1{

			skip = (pageNumber - 1) * limit

		}
	

		var user model.Accounting

		err1 := accountingCollection.FindOne(ctx, bson.M{"UserId": objId.Hex()}).Decode(&user)
		if err1 != nil {
			logs.Error(err1.Error())
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

		skipCond := map[string]interface{}{
			"$skip": skip,
		}

		limitCond := map[string]interface{}{
			"$limit": limit,
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

		
 		operationsCond = append(operationsCond,unwindCond,matchCond,sortcond,skipCond,limitCond)


		cursor, err3 := accountingCollection.Aggregate(ctx, operationsCond)

		if err3 != nil {
			logs.Error(err3.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": err3.Error()})
			return
		}
		
		var results []model.Accounting2
		if err3 = cursor.All(ctx, &results); err3 != nil {
			logs.Error(err3.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": err3.Error()})
			return
		}


		phistory:= []interface{}{}
		var totalDocsCount int 
		count := 0

		for _ ,p :=  range  results {

			phistory = append(phistory, p.PenaltyDetail)
			count += 1

		}

		totalDocsCount = count
		totalPageCount := totalDocsCount/offset

		respModel := map[string]interface{}{

			"UserId": user.UserId,
			"Name" : user.Firstname + user.LastName,
			"Email" : user.Email,
			"PenaltyHistory": phistory,
			"TimePenaltyPay": user.TimePenaltyPay,
			"TotalPenalty": user.TotalPenalty,
			"TotalPageCount": totalPageCount,
			"TotalDocCount": totalDocsCount,
						
		}	

		pgNumber := "Page"+ strPageNumber
	
		c.JSON(http.StatusOK, gin.H{pgNumber : respModel})

	}


}
