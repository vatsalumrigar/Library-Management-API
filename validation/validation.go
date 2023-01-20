package validation

import (
	database "PR_2/databases"
	"context"
	"errors"
	"regexp"
	logs "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func ValidateBmodel(ctx context.Context, e1 string, e2 string, p string) error {

	author_email := e1
	publisher_email := e2
	published_on := p

	email_pattern := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	published_on_pattern := regexp.MustCompile("(0?[1-9]|[12][0-9]|3[01])/(0?[1-9]|1[012])/((19|20)[0-9][0-9])")

	if author_email == "" || publisher_email == "" {

		if author_email == ""{
			logs.Error("error: author email empty")
		} else {
			logs.Error("error: author email empty")
		}
		return errors.New("author or publisher email cannot be empty")

	}

	if !email_pattern.Match([]byte(author_email)) {

		// fmt.Println("error: change author email : ",author_email)
		logs.Error("error: change author email : ",author_email)
		return errors.New("enter valid author email address")

	}

	if !email_pattern.Match([]byte(publisher_email)) {

		// fmt.Println("error: change publisher email : ",publisher_email)
		logs.Error("error: change publisher email : ",publisher_email)
		return errors.New("enter valid publisher email address")

	}

	if !published_on_pattern.Match([]byte(published_on)) {
		logs.Error("error: date format should be dd/mm/yyyy : ",published_on)
		return errors.New("date format should be dd/mm/yyyy ")

	}

	return nil

}

func ValidateUmodel(ctx context.Context, e string, u string, m string, d string, s string) error {

	
	email := e
	username := u
	mobile_no := m
	dob := d
	status := s
	
	//var DB = database.ConnectDB()

	userCollection := database.GetCollection("User")
	

	count_username, _ := userCollection.CountDocuments(ctx, bson.M{"Username": username})
	count_email, _ := userCollection.CountDocuments(ctx, bson.M{"Email": email})
	count_mobile_no, _ := userCollection.CountDocuments(ctx, bson.M{"MobileNo": mobile_no})

	email_pattern := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	mobile_pattern := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	dob_pattern := regexp.MustCompile("(0?[1-9]|[12][0-9]|3[01])/(0?[1-9]|1[012])/((19|20)[0-9][0-9])")


	if username == "" {

		logs.Error("error: username cannot be empty")
		return errors.New("username cannot be empty")

	}

	if count_username >= 1 {
		
		logs.Error("error: username already exists")
		return errors.New("username already exists")
		
	}

	if email == "" {
		logs.Error("error: user email cannot be empty")
		return errors.New("user email cannot be empty")

	}

	if count_email >= 1 {
		logs.Error("error: user email already exists")
		return errors.New("user email already exists")

	}

	if !email_pattern.Match([]byte(email)) {
		logs.Error("error: enter valid email address")
		return errors.New("enter valid email address")

	}

	if count_mobile_no >= 1 {
		logs.Error("error: mobile number already exists")
		return errors.New("mobile number already exists")

	}

	if !mobile_pattern.Match([]byte(mobile_no)) {
		logs.Error("error: enter valid mobile number")
		return errors.New("enter valid mobile number")

	}

	if !dob_pattern.Match([]byte(dob)) {
		logs.Error("error: date format incorrect")
		return errors.New("date format incorrect")
	}

	if status != "Available" && status != "Unavailable" {
		logs.Error("error: user status should be: Available/Unavailabe")
		return errors.New("user status should be: Available/Unavailabe")
	}

	return nil

}

