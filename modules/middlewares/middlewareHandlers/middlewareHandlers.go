package middlewareHandlers

import (
	"strings"

	"github.com/NineKanokpol/Nine-shop-test/config"
	"github.com/NineKanokpol/Nine-shop-test/modules/entities"
	"github.com/NineKanokpol/Nine-shop-test/modules/middlewares/middlewaresUsecases"
	"github.com/NineKanokpol/Nine-shop-test/pkg/nineauth"
	"github.com/NineKanokpol/Nine-shop-test/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type middlewareHandlersErrCode string

const (
	RouterCheckErr   middlewareHandlersErrCode = "middleware-001"
	jwtAuthErr       middlewareHandlersErrCode = "middleware-002"
	paramsCheckErr   middlewareHandlersErrCode = "middleware-003"
	authorizationErr middlewareHandlersErrCode = "middleware-004"
	ApiKeyErr        middlewareHandlersErrCode = "middleware-005"
)

type IMiddlewaresHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
	ParamsCheck() fiber.Handler
	Authorize(expectedRoleId ...int) fiber.Handler
	ApiKeyAuth() fiber.Handler
}

type middlewaresHandler struct {
	cfg                 config.IConfig
	middlewaresUsecases middlewaresUsecases.IMiddlewaresUsecase
}

func MiddlewaresHandler(cfg config.IConfig, middlewaresUsecase middlewaresUsecases.IMiddlewaresUsecase) IMiddlewaresHandler {
	return &middlewaresHandler{
		cfg:                 cfg,
		middlewaresUsecases: middlewaresUsecase,
	}
}

func (h *middlewaresHandler) Cors() fiber.Handler {
	//return type เป็น fiber.Handler
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}

// *RouterCheck หากพิมพ์ path มั่วๆ
func (h *middlewaresHandler) RouterCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return entities.NewResponse(c).Error(
			fiber.ErrNotFound.Code,
			string(RouterCheckErr),
			"router not found",
		).Res()
	}
}

func (h *middlewaresHandler) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} [${ip}] ${status} - ${method} ${path}\n", // \n ขึ้นบรรทัดใหม่
		TimeFormat: "01/02/2006",
		TimeZone:   "Bangkok/Asia",
	})
}

func (h *middlewaresHandler) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		//*ฟิล์ด Authorization แล้วก็มี "Bearer xxxxxx"
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		result, err := nineauth.ParseToken(h.cfg.Jwt(), token)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code, //401 UnAuth
				string(jwtAuthErr),
				err.Error(),
			).Res()
		}

		claims := result.Claims
		if !h.middlewaresUsecases.FindAccessToken(claims.Id, token) {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				"no permission to access",
			).Res()
		}

		// Set UserId
		c.Locals("userId", claims.Id)
		c.Locals("userRoleId", claims.RoleId)
		//* c.Next คือการส่งต่อไปที่ฟังก์ชั่นถัดไป
		return c.Next()
	}
}

// *check params ว่า id ใน user ของ access token และการผ่านมาเรียกดูตรงกันไหม
func (h *middlewaresHandler) ParamsCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId")
		if c.Locals("userRoleId").(int) == 2 {
			return c.Next()
		}
		if c.Params("user_id") != userId {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(paramsCheckErr),
				"never gonna give you up",
			).Res()
		}
		return c.Next()
	}
}

// *expectRoleId ...int รับ params แบบไม่จำกัด tyep array
func (h *middlewaresHandler) Authorize(expectedRoleId ...int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleId, ok := c.Locals("userRoleId").(int) //แปลง type
		if !ok {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(authorizationErr),
				"user_id is not int type",
			).Res()
		}

		roles, err := h.middlewaresUsecases.FindRole()
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(authorizationErr),
				err.Error(),
			).Res()
		}
		sum := 0
		for _, v := range expectedRoleId {
			sum += v
		}
		expectedValueBinary := utils.BinaryConverter(sum, len(roles))
		userValueBinary := utils.BinaryConverter(userRoleId, len(roles))
		for i := range userValueBinary {
			if userValueBinary[i]&expectedValueBinary[i] == 1 {
				return c.Next()
			}
		}
		return entities.NewResponse(c).Error(
			fiber.ErrUnauthorized.Code,
			string(authorizationErr),
			"no permission access",
		).Res()
	}
}

// middlewareApikey
func (h *middlewaresHandler) ApiKeyAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Get("X-Api-Key")
		if _, err := nineauth.ParseApiKey(h.cfg.Jwt(), key); err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(ApiKeyErr),
				"apiKey is invalid or required",
			).Res()
		}
		return c.Next()
	}
}
