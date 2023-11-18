package servers

import (
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
}

type moduleFactory struct {
	router fiber.Router
	server *server
	mid    middlewareHandlers.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, mid middlewareHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		router: r,
		server: s,
		mid:    mid,
	}
}

func InitMiddlewares(s *server) middlewareHandlers.IMiddlewaresHandler {
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
	router.Post("/signup", handler.SighUpCustomer)
	router.Post("/signin", handler.SignIn)
	router.Post("/refresh", handler.RefreshPassport)

}
