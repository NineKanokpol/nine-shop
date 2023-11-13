package servers

import (
	"github.com/NineKanokpol/Nine-shop-test/modules/middlewares/middlewareHandlers"
	"github.com/NineKanokpol/Nine-shop-test/modules/middlewares/middlewaresRepositories"
	"github.com/NineKanokpol/Nine-shop-test/modules/middlewares/middlewaresUsecases"
	"github.com/NineKanokpol/Nine-shop-test/modules/monitor/monitorHandlers"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitoredModule()
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
