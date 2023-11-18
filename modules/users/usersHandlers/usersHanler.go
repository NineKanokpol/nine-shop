package usersHandlers

import (
	"github.com/NineKanokpol/Nine-shop-test/config"
	"github.com/NineKanokpol/Nine-shop-test/modules/entities"
	"github.com/NineKanokpol/Nine-shop-test/modules/users"
	"github.com/NineKanokpol/Nine-shop-test/modules/users/usersUsecase"
	"github.com/gofiber/fiber/v2"
)

type usersHandlerErrCode string

const (
	signUpCustomerErr  usersHandlerErrCode = "users-001"
	signInErr          usersHandlerErrCode = "users-002"
	refreshPassportErr usersHandlerErrCode = "users-003"
)

type IUsersHandler interface {
	SighUpCustomer(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
	RefreshPassport(c *fiber.Ctx) error
}

type usersHandler struct {
	cfg         config.IConfig
	userUsecase usersUsecase.IUsersUseCase
}

func UsersHandler(cfg config.IConfig, usersUsecase usersUsecase.IUsersUseCase) IUsersHandler {
	return &usersHandler{
		cfg:         cfg,
		userUsecase: usersUsecase,
	}
}

func (h *usersHandler) SighUpCustomer(c *fiber.Ctx) error {
	//Request body parser
	req := new(users.UserRegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpCustomerErr),
			err.Error(),
		).Res()
	}

	//Email valiation
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpCustomerErr),
			"email pattern is invalid",
		).Res()
	}

	//Insert
	result, err := h.userUsecase.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code, //500
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		}
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}

func (h *usersHandler) SignIn(c *fiber.Ctx) error {
	req := new(users.UserCredentials)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInErr),
			err.Error(),
		).Res()
	}
	passport, err := h.userUsecase.GetPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *usersHandler) RefreshPassport(c *fiber.Ctx) error {
	req := new(users.UserRefreshCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(refreshPassportErr),
			err.Error(),
		).Res()
	}
	passport, err := h.userUsecase.RefreshTokenPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(refreshPassportErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}
