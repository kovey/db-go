package serv

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/kovey/cli-go/app"
	"github.com/kovey/cli-go/gui"
	"github.com/kovey/db-go/ksql/core"
	"github.com/kovey/db-go/ksql/diff"
	"github.com/kovey/db-go/ksql/mk"
	"github.com/kovey/db-go/ksql/orm"
	"github.com/kovey/db-go/ksql/version"
	"github.com/kovey/debug-go/debug"
)

type serv struct {
	app.ServBase
}

func (s *serv) getFromDsn() string {
	if os.Getenv("DB_DRIVER") != "mysql" {
		return ""
	}

	conf := mysql.NewConfig()
	conf.Loc = time.Local
	conf.User = os.Getenv("DB_USER")
	conf.Passwd = os.Getenv("DB_PASSWORD")
	conf.Net = "tcp"
	conf.Addr = fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
	conf.DBName = os.Getenv("DB_NAME")
	conf.Params = map[string]string{
		"charset": os.Getenv("DB_CHARSET"),
	}

	return conf.FormatDSN()
}

func (s *serv) getToDsn() string {
	if os.Getenv("DB_DRIVER") != "mysql" {
		return ""
	}

	conf := mysql.NewConfig()
	conf.Loc = time.Local
	conf.User = os.Getenv("TO_DB_USER")
	conf.Passwd = os.Getenv("TO_DB_PASSWORD")
	conf.Net = "tcp"
	conf.Addr = fmt.Sprintf("%s:%s", os.Getenv("TO_DB_HOST"), os.Getenv("TO_DB_PORT"))
	conf.DBName = os.Getenv("TO_DB_NAME")
	conf.Params = map[string]string{
		"charset": os.Getenv("TO_DB_CHARSET"),
	}

	return conf.FormatDSN()
}

func (s *serv) Flag(a app.AppInterface) error {
	a.FlagLong("driver", os.Getenv("DB_DRIVER"), app.TYPE_STRING, "driver: mysql")
	a.FlagLong("from", s.getFromDsn(), app.TYPE_STRING, "from dsn")
	a.FlagLong("to", s.getToDsn(), app.TYPE_STRING, "to dsn")
	a.FlagLong("fromdb", os.Getenv("DB_NAME"), app.TYPE_STRING, "from db name")
	a.FlagLong("todb", os.Getenv("TO_DB_NAME"), app.TYPE_STRING, "from db name")
	a.FlagLong("dir", "", app.TYPE_STRING, "migrates dir when diff")
	a.Flag("n", "", app.TYPE_STRING, "migrate name when make use to migplug")
	a.Flag("v", "", app.TYPE_STRING, "migrate version when make use to migplug")
	return nil
}

func (s *serv) Usage() {
	fmt.Println(`
ksql tools to manage sql migrate、create orm model.

Usage:
	ksql <command> [arguments]
The commands are:
	migrate	 migrate sql from dev to prod
	diff     diff table changed from dev to prod, create changed sql file
	migplug  migrate sql from migration plugins
	orm      create orm model from database
	config   manage config file
	version  show ksql version
Use "ksql help <command>" for more information about a command.
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
		for _, flag := range []string{"from", "driver", "v"} {
			if err := s.checkFlag(a, flag); err != nil {
				return err
			}
		}

		from, _ := a.Get("from")
		driver, _ := a.Get("driver")
		version, _ := a.Get("v")
		var plugin = ""
		if p, err := a.Get("dir"); err == nil {
			plugin = p.String()
		}
		if plugin == "" {
			plugin = os.Getenv("PLUGIN_MIGRATOR_PATH")
			if plugin == "" {
				return fmt.Errorf("plugin is empty")
			}
		}

		stat, err := os.Stat(plugin)
		if err != nil {
			return err
		}
		if !stat.IsDir() {
			return fmt.Errorf("%s is not dir", plugin)
		}
		return core.LoadPlugin(driver.String(), from.String(), fmt.Sprintf("%s/%s/migrate.so", plugin, version), core.Type_Up)
	case "down":
		for _, flag := range []string{"from", "driver", "v"} {
			if err := s.checkFlag(a, flag); err != nil {
				return err
			}
		}

		for _, flag := range []string{"from", "driver"} {
			if err := s.checkFlag(a, flag); err != nil {
				return err
			}
		}

		from, _ := a.Get("from")
		driver, _ := a.Get("driver")
		version, _ := a.Get("v")
		var plugin = ""
		if p, err := a.Get("dir"); err == nil {
			plugin = p.String()
		}
		if plugin == "" {
			plugin = os.Getenv("PLUGIN_MIGRATOR_PATH")
			if plugin == "" {
				return fmt.Errorf("plugin is empty")
			}
		}
		stat, err := os.Stat(plugin)
		if err != nil {
			return err
		}
		if !stat.IsDir() {
			return fmt.Errorf("%s is not dir", plugin)
		}
		return core.LoadPlugin(driver.String(), from.String(), fmt.Sprintf("%s/%s/migrate.so", plugin, version), core.Type_Down)
	case "show":
		for _, flag := range []string{"to", "driver", "v"} {
			if err := s.checkFlag(a, flag); err != nil {
				return err
			}
		}

		to, _ := a.Get("to")
		driver, _ := a.Get("driver")
		version, _ := a.Get("v")
		var plugin = ""
		if p, err := a.Get("plugin"); err == nil {
			plugin = p.String()
		}
		if plugin == "" {
			plugin = os.Getenv("PLUGIN_MIGRATOR_PATH")
			if plugin == "" {
				return fmt.Errorf("plugin is empty")
			}
		}
		stat, err := os.Stat(plugin)
		if err != nil {
			return err
		}
		if !stat.IsDir() {
			return fmt.Errorf("%s is not dir", plugin)
		}
		return core.Show(driver.String(), to.String(), fmt.Sprintf("%s/%s/migrate.so", plugin, version))
	case "make":
		return s._make(a)
	case "build":
		return s.build(a)
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
	case "build":
		buildHelp()
	default:
		return fmt.Errorf("mt[%s] unsupport", method)
	}

	return nil
}

func (s *serv) diff(a app.AppInterface) error {
	for _, flag := range []string{"from", "to", "fromdb", "todb", "driver"} {
		if err := s.checkFlag(a, flag); err != nil {
			return err
		}
	}

	driver, _ := a.Get("driver")
	if driver.String() != "mysql" {
		return fmt.Errorf("driver[%s] is not mysql", driver)
	}

	var dirVal = ""
	if dir, err := a.Get("dir"); err == nil {
		dirVal = dir.String()
	}
	if dirVal == "" {
		dirVal = os.Getenv("DIFF_SQL_PATH")
	}
	if dirVal == "" {
		return fmt.Errorf("dir is empty")
	}
	from, _ := a.Get("from")
	to, _ := a.Get("to")
	fromdb, _ := a.Get("fromdb")
	todb, _ := a.Get("todb")
	if err := s.mkdir(dirVal); err != nil {
		return err
	}

	ops, err := diff.Diff(context.Background(), driver.String(), from.String(), to.String(), fromdb.String(), todb.String())
	if err != nil {
		return err
	}

	file, err := os.Create(dirVal + fmt.Sprintf("/migrate_%d.sql", time.Now().UnixNano()))
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
	for _, flag := range []string{"to", "todb", "driver"} {
		if err := s.checkFlag(a, flag); err != nil {
			return err
		}
	}

	driver, _ := a.Get("driver")
	if driver.String() != "mysql" {
		return fmt.Errorf("driver[%s] is not mysql", driver)
	}

	var dirVal = ""
	if dir, err := a.Get("dir"); err == nil {
		dirVal = dir.String()
	}
	if dirVal == "" {
		dirVal = os.Getenv("DIFF_SQL_PATH")
		if dirVal == "" {
			return fmt.Errorf("dir is empty")
		}
	}
	stat, err := os.Stat(dirVal)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return fmt.Errorf("[%s] not dir", dirVal)
	}

	to, _ := a.Get("to")
	todb, _ := a.Get("todb")

	return diff.Migrate(driver.String(), to.String(), todb.String(), dirVal)
}

func (s *serv) _make(a app.AppInterface) error {
	for _, flag := range []string{"n", "v"} {
		if err := s.checkFlag(a, flag); err != nil {
			return err
		}
	}

	name, _ := a.Get("n")
	version, _ := a.Get("v")
	to, _ := a.Get("from")
	var dirVal = ""
	if dir, err := a.Get("dir"); err == nil {
		dirVal = dir.String()
	}
	if dirVal == "" {
		dirVal = os.Getenv("PLUGIN_MIGRATOR_PATH")
		if dirVal == "" {
			return fmt.Errorf("dir is empty")
		}
	}
	d, _ := a.Get("driver")
	return mk.Make(name.String(), version.String(), dirVal, to.String(), d.String())
}

func (s *serv) orm(a app.AppInterface) error {
	for _, flag := range []string{"from", "fromdb"} {
		if err := s.checkFlag(a, flag); err != nil {
			return err
		}
	}

	from, _ := a.Get("from")
	var dirVal = ""
	if dir, err := a.Get("dir"); err == nil {
		dirVal = dir.String()
	}
	if dirVal == "" {
		dirVal = os.Getenv("MODELS_PATH")
		if dirVal == "" {
			return fmt.Errorf("dir is empty")
		}
	}

	d, _ := a.Get("driver")
	db, _ := a.Get("fromdb")
	return orm.Orm(d.String(), from.String(), dirVal, db.String())
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
		s.Usage()
		return nil
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
	case "config":
		return s.config(a)
	case "help":
		return s.help(a)
	default:
		s.Usage()
		return nil
	}
}

func (s *serv) config(a app.AppInterface) error {
	if flag, err := a.Get("c"); err == nil && flag.IsInput() {
		_, err := os.Stat(".env")
		if !os.IsNotExist(err) {
			return nil
		}

		return os.WriteFile(".env", []byte(config_tpl), 0644)
	}

	if flag, err := a.Get("e"); err == nil && flag.IsInput() {
		cmd := exec.Command("vim", ".env")
		return cmd.Run()
	}

	if flag, err := a.Get("l"); err == nil && flag.IsInput() {
		content, err := os.ReadFile(".env")
		if err != nil {
			fmt.Println("no configs")
			return nil
		}

		fmt.Println(string(content))
		return nil
	}

	configHelp()
	return nil
}

func (s *serv) build(a app.AppInterface) error {
	for _, flag := range []string{"v"} {
		if err := s.checkFlag(a, flag); err != nil {
			return err
		}
	}

	version, _ := a.Get("v")
	var dirVal = ""
	if dir, err := a.Get("dir"); err == nil {
		dirVal = dir.String()
	}
	if dirVal == "" {
		dirVal = os.Getenv("PLUGIN_MIGRATOR_PATH")
		if dirVal == "" {
			return fmt.Errorf("dir is empty")
		}
	}

	stat, err := os.Stat(dirVal)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return fmt.Errorf("%s is not dir", dirVal)
	}

	cmd := exec.Command("go", "build", "-C", fmt.Sprintf("%s/%s", dirVal, version), "-buildmode=plugin", "-o", "migrate.so", "migrate.go")
	cmd.Stderr = os.Stdout
	return cmd.Run()
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
	case "config":
		configHelp()
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
