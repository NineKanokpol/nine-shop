package servers

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/NineKanokpol/Nine-shop-test/config"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// /มี interface ต้องมี struct
type IServer interface {
	Start()
}

type server struct {
	///app ไม่ pass เพราะจะ init ตอน start sever
	app *fiber.App
	db  *sqlx.DB
	cfg config.IConfig
}

func NewServer(cfg config.IConfig, db *sqlx.DB) IServer {
	//*ต้อง return เป็น obj
	return &server{
		cfg: cfg,
		db:  db,
		app: fiber.New(fiber.Config{
			AppName:      cfg.App().Name(),
			BodyLimit:    cfg.App().BodyLimit(),
			ReadTimeout:  cfg.App().ReadTimeout(),
			WriteTimeout: cfg.App().WriteTimeout(),
			JSONEncoder:  json.Marshal, //*แนะนำให้ใช้เพราะทำให้ server go fiber เร็วขึ้น
			JSONDecoder:  json.Unmarshal,
		}),
	}
}

// *เขียน struct ให้มา dot ข้างหน้า
// * vvvvvvvv
func (s *server) Start() {
	//Middlewares
	middlewares := InitMiddlewares(s)

	//*ต้องประกาศหลังจากที่ user ยิง api เข้ามาหา server
	s.app.Use(middlewares.Logger())

	//*Use() ประกาศให้ middlewares เรียกใช้และประกาศทุก end-point
	s.app.Use(middlewares.Cors())

	//Modules
	//https://localhost:3000/
	v1 := s.app.Group("v1")
	///return เป็น interface เอาค่ามารับเพื่อเอาค่าไปใช้ต่อ
	modules := InitModule(v1, s, middlewares)

	modules.MonitoredModule()
	modules.UsersModule()
	modules.AppinfoModule()

	s.app.Use(middlewares.RouterCheck())

	//Graceful shutdown ถ้า server ถูก interrupt จะคืน resource ก่อน ค่อยๆปิดฟังก์ชั่นต่างๆ ก่อนจะปิดตัวแอปลง
	c := make(chan os.Signal, 1) //ประกาศมารับสัญญาณ เช็คตลอดว่าถูก interrupt รึเปล่า โดยต้องเขียนแบบ go concurrency
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		log.Println("server is shutting down....")
		_ = s.app.Shutdown()
	}()

	//listen to host:port
	log.Printf("server starting on %v", s.cfg)
	s.app.Listen(s.cfg.App().Url())
}
