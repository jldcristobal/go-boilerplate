package controller

import (
	"net/http"
	"strconv"

	"github.com/breach-simulator/dto"
	"github.com/breach-simulator/entity"
	"github.com/breach-simulator/helper"
	"github.com/breach-simulator/service"
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

// this is where you put your services
type authController struct {
	authService service.AuthService
	jwtService  service.JWTService
}

// create a new instance of AuthController
func NewAuthController(authservice service.AuthService, jwtService service.JWTService) AuthController {
	return &authController{
		authService: authservice,
		jwtService:  jwtService,
	}
}

func (c *authController) Login(ctx *gin.Context) {
	var loginDTO dto.LoginDTO
	err := ctx.ShouldBind(&loginDTO)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	authResult := c.authService.VerifyCredential(loginDTO.Email, loginDTO.Password)
	if v, ok := authResult.(entity.User); ok {
		generatedToken := c.jwtService.GenerateToken(strconv.FormatUint(v.ID, 10))
		v.Token = generatedToken
		response := helper.BuildSuccessResponse(true, "OK", v)
		ctx.JSON(http.StatusOK, response)
		return
	}
	response := helper.BuildErrorResponse("Your login credentials are invalid", "Invalid login credentials", helper.EmptyObj{})
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
}

func (c *authController) Register(ctx *gin.Context) {
	var userCreateDTO dto.UserCreateDTO
	err := ctx.ShouldBind(&userCreateDTO)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	if !c.authService.IsDuplicateEmail(userCreateDTO.Email) {
		response := helper.BuildErrorResponse("Failed to process request", "Email already exists", helper.EmptyObj{})
		ctx.JSON(http.StatusConflict, response)
		return
	} else {
		createdUser := c.authService.CreateUser(userCreateDTO)
		token := c.jwtService.GenerateToken(strconv.FormatUint(createdUser.ID, 10))
		createdUser.Token = token
		response := helper.BuildSuccessResponse(true, "OK", createdUser)
		ctx.JSON(http.StatusCreated, response)
	}
}
