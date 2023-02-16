package router

import (
	accounting "PR_2/API/Accounting"
	admin "PR_2/API/Admin"
	appsetting "PR_2/API/AppSetting"
	book "PR_2/API/Books"
	librarians "PR_2/API/Librarian"
	user "PR_2/API/User"
	controllers "PR_2/Controller"

	"github.com/gin-gonic/gin"

	_ "PR_2/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {

	// configure the Gin server
	router := gin.Default()
	 
	router.POST("Book/", book.CreateBook)
	router.GET("getOneBook/:bookId", book.ReadOneBook) 
	router.GET("getAllBook/", book.ReadAllBook)
	router.PUT("/updateBook/:bookId", book.UpdateBook)
	router.DELETE("/deleteBook/:bookId", book.DeleteBook)
	router.GET("getIssuedBook/",book.IssuedBook)
	router.GET("getHistoryBook/",book.HistoryBook)
	router.GET("getQuantityBook/",book.QuantityBook)
	router.PATCH("operationBook/",book.OperationBook)


	router.GET("getOneUser/", user.ReadOneUser) 
	router.GET("getAllUser/", user.ReadAllUser)
	router.PUT("/updateUser/", user.UpdateUser)
	router.DELETE("/deleteUser/", user.DeleteUser)
	
	router.PATCH("UserBookTaken/", user.UserBooksTaken)
	router.PATCH("UserBookReturn/", user.UserBooksReturn)

	router.PATCH("UserSetNewPassword/",user.SetNewPasswordUser)

	router.PATCH("PenaltyUser/", user.UserPenaltyCheck)
	router.PATCH("PenaltyPay/", user.UserPenaltyPay)

	router.GET("UserParam/", user.UserParam)

	router.POST("Admin/", admin.CreateAdmin)
	router.GET("getOneAdmin/:adminId", admin.ReadOneAdmin) 
	router.PUT("/updateAdmin/:adminId", admin.UpdateAdmin)
	router.DELETE("/deleteAdmin/:adminId", admin.DeleteAdmin)

	router.GET("getOneLibrarian/:librarianId", librarians.ReadOneLibrarian) 
	router.GET("getAllLibrarian/", librarians.ReadAllLibrarian)
	router.PUT("/updateLibrarian/:librarianId", librarians.UpdateLibrarian)
	router.DELETE("/deleteLibrarian/:librarianId", librarians.DeleteLibrarian)

	router.POST("/users/signup", controllers.SignUp)
    router.POST("/users/login", controllers.Login)

	router.POST("Accounting/penaltycheck", accounting.AccountingPenaltyCheck)
	router.POST("Accounting/penaltypay",accounting.AccountingPenaltyPay)

	router.POST("CreateSetting/",appsetting.CreateSetting)
	router.PUT("UpdateSetting/",appsetting.UpdateSetting)

	// docs route
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// run the Gin server
	router.Run("localhost:3000")

	return router	
}

