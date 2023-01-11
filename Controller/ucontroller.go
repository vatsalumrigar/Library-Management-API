package controllers

import (
	middleware "PR_2/Middleware"
	model "PR_2/model"
	"fmt"
	"log"

	"net/http"
	//"time"

	"github.com/gin-gonic/gin"
	// "github.com/go-playground/validator/v10"

	database "PR_2/databases"
	helper "PR_2/helper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	//"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

//HashPassword is used to encrypt the password before it is stored in the DB
func HashPassword(password string) string {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    if err != nil {
        log.Panic(err)
    }

    return string(bytes)
}

//VerifyPassword checks the input password while verifying it with the passward in the DB.
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
    err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
    check := true
    msg := ""

    if err != nil {
        msg = "login or passowrd is incorrect"
        check = false
		fmt.Println(msg)
    }

    return check, msg
}

//CreateUser is the api used to tget a single user
func SignUp(c *gin.Context) {

		var userCollection = database.GetCollection("User")
		ctx, cancel := database.DbContext(10)
        var user model.User

        if err := c.BindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        if user.UserType == "Librarian" {

            if middleware.Authentication(c) {

                ufname, _ := c.Get("first_name")
                userfname := ufname.(string)

                if userfname == "Admin" {

                    uId,err1 := c.Get("uid")
                    fmt.Printf("uId: %v\n", uId)

                    if !err1 {
                        c.JSON(http.StatusNotFound, gin.H{"message": err1})
                        return
                    }
                }else {
                    c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error":"enter admin token in header"})
                    return
                }
        
            }else {
                c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error":"admin not logged in"})
                return
            }

        }

        count, err := userCollection.CountDocuments(ctx, bson.M{"Email": user.Email})
		fmt.Println(count)
        defer cancel()
        if err != nil {
            log.Panic(err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
            return
        }

        password := HashPassword(user.Password)

		user.Password = password
		fmt.Println(password)

        count, err = userCollection.CountDocuments(ctx, bson.M{"Mobile_No": user.MobileNo})
        defer cancel()
        if err != nil {
            log.Panic(err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone number"})
            return
        }

        if count > 0 {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})
            return
        }

        /*user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
        user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
        user.ID = primitive.NewObjectID()
        user.User_id = user.ID.Hex()
        token, refreshToken, _ := helper.GenerateAllTokens(user.Email, user.Firstname, user.Lastname, user.User_id)
        user.Token = token
        user.Refresh_Token = refreshToken*/

        user.ID = primitive.NewObjectID()

        resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)

        if insertErr != nil {
            msg := "User item was not created"
            c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
            return
        }
        defer cancel()

        c.JSON(http.StatusOK, resultInsertionNumber)

    //}
}

//Login is the api used to tget a single user
func Login(c *gin.Context) {

	var userCollection = database.GetCollection("User")
	
		ctx, cancel := database.DbContext(10)
        var user model.User
        var foundUser model.User

        if err := c.BindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        

        err := userCollection.FindOne(ctx, bson.M{"Email": user.Email}).Decode(&foundUser)

        defer cancel()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "login or passowrd is incorrect"})
            return
        }

        passwordIsValid, msg := VerifyPassword(user.Password, foundUser.Password)
        defer cancel()

        if !passwordIsValid {
            c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
            return
        }

        token, refreshToken, _ := helper.GenerateAllTokens(foundUser.Email, foundUser.Firstname, foundUser.Lastname, foundUser.ID.Hex())

        //helper.UpdateAllTokens(token, refreshToken, foundUser.ID.Hex())

		var res = map[string]interface{}{
			"token" : token,
			"refreshtoken": refreshToken, 
		}

        // fmt.Printf("foundUser.ID.Hex(): %v\n", foundUser.ID.Hex())


        c.JSON(http.StatusOK, res)

    
}