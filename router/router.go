package router

import (
	accounting "PR_2/API/Accounting"
	admin "PR_2/API/Admin"
	book "PR_2/API/Books"
	librarians "PR_2/API/Librarian"
	user "PR_2/API/User"
	controllers "PR_2/Controller"

	"github.com/gin-gonic/gin"
)

func Router(){

	router := gin.Default()
	 
	router.POST("Book/", book.CreateBook)
	router.GET("getOneBook/:bookId", book.ReadOneBook) 
	router.GET("getAllBook/", book.ReadAllBook)
	router.PUT("/updateBook/:bookId", book.UpdateBook)
	router.DELETE("/deleteBook/:bookId", book.DeleteBook)
	router.GET("getIssuedBook/",book.IssuedBook)
	router.GET("getHistoryBook/",book.HistoryBook)

	router.POST("User/", user.CreateUser)
	router.GET("getOneUser/", user.ReadOneUser) 
	router.GET("getAllUser/", user.ReadAllUser)
	router.PUT("/updateUser/", user.UpdateUser)
	router.DELETE("/deleteUser/", user.DeleteUser)
	
	router.PATCH("UserBookTaken/", user.UserBooksTaken)
	router.PATCH("UserBookReturn/", user.UserBooksReturn)

	router.POST("UserLogin/", user.LoginUser)
	router.POST("UserLogout/",user.LogoutUser)

	router.PATCH("PenaltyUser/", user.UserPenaltyCheck)
	router.PATCH("PenaltyPay/", user.UserPenaltyPay)

	router.GET("UserParam/", user.UserParam)

	router.POST("Admin/", admin.CreateAdmin)
	router.GET("getOneAdmin/:adminId", admin.ReadOneAdmin) 
	router.GET("getAllAdmin/", admin.ReadAllAdmin)
	router.PUT("/updateAdmin/:adminId", admin.UpdateAdmin)
	router.DELETE("/deleteAdmin/:adminId", admin.DeleteAdmin)

	router.POST("Librarian/", librarians.CreateLibrarian)
	router.GET("getOneLibrarian/:librarianId", librarians.ReadOneLibrarian) 
	router.GET("getAllLibrarian/", librarians.ReadAllLibrarian)
	router.PUT("/updateLibrarian/:librarianId", librarians.UpdateLibrarian)
	router.DELETE("/deleteLibrarian/:librarianId", librarians.DeleteLibrarian)

	router.POST("/users/signup", controllers.SignUp)
    router.POST("/users/login", controllers.Login)

	router.POST("Accounting/penaltycheck", accounting.AccountingPenaltyCheck)
	router.POST("Accounting/penaltypay",accounting.AccountingPenaltyPay)

	router.Run("localhost:3000")

}