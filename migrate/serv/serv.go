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
	a.Flag("m", "version", app.TYPE_STRING, "method: migrate|diff|migplug|make|orm|version")
	a.Flag("d", "mysql", app.TYPE_STRING, "driver: mysql")
	a.Flag("from", "", app.TYPE_STRING, "from dsn")
	a.Flag("to", "", app.TYPE_STRING, "to dsn")
	a.Flag("fromdb", "", app.TYPE_STRING, "from db name")
	a.Flag("todb", "", app.TYPE_STRING, "from db name")
	a.Flag("dir", "", app.TYPE_STRING, "migrates dir when diff")
	a.Flag("p", "", app.TYPE_STRING, "migplug plugin so file path")
	a.Flag("mt", "show", app.TYPE_STRING, "migplug type up|down|show")
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
	for _, flag := range []string{"p", "mt", "to", "d"} {
		if err := s.checkFlag(a, flag); err != nil {
			return err
		}
	}

	p, _ := a.Get("p")
	mt, _ := a.Get("mt")
	to, _ := a.Get("to")
	driver, _ := a.Get("d")
	switch mt.String() {
	case "up":
		return core.LoadPlugin(driver.String(), to.String(), p.String(), core.Type_Up)
	case "down":
		return core.LoadPlugin(driver.String(), to.String(), p.String(), core.Type_Down)
	case "show":
		return core.Show(driver.String(), to.String(), p.String())
	default:
		return fmt.Errorf("mt[%s] unsupport", mt)
	}
}

func (s *serv) diff(a app.AppInterface) error {
	for _, flag := range []string{"from", "to", "fromdb", "todb", "d", "dir"} {
		if err := s.checkFlag(a, flag); err != nil {
			return err
		}
	}

	driver, _ := a.Get("d")
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
	for i, op := range ops {
		if i > 0 {
			if _, err := buff.WriteString("\n"); err != nil {
				debug.Erro("write error: %s", err)
			}
		}

		if _, err := buff.WriteString(op.Prepare()); err != nil {
			debug.Erro("write error: %s", err)
		}
	}

	return buff.Flush()
}

func (s *serv) migrate(a app.AppInterface) error {
	for _, flag := range []string{"to", "todb", "d", "dir"} {
		if err := s.checkFlag(a, flag); err != nil {
			return err
		}
	}

	driver, _ := a.Get("d")
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
	for _, flag := range []string{"n", "v", "to", "dir", "d"} {
		if err := s.checkFlag(a, flag); err != nil {
			return err
		}
	}

	name, _ := a.Get("n")
	version, _ := a.Get("v")
	to, _ := a.Get("to")
	dir, _ := a.Get("dir")
	d, _ := a.Get("d")
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
	d, _ := a.Get("d")
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
		fmt.Println(`
Usage:
	ksql-tool migrate [-dir] [-driver] [-todb] [-to]
		-dir     sql directory
		-driver  database driver(mysql)
		-todb    database name
		-to      database dsn
		`)
	case "diff":
		// "from", "to", "fromdb", "todb", "d", "dir"
		fmt.Println(`
Usage:
	ksql-tool diff [-dir] [-driver] [-fromdb] [-from] [-todb] [-to]
		-dir     created sql directory
		-driver  database driver(mysql)
		-todb    to database name
		-to      to database dsn
		-fromdb  from database name
		-from    from database dsn
		`)
	case "migplug":
		// "p", "mt", "to", "d"
		fmt.Println(`
Usage:
	ksql-tool diff [-plugin] [-driver] [-m] [-to]
		-plugin  plugin directory
		-driver  database driver(mysql)
		-m       plugin method(up|down|show)
		-to      to database dsn
		`)
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
