package serv

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kovey/cli-go/app"
	"github.com/kovey/cli-go/gui"
	"github.com/kovey/db-go/migrate/core"
	"github.com/kovey/db-go/migrate/diff"
	"github.com/kovey/db-go/migrate/mk"
	"github.com/kovey/db-go/migrate/orm"
	"github.com/kovey/db-go/migrate/version"
	"github.com/kovey/debug-go/debug"
)

type serv struct {
	app.ServBase
}

func (s *serv) Flag(a app.AppInterface) error {
	a.FlagLong("driver", "mysql", app.TYPE_STRING, "driver: mysql")
	a.FlagLong("from", "", app.TYPE_STRING, "from dsn")
	a.FlagLong("to", "", app.TYPE_STRING, "to dsn")
	a.FlagLong("fromdb", "", app.TYPE_STRING, "from db name")
	a.FlagLong("todb", "", app.TYPE_STRING, "from db name")
	a.FlagLong("dir", "", app.TYPE_STRING, "migrates dir when diff")
	a.FlagLong("plugin", "", app.TYPE_STRING, "migplug plugin so file path")
	a.Flag("n", "", app.TYPE_STRING, "migrate name when make use to migplug")
	a.Flag("v", "", app.TYPE_STRING, "migrate version when make use to migplug")
	return nil
}

func (s *serv) Usage() {
	fmt.Println(`
ksql-migrate tools to manage sql migrate、create orm model.

Usage:
	ksql-tool <command> [arguments]
The commands are:
	migrate	 migrate sql from dev to prod
	diff     diff table changed from dev to prod, create changed sql file
	migplug  migrate sql from migration plugins
	orm      create orm model from database
	version  show ksql-tool version
Use "ksql-tool help <command>" for more information about a command.
	`)
}

func (s *serv) Init(app.AppInterface) error {
	return nil
}

func (s *serv) checkFlag(a app.AppInterface, flag string) error {
	val, err := a.Get(flag)
	if err != nil {
		return err
	}

	if val.String() == "" {
		return fmt.Errorf("%s is empty", flag)
	}

	return nil
}

func (s *serv) migplug(a app.AppInterface) error {
	method, err := a.Arg(1, app.TYPE_STRING)
	if err != nil {
		return err
	}

	switch method.String() {
	case "up":
		for _, flag := range []string{"to", "driver"} {
			if err := s.checkFlag(a, flag); err != nil {
				return err
			}
		}

		to, _ := a.Get("to")
		driver, _ := a.Get("driver")
		p, err := a.Get("plugin")
		if err != nil {
			return err
		}
		return core.LoadPlugin(driver.String(), to.String(), p.String(), core.Type_Up)
	case "down":
		for _, flag := range []string{"to", "driver"} {
			if err := s.checkFlag(a, flag); err != nil {
				return err
			}
		}

		for _, flag := range []string{"to", "driver"} {
			if err := s.checkFlag(a, flag); err != nil {
				return err
			}
		}

		to, _ := a.Get("to")
		driver, _ := a.Get("driver")
		p, err := a.Get("plugin")
		if err != nil {
			return err
		}
		return core.LoadPlugin(driver.String(), to.String(), p.String(), core.Type_Down)
	case "show":
		for _, flag := range []string{"to", "driver"} {
			if err := s.checkFlag(a, flag); err != nil {
				return err
			}
		}

		to, _ := a.Get("to")
		driver, _ := a.Get("driver")
		p, err := a.Get("plugin")
		if err != nil {
			return err
		}
		return core.Show(driver.String(), to.String(), p.String())
	case "make":
		return s._make(a)
	case "help":
		return s.helpMigplug(a)
	default:
		return fmt.Errorf("mt[%s] unsupport", method)
	}
}

func (s *serv) helpMigplug(a app.AppInterface) error {
	method, err := a.Arg(2, app.TYPE_STRING)
	if err != nil {
		return err
	}

	switch method.String() {
	case "up":
		upHelp()
	case "down":
		downHelp()
	case "show":
		showHelp()
	case "make":
		makeHelp()
	default:
		return fmt.Errorf("mt[%s] unsupport", method)
	}

	return nil
}

func (s *serv) diff(a app.AppInterface) error {
	for _, flag := range []string{"from", "to", "fromdb", "todb", "driver", "dir"} {
		if err := s.checkFlag(a, flag); err != nil {
			return err
		}
	}

	driver, _ := a.Get("driver")
	if driver.String() != "mysql" {
		return fmt.Errorf("driver[%s] is not mysql", driver)
	}

	dir, _ := a.Get("dir")
	from, _ := a.Get("from")
	to, _ := a.Get("to")
	fromdb, _ := a.Get("fromdb")
	todb, _ := a.Get("todb")
	if err := s.mkdir(dir.String()); err != nil {
		return err
	}

	ops, err := diff.Diff(context.Background(), driver.String(), from.String(), to.String(), fromdb.String(), todb.String())
	if err != nil {
		return err
	}

	file, err := os.Create(dir.String() + fmt.Sprintf("/migrate_%d.sql", time.Now().UnixNano()))
	if err != nil {
		return err
	}

	defer file.Close()
	buff := bufio.NewWriter(file)
	buff.WriteString(fmt.Sprintf("-- from database: %s\n", fromdb))
	buff.WriteString(fmt.Sprintf("-- to database:   %s\n", todb))
	buff.WriteString(fmt.Sprintf("-- tool version:  %s\n", version.Version()))
	buff.WriteString(fmt.Sprintf("-- create time:   %s", time.Now().Format(time.DateTime)))
	for _, op := range ops {
		if _, err := buff.WriteString("\n"); err != nil {
			debug.Erro("write error: %s", err)
		}

		if _, err := buff.WriteString(op.Prepare()); err != nil {
			debug.Erro("write error: %s", err)
		}
	}

	return buff.Flush()
}

func (s *serv) migrate(a app.AppInterface) error {
	for _, flag := range []string{"to", "todb", "driver", "dir"} {
		if err := s.checkFlag(a, flag); err != nil {
			return err
		}
	}

	driver, _ := a.Get("driver")
	if driver.String() != "mysql" {
		return fmt.Errorf("driver[%s] is not mysql", driver)
	}

	dir, _ := a.Get("dir")
	stat, err := os.Stat(dir.String())
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return fmt.Errorf("[%s] not dir", dir.String())
	}

	to, _ := a.Get("to")
	todb, _ := a.Get("todb")

	return diff.Migrate(driver.String(), to.String(), todb.String(), dir.String())
}

func (s *serv) _make(a app.AppInterface) error {
	for _, flag := range []string{"n", "v", "dir"} {
		if err := s.checkFlag(a, flag); err != nil {
			return err
		}
	}

	name, _ := a.Get("n")
	version, _ := a.Get("v")
	to, _ := a.Get("to")
	dir, _ := a.Get("dir")
	d, _ := a.Get("driver")
	return mk.Make(name.String(), version.String(), dir.String(), to.String(), d.String())
}

func (s *serv) orm(a app.AppInterface) error {
	for _, flag := range []string{"to", "dir", "d", "todb"} {
		if err := s.checkFlag(a, flag); err != nil {
			return err
		}
	}

	to, _ := a.Get("to")
	dir, _ := a.Get("dir")
	d, _ := a.Get("driver")
	db, _ := a.Get("todb")
	return orm.Orm(d.String(), to.String(), dir.String(), db.String())
}

func (s *serv) ver() {
	ta := gui.NewTable()
	ta.Add(0, "ksql")
	ta.Add(0, "migrate tools")
	ta.Add(1, "version")
	ta.Add(1, version.Version())
	ta.Add(2, "author")
	ta.Add(2, "kovey")
	ta.Show()
}

func (s *serv) Run(a app.AppInterface) error {
	method, err := a.Arg(0, app.TYPE_STRING)
	if err != nil {
		return err
	}
	switch method.String() {
	case "migrate":
		return s.migrate(a)
	case "diff":
		return s.diff(a)
	case "migplug":
		return s.migplug(a)
	case "version":
		s.ver()
		return nil
	case "orm":
		return s.orm(a)
	case "help":
		return s.help(a)
	default:
		s.Usage()
		return nil
	}
}

func (s *serv) help(a app.AppInterface) error {
	method, err := a.Arg(1, app.TYPE_STRING)
	if err != nil {
		return err
	}

	switch method.String() {
	case "migrate":
		migrateHelp()
	case "diff":
		diffHelp()
	case "migplug":
		migplugHelp()
	case "orm":
		ormHelp()
	default:
		s.Usage()
	}

	return nil
}

func (s *serv) mkdir(dir string) error {
	err := os.MkdirAll(dir, 0755)
	if err == os.ErrExist {
		return nil
	}

	return err
}

func (s *serv) Shutdown(app.AppInterface) error {
	return nil
}
