package serv

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kovey/cli-go/app"
	"github.com/kovey/cli-go/gui"
	"github.com/kovey/db-go/v3/migrate/core"
	"github.com/kovey/db-go/v3/migrate/diff"
	"github.com/kovey/db-go/v3/migrate/mk"
	"github.com/kovey/db-go/v3/migrate/orm"
	"github.com/kovey/db-go/v3/migrate/version"
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
	ta.Add("ksql migrate tools")
	ta.Add(fmt.Sprintf("version: %s", version.Version()))
	ta.Add(fmt.Sprintf("Major: %d", version.MAJOR))
	ta.Add(fmt.Sprintf("Minor: %d", version.MINOR))
	ta.Add(fmt.Sprintf("Build: %d", version.BUILD))
	ta.Show()
}

func (s *serv) Run(a app.AppInterface) error {
	method, err := a.Get("m")
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
	case "make":
		return s._make(a)
	case "orm":
		return s.orm(a)
	default:
		s.ver()
		return nil
	}
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
