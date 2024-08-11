package ninelogger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/NineKanokpol/Nine-shop-test/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type INineLogger interface {
	Print() INineLogger
	Save()
	SetQuery(c *fiber.Ctx) //ดู quey param หลังเครื่องหมาย ? http://localhost:3000/v1/products/:id/
	SetBody(c *fiber.Ctx)
	SetResponse(res any)
}

type nineLogger struct {
	Time       string `json:"time"`
	Ip         string `json:"ip"`
	Method     string `json:"method"`
	StatusCode int    `json:"status_code"`
	Path       string `json:"path"`
	Query      any    `json:"query"`
	Body       any    `json:"body"`
	Response   any    `json:"response"`
}

func InitNineLogger(c *fiber.Ctx, res any) INineLogger {
	log := &nineLogger{
		Time:       time.Now().Local().Format("2006-01-02 15:04:05"),
		Ip:         c.IP(), //ขึ้นอยู่กับ reverse proxy
		Method:     c.Method(),
		Path:       c.Path(),
		StatusCode: c.Response().StatusCode(),
	}
	log.SetQuery(c)
	log.SetBody(c)
	log.SetResponse(res)
	return log
}

func (l *nineLogger) Print() INineLogger {
	utils.Debug(l)
	return l
}
func (l *nineLogger) Save() {
	///input struct เข้าไปและให้ output ออกมาเป็น string format
	data := utils.Output(l)
	filename := fmt.Sprintf("./assets/logs/ninelogger_%v.txt", strings.ReplaceAll(time.Now().Format("2006-01-02"), "-", ""))
	//*0666 คือ file permission  //read write create และ append คือถ้าวันเดียวกันให้ replace ไปเรื่อยๆ
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	file.WriteString(string(data) + "")

}

func (l *nineLogger) SetQuery(c *fiber.Ctx) {
	var body any
	if err := c.QueryParser(&body); err != nil {
		log.Printf("query error parsing error %v", err)
	}
	l.Query = body
}

func (l *nineLogger) SetBody(c *fiber.Ctx) {
	//*setbody จะเกิดขึ้นต่อเมื่อ method เป็น hash ,patch put post
	var body any
	if err := c.BodyParser(&body); err != nil {
		log.Printf("body error parsing error %v", err)
	}

	switch l.Path {
	case "v1/users/signup":
		l.Body = "never gonna give you up"
	default:
		l.Body = body
	}
}

func (l *nineLogger) SetResponse(res any) {
	l.Response = res
}
