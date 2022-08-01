package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/breach-simulator/dto"
	"github.com/breach-simulator/entity"
	"github.com/breach-simulator/helper"
	"github.com/breach-simulator/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type BookController interface {
	All(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	Insert(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type bookController struct {
	jwtService  service.JWTService
	bookService service.BookService
}

func NewBookController(bookService service.BookService, jwtService service.JWTService) BookController {
	return &bookController{
		bookService: bookService,
		jwtService:  jwtService,
	}
}

func (c *bookController) All(ctx *gin.Context) {
	var books []entity.Book = c.bookService.All()
	res := helper.BuildSuccessResponse(true, "OK", books)
	ctx.JSON(http.StatusOK, res)
}

func (c *bookController) FindByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("No param ID was found", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	var book entity.Book = c.bookService.FindByID(id)
	if (book == entity.Book{}) {
		res := helper.BuildErrorResponse("Data not found", "No Data", helper.EmptyObj{})
		ctx.JSON(http.StatusNotFound, res)
	} else {
		res := helper.BuildSuccessResponse(true, "OK", book)
		ctx.JSON(http.StatusOK, res)
	}
}

func (c *bookController) Insert(ctx *gin.Context) {
	var bookCreateDTO dto.BookCreateDTO
	errDTO := ctx.ShouldBind(&bookCreateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	authHeader := ctx.GetHeader("Authorization")
	userID := c.getUserIDByToken(authHeader)
	convertedUserID, err := strconv.ParseUint(userID, 10, 64)
	if err == nil {
		bookCreateDTO.UserID = convertedUserID
	}
	book := c.bookService.Insert(bookCreateDTO)
	res := helper.BuildSuccessResponse(true, "OK", book)
	ctx.JSON(http.StatusCreated, res)
}

func (c *bookController) Update(ctx *gin.Context) {
	var bookUpdateDTO dto.BookUpdateDTO
	errDTO := ctx.ShouldBind(&bookUpdateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	authHeader := ctx.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.bookService.IsAllowedToEdit(userID, bookUpdateDTO.ID) {
		id, errID := strconv.ParseUint(userID, 10, 64)
		if errID == nil {
			bookUpdateDTO.UserID = id
		}
		book := c.bookService.Update(bookUpdateDTO)
		res := helper.BuildSuccessResponse(true, "OK", book)
		ctx.JSON(http.StatusOK, res)
	} else {
		res := helper.BuildErrorResponse("Forbidden", "You are not the owner", helper.EmptyObj{})
		ctx.JSON(http.StatusForbidden, res)
	}
}

func (c *bookController) Delete(ctx *gin.Context) {
	var book entity.Book
	id, errID := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if errID != nil {
		res := helper.BuildErrorResponse("Failed to get ID", "No param ID was found", helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
	}
	book.ID = id
	authHeader := ctx.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.bookService.IsAllowedToEdit(userID, book.ID) {
		c.bookService.Delete(book)
		res := helper.BuildSuccessResponse(true, "Deleted", helper.EmptyObj{})
		ctx.JSON(http.StatusOK, res)
	} else {
		res := helper.BuildErrorResponse("Forbidden", "You are not the owner", helper.EmptyObj{})
		ctx.JSON(http.StatusForbidden, res)
	}
}

func (c *bookController) getUserIDByToken(authToken string) string {
	token, err := c.jwtService.ValidateToken(authToken)
	if err != nil {
		panic(err.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	return fmt.Sprintf("%v", claims["user_id"])
}
