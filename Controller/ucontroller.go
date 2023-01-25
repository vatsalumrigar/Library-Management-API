package controllers

import (
	middleware "PR_2/Middleware"
	localization "PR_2/localise"
	model "PR_2/model"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	database "PR_2/databases"
	helper "PR_2/helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	logs "github.com/sirupsen/logrus"
)

var (
    lowerCharSet   = "abcdedfghijklmnopqrst"
    upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    specialCharSet = "!@#$%&*"
    numberSet      = "0123456789"
    allCharSet     = lowerCharSet + upperCharSet + specialCharSet + numberSet
)

//HashPassword is used to encrypt the password before it is stored in the DB
func HashPassword(password string) string {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    if err != nil {
        log.Panic(err)
        logs.Error(err)
    }

    return string(bytes)
}

//VerifyPassword checks the input password while verifying it with the passward in the DB.
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
    err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
    check := true
    msg := ""

    if err != nil {
        msg = "passowrd is incorrect"
        check = false
    }

    return check, msg
}


// @Summary signup user
// @ID user-signup
// @Accept json
// @Produce json
// @Param language header string true "languageToken"
// @Param payload body model.User true "Payload for Signup API"
// @Success 200 {object} model.User
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 409 {object} error
// @Failure 500 {object} error
// @Router /users/signup [post]
func SignUp(c *gin.Context) {

    languageToken := c.Request.Header.Get("lan")

	var userCollection = database.GetCollection("User")
	ctx, cancel := database.DbContext(10)
    var user model.User

    if err := c.BindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, localization.GetMessage(languageToken,"400"))
        logs.Error(err.Error())
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
                    c.JSON(http.StatusNotFound, localization.GetMessage(languageToken,"404"))
                    logs.Error(err1)
                    return
                }
            }else {
                c.AbortWithStatusJSON(http.StatusConflict, localization.GetMessage(languageToken,"SignUp.409.error1"))
                logs.Error("enter admin token in header")
                return
            }
    
        }else {
            c.AbortWithStatusJSON(http.StatusConflict, localization.GetMessage(languageToken,"SignUp.409.error2"))
            logs.Error("admin not logged in")
            return
        }

    }

    count, err := userCollection.CountDocuments(ctx, bson.M{"Email": user.Email})
    fmt.Printf("count: %v\n", count)
    defer cancel()
    if err != nil {
        c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"SignUp.500.error1"))
        logs.Error("error occured while checking for the email")
        return
    }

        // password := HashPassword(user.Password)

    rand.Seed(time.Now().UnixNano())
    minSpecialChar := 1
    minNum := 1
    minUpperCase := 1
    minLowerCase := 1
    minlength := 8
    maxlength := 15
    passwordLength := rand.Intn(maxlength-minlength) + minlength

    password := generatePassword(passwordLength, minSpecialChar, minNum, minUpperCase, minLowerCase)

	user.Password = password
	fmt.Println(password)

    user.IsFirstLogin = true
  
    count, err = userCollection.CountDocuments(ctx, bson.M{"Mobile_No": user.MobileNo})
    defer cancel()

    if err != nil {

        c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"SignUp.500.error2"))
        logs.Error("error occured while checking for the phone number")
        return

    }

    if count > 0 {

        c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"SignUp.500.error3"))
        logs.Error("the email or phone number already exists")
        return
            
    }

    user.ID = primitive.NewObjectID()

    resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)

    if insertErr != nil {
        msg := "User item was not created"
        c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"SignUp.500.error4"))
		logs.Error(msg)
        return
    }

    defer cancel()

    c.JSON(http.StatusOK, resultInsertionNumber)

}

// @Summary login user
// @ID user-login
// @Accept json
// @Produce json
// @Param language header string true "languageToken"
// @Param payload body model.Login true "Payload for login API"
// @Success 200 {object} map[string][]string
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /users/login [post]
func Login(c *gin.Context) {

    languageToken := c.Request.Header.Get("lan")
    
	var userCollection = database.GetCollection("User")
	
	ctx, cancel := database.DbContext(10)
    var user model.Login
    var foundUser model.User

    if err := c.BindJSON(&user); err != nil {

        c.JSON(http.StatusBadRequest, localization.GetMessage(languageToken,"400"))
        logs.Error(err.Error())
        return

    }

    err := userCollection.FindOne(ctx, bson.M{"Email": user.Email}).Decode(&foundUser)

    defer cancel()
    if err != nil {
        c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"500"))
        logs.Error(err.Error())
        return
    }

    if foundUser.IsFirstLogin {

        if foundUser.Password != user.Password {

            c.JSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"Login.500"))
            logs.Error("passowrd is incorrect")
            return

        }

        c.AbortWithStatusJSON(http.StatusBadRequest, localization.GetMessage(languageToken,"Login.409"))
        logs.Error("password expired, please create new password through setnewpassword")
        return

    }

    passwordIsValid, msg := VerifyPassword(user.Password, foundUser.Password)
    defer cancel()

    if !passwordIsValid {

        c.AbortWithStatusJSON(http.StatusInternalServerError, localization.GetMessage(languageToken,"Login.500"))
        logs.Error(msg)
        return

    }

    match := bson.M{"Email": foundUser.Email}
    update := bson.M{"Login": true}
 
    _, err1 := userCollection.UpdateOne(ctx,match,bson.M{"$set":update})
    
    if err1 != nil {
        c.JSON(http.StatusInternalServerError,localization.GetMessage(languageToken,"500"))
        logs.Error(err1.Error())
        return
     }

    token, refreshToken, _ := helper.GenerateAllTokens(foundUser.Email, foundUser.Firstname, foundUser.Lastname, foundUser.ID.Hex())


	var res = map[string]interface{}{
		"token" : token,
		"refreshtoken": refreshToken, 
	}


    c.JSON(http.StatusOK, res)

    
}

func generatePassword( passwordLength, minSpecialChar, minNum, minUpperCase, minLowerCase int) string {
    
    var password strings.Builder

    //Set special character
    for i := 0; i < minSpecialChar; i++ {
        random := rand.Intn(len(specialCharSet))
        password.WriteString(string(specialCharSet[random]))
    }

    //Set numeric
    for i := 0; i < minNum; i++ {
        random := rand.Intn(len(numberSet))
        password.WriteString(string(numberSet[random]))
    }

    //Set uppercase
    for i := 0; i < minUpperCase; i++ {
        random := rand.Intn(len(upperCharSet))
        password.WriteString(string(upperCharSet[random]))
    }

	//Set lowercase
	for i := 0; i < minLowerCase; i++ {
        random := rand.Intn(len(lowerCharSet))
        password.WriteString(string(lowerCharSet[random]))
    }

    remainingLength := passwordLength - minSpecialChar - minNum - minUpperCase
    for i := 0; i < remainingLength; i++ {
        random := rand.Intn(len(allCharSet))
        password.WriteString(string(allCharSet[random]))
    }
    inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}