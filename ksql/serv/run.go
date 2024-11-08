package serv

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/kovey/cli-go/app"
	"github.com/kovey/debug-go/debug"
)

func Run() {
	cli := app.NewApp("ksql")
	cli.UsageWhenErr()
	cli.SetServ(&serv{})
	if err := cli.Run(); err != nil {
		debug.Erro("run[%s] error: %s", cli.Name(), err.Error())
	}
}
