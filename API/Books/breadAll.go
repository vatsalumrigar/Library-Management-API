package book

import (
	database "PR_2/databases"
	model "PR_2/model"
	"net/http"
	logs "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// @Summary read books from book collection 
// @ID read-books
// @Accept json
// @Produce json
// @Param queryWord query string false "Book Title"
// @Param payload body model.FilterModel false "Query Payload for Read All Books API"
// @Success 200 {object} model.Books
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /getAllBook/ [get]
func ReadAllBook(c *gin.Context) {

	bookCollection := database.GetCollection("Books")
	ctx, cancel := database.DbContext(10)

	defer cancel()

	var filtermodel model.FilterModel

	if err:= c.ShouldBind(&filtermodel); err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	match := bson.M{}


	queryWord := c.Request.URL.Query().Get("queryWord")

	if queryWord != "" {

		match = bson.M{"Title": bson.M{"$regex": queryWord, "$options": "i"}}

	}

	
	if filtermodel != (model.FilterModel{})  {

		if filtermodel.Title != "" {
			match["Title"] = filtermodel.Title
		}
		if filtermodel.Genre != "" {
			match["Genre"] = filtermodel.Genre
		}
		if filtermodel.Author != "" {
			match["Author.Name"] = filtermodel.Author
		}
		if filtermodel.Publisher != "" {
			match["Publisher.Company_Name"] = filtermodel.Publisher
		}

		//fmt.Printf("match: %v\n", match)
	}



	// match := bson.D{{"Quantities",bson.D{{"$gt",30}}}} // use greater than

	// opts := options.Find().SetLimit(3).SetSort(bson.M{"Quantities": 1}) // limit //sort

	// opts := options.Find().SetSkip(1).SetSort(bson.M{"Quantities": 1}) //skip // sort

	// opts := options.Find().SetProjection(bson.M{"Author":0,"Publisher":0}) // projection //exclude field in result

	// opts := options.Find().SetProjection(bson.M{"_Id": 1,"Publisher": 1}) //projection //include fields in result

	//pipeline
	// groupStage := bson.D{
	// 	{Key: "$group", Value: bson.D{
	// 		{Key: "_id", Value: "$Title"},
	// 		{Key: "average_price", Value: bson.D{{Key: "$avg", Value: "$Cost"}}},
	// 	}}}


	//aggregation
	// groupStage := bson.D{{
	// 	Key: "$addFields",
	// 	Value: bson.D{{
	// 		Key:   "totalspend",
	// 		Value: bson.D{{Key: "$multiply", Value: bson.A{"$Quantities", "$Cost"}}},
	// 	}},
	// }}

	// unsetStage := bson.D{{Key: "$unset", Value: bson.A{"_id", "Author","Publisher"}}}
		
	// // pass the pipeline to the Aggregate() method
	// cursor1, err := bookCollection.Aggregate(ctx, mongo.Pipeline{groupStage,unsetStage})
	// if err != nil {
	// 	panic(err)
	// }

	// var results []bson.M
	// if err = cursor1.All(ctx, &results); err != nil {
	// 	panic(err)
	// }
	
	
	cursor, _ := bookCollection.Find(ctx, match)
	var result []model.Books

	for cursor.Next(ctx){																

		var resl model.Books
		err := cursor.Decode(&resl)

		if err != nil {
			logs.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message":err.Error()})
			return
		}

		result = append(result, resl)


	}

	c.JSON(http.StatusOK, gin.H{"Data": result})

	

}