package helper

import (
	"fmt"
	"os"
	"time"

	database "PR_2/databases"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
    logs "github.com/sirupsen/logrus"
)

// SignedDetails
type SignedDetails struct {
    Email      string
    Firstname string
    Lastname  string
    Uid        string
    jwt.StandardClaims
}

//var userCollection = database.GetCollection("User")

var SECRET_KEY string = os.Getenv("SECRET_KEY")

// GenerateAllTokens generates both teh detailed token and refresh token
func GenerateAllTokens(email string, firstName string, lastName string, uid string) (signedToken string, signedRefreshToken string, err error) {

    claims := &SignedDetails{
        Email: email,
        Firstname: firstName,
        Lastname: lastName,
        Uid: uid,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Local().Add(time.Hour * 1).Unix(),
        },
    }

    refreshClaims := &SignedDetails{
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
        },
    }


    token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))


    refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

    if err != nil {
        logs.Error(err)
        return
    }

    return token, refreshToken, err
}



//ValidateToken validates the jwt token
func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {

    token, err := jwt.ParseWithClaims(
        signedToken,
        &SignedDetails{},
        func(token *jwt.Token) (interface{}, error) {
            return []byte(SECRET_KEY), nil
        },
    )
    if err != nil {
        msg = err.Error()
        logs.Error(msg)
        return
    }

    claims, ok := token.Claims.(*SignedDetails)
    if !ok {
        msg = "the token is invalid"
        logs.Error(msg)
        msg = err.Error()
        return
    }


    if claims.ExpiresAt < time.Now().Local().Unix() {
        msg = "token is expired"
        logs.Error(msg)
        msg = err.Error()
        return
    }

    return claims, msg
}
//UpdateAllTokens renews the user tokens when they login
func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {

    var userCollection = database.GetCollection("User")

    ctx, cancel := database.DbContext(10)

    var updateObj primitive.D

    updateObj = append(updateObj, bson.E{Key:"Token", Value: signedToken})
    updateObj = append(updateObj, bson.E{Key:"Refresh_Token", Value: signedRefreshToken})

    Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
    updateObj = append(updateObj, bson.E{Key:"Updated_At", Value: Updated_at})

    upsert := true
    filter := bson.M{"_Id": userId}
    opt := options.UpdateOptions{
        Upsert: &upsert,
    }

    fmt.Printf("userId: %v\n", userId)

    _, err := userCollection.UpdateOne(
        ctx,
        filter,
		bson.D{
            {Key: "$set", Value: updateObj},
        },
        &opt,
    )
    
    defer cancel()

    if err != nil {
        logs.Error(err)
        return
    }

}											