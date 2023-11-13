package monitorHandlers

import (
	"github.com/NineKanokpol/Nine-shop-test/config"
	"github.com/NineKanokpol/Nine-shop-test/modules/entities"
	"github.com/NineKanokpol/Nine-shop-test/modules/monitor"
	"github.com/gofiber/fiber/v2"
)

//*รับ request จาก network

type IMonitorHandlers interface {
	HealthCheck(c *fiber.Ctx) error
}

type monitorHandler struct {
	cfg config.IConfig
}

func MonitorHandlers(cfg config.IConfig) IMonitorHandlers {
	return &monitorHandler{
		cfg: cfg,
	}
}

func (h *monitorHandler) HealthCheck(c *fiber.Ctx) error {
	res := &monitor.Monitor{
		Name:    h.cfg.App().Name(),
		Version: h.cfg.App().Version(),
	}
	return entities.NewResponse(c).Success(fiber.StatusOK,res).Res()
}
