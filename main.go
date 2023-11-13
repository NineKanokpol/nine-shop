package main

import (
	"os"

	"github.com/NineKanokpol/Nine-shop-test/config"
	"github.com/NineKanokpol/Nine-shop-test/modules/servers"
	"github.com/NineKanokpol/Nine-shop-test/pkg/database"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg := config.LoadConfig(envPath())

	db := database.DbConnection(cfg.Db())
	defer db.Close() //* defer ทำงานท้ายสุดก่อนที่ func จะทำการ return ออกไป

	servers.NewServer(cfg, db).Start()
}
