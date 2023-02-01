package book

import (
	database "PR_2/databases"
	model "PR_2/model"
	validation "PR_2/validation"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	logs "github.com/sirupsen/logrus"
	localization "PR_2/localise"
)

// @Summary create book in book collection 
// @ID create-book
// @Accept json
// @Produce json
// @Param payload body model.Books true "Query Payload for create Book API"
// @Param language header string true "languageToken"
// @Success 201 {object} model.Books
// @Failure 400 {object} error
// @Failure 409 {object} error
// @Failure 500 {object} error
// @Router /Book/ [post]
func CreateBook(c *gin.Context) {

	languageToken := c.Request.Header.Get("lan")

	bookCollection := database.GetCollection("Books")
	ctx, cancel := database.DbContext(10)

	book := new(model.Books)
	//author := new(bmodel.Authors)
	
	defer cancel()

	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, localization.GetMessage(languageToken,"400"))
		logs.Error(err.Error())
		return
	} 
	
	for _ , author := range book.Author {
	//fmt.Printf("%d : %s, ", i,  author.Author_Email)
	err := validation.ValidateBmodel(ctx, author.Author_Email, book.Publisher.Publisher_Email, book.Publisher.PublishedOn )

		if err != nil {

			logs.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusConflict, localization.GetMessage(languageToken,"409"))
			return
	
		}
	}
	//err := validateBmodel(ctx, book.Author[0].Email, book.Publisher.Email )


	//appendBookQty(ctx, book.Title, book.Quantities)
	//bookId := c.Param("bookId")
	var result model.Books
	//objId, _ := primitive.ObjectIDFromHex(bookId)
	
	title := book.Title
	quantity := book.Quantities
	count_title, _ := bookCollection.CountDocuments(ctx , bson.M{"Title": title} )

	if count_title >= 1 {

		get1 := bookCollection.FindOne(ctx, bson.M{"Title": title}).Decode(&result)
		if get1 != nil {
			logs.Error(get1.Error())
			c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
			return
		}
	
		qty := quantity + result.Quantities
		update := bson.M{
			
				"$set": bson.M{
					"quantities": qty ,
				},
			
		}
		
		edited := bson.M {

			"Title" : book.Title,
			"Author" : book.Author,
			"Description" : book.Description,
			"Publisher" : book.Publisher,
			"Genre" : book.Genre,
			"Quantities" : qty,
			"Status" : book.Status,	
			"Penalty" : book.Penalty,
			"Cost": book.Cost,
			
		}
		err, _ := bookCollection.UpdateOne(ctx,bson.M{"Title": title},bson.M{"$set": edited})
	
		if err != nil {
			logs.Error(err)
			c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": localization.GetMessage(languageToken,"CreateBook.201.message1"), "Data": update})

	} else {

		addedBook := model.Books{

			ID: primitive.NewObjectID(),
			Title: book.Title,
			Author: book.Author, 
			Description: book.Description,
			Publisher: book.Publisher,
			Genre: book.Genre,
			Quantities: book.Quantities,
			Status: book.Status,
			Penalty: book.Penalty,

		}

		result1, err := bookCollection.InsertOne(ctx, addedBook)

		res1 := map[string]interface{}{"data": result1}

		if err != nil {
			logs.Error(err.Error())
			c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": localization.GetMessage(languageToken,"CreateBook.201.message2"), "Data": res1})

	}

}

