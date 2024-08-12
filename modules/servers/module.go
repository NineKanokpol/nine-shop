package servers

import (
	"github.com/NineKanokpol/Nine-shop-test/modules/appinfo/appinfoHandlers"
	"github.com/NineKanokpol/Nine-shop-test/modules/appinfo/appinfoRepositories"
	"github.com/NineKanokpol/Nine-shop-test/modules/appinfo/appinfoUseCases"
	"github.com/NineKanokpol/Nine-shop-test/modules/middlewares/middlewareHandlers"
	"github.com/NineKanokpol/Nine-shop-test/modules/middlewares/middlewaresRepositories"
	"github.com/NineKanokpol/Nine-shop-test/modules/middlewares/middlewaresUsecases"
	"github.com/NineKanokpol/Nine-shop-test/modules/monitor/monitorHandlers"
	"github.com/NineKanokpol/Nine-shop-test/modules/users/usersHandlers"
	"github.com/NineKanokpol/Nine-shop-test/modules/users/usersRepositories"
	"github.com/NineKanokpol/Nine-shop-test/modules/users/usersUsecase"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitoredModule()
	UsersModule()
	AppinfoModule()
}

type moduleFactory struct {
	router fiber.Router
	server *server
	mid    middlewareHandlers.IMiddlewaresHandler
}

// init module ต้อง init middleware ไปด้วย
func InitModule(r fiber.Router, s *server, mid middlewareHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		router: r,
		server: s,
		mid:    mid,
	}
}

// /init middleware
func InitMiddlewares(s *server) middlewareHandlers.IMiddlewaresHandler {
	///ชั้นใน -> ชั้นนอก
	respository := middlewaresRepositories.MiddlewaresRepository(s.db)
	usecase := middlewaresUsecases.MiddlewaresUsecase(respository)
	return middlewareHandlers.MiddlewaresHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitoredModule() {
	handler := monitorHandlers.MonitorHandlers(m.server.cfg)

	m.router.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UsersRespository(m.server.db)
	usecase := usersUsecase.UsersUseCase(m.server.cfg, repository)
	handler := usersHandlers.UsersHandler(m.server.cfg, usecase)

	//v1/users/sign
	router := m.router.Group("/users")
	router.Post("/signup", m.mid.ApiKeyAuth(), handler.SighUpCustomer)
	router.Post("/signin", m.mid.ApiKeyAuth(), handler.SignIn)
	router.Post("/refresh", m.mid.ApiKeyAuth(), handler.RefreshPassport)
	router.Post("/signout", m.mid.ApiKeyAuth(), handler.SignOut)
	router.Post("/signup-admin", m.mid.JwtAuth(), m.mid.Authorize(2), handler.SignOut)

	//* part parameter :user_id //ขั้นตอน 3 เช็ค id user
	router.Get("/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.GetUserProfile)
	// 2 admin , 1 customer
	router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateAdminToken)

	//Initial admin ขึ้นมา 1 คน ใน db (insert ใน SQL)
	//Gen Admin key
	//ทุกครั้งที่ทำการสมัครแอดมินเพิ่ม ให้ส่ง admin token มาด้วยทุกครั้ง ผ่าน middleware

}

func (m *moduleFactory) AppinfoModule() {
	repository := appinfoRepositories.AppinfoRepository(m.server.db)
	usecase := appinfoUseCases.AppinfoUseCase(repository)
	handler := appinfoHandlers.AppinfoHandler(m.server.cfg, usecase)

	router := m.router.Group("/appinfo")

	router.Get("/apikey", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateApiKey)
	router.Get("/categories", m.mid.ApiKeyAuth(), handler.FindCategory)
}
