package serv

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/kovey/cli-go/app"
	"github.com/kovey/cli-go/gui"
	"github.com/kovey/db-go/ksql/core"
	"github.com/kovey/db-go/ksql/diff"
	"github.com/kovey/db-go/ksql/dir"
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
	app.GetHelp().Title = "ksql tools to manage sql migrate、create orm model"

	a.FlagArg("migrate", "migrate sql from dev to prodmigrate sql from dev to prod")
	a.FlagArg("diff", "diff table changed from dev to prod, create changed sql file")
	a.FlagArg("migplug", "migrate sql from migration plugins")
	a.FlagArg("orm", "create orm model from database")
	a.FlagArg("config", "manage config file")
	a.FlagArg("version", "show ksql version")

	a.FlagLong("dir", ".env.DIFF_SQL_PATH", app.TYPE_STRING, "sql directory", "migrate")
	a.FlagLong("driver", ".env.DB_DRIVER", app.TYPE_STRING, "database driver", "migrate")
	a.FlagLong("todb", ".env.TO_DB_NAME", app.TYPE_STRING, "database driver name", "migrate")
	a.FlagLong("to", ".env.TO_DB_*", app.TYPE_STRING, "to database dsn", "migrate")

	a.FlagLong("dir", ".env.DIFF_SQL_PATH", app.TYPE_STRING, "created sql directory", "diff")
	a.FlagLong("driver", ".env.DB_DRIVER", app.TYPE_STRING, "database driver", "diff")
	a.FlagLong("todb", ".env.TO_DB_NAME", app.TYPE_STRING, "database driver name", "diff")
	a.FlagLong("to", ".env.TO_DB_*", app.TYPE_STRING, "to database dsn", "diff")
	a.FlagLong("fromdb", ".env.DB_NAME", app.TYPE_STRING, "from database name", "diff")
	a.FlagLong("from", ".env.DB_*", app.TYPE_STRING, "from database dsn", "diff")

	a.FlagArg("up", "upgrade migrator", "migplug")
	a.FlagArg("down", "downgrade migrator", "migplug")
	a.FlagArg("show", "show migrator status", "migplug")
	a.FlagArg("make", "create upgrade and downgrade file", "migplug")
	a.FlagArg("build", "build migrator to ksql plugins", "migplug")

	a.FlagLong("dir", ".env.PLUGIN_MIGRATOR_PATH", app.TYPE_STRING, "migrators directory", "migplug", "up")
	a.FlagLong("driver", ".env.DB_DRIVER", app.TYPE_STRING, "database driver", "migplug", "up")
	a.FlagLong("from", ".env.DB_*", app.TYPE_STRING, "from database dsn", "migplug", "up")
	a.Flag("v", nil, app.TYPE_STRING, "migrator version", "migplug", "up")

	a.FlagLong("dir", ".env.PLUGIN_MIGRATOR_PATH", app.TYPE_STRING, "migrators directory", "migplug", "down")
	a.FlagLong("driver", ".env.DB_DRIVER", app.TYPE_STRING, "database driver", "migplug", "down")
	a.FlagLong("from", ".env.DB_*", app.TYPE_STRING, "from database dsn", "migplug", "down")
	a.Flag("v", nil, app.TYPE_STRING, "migrator version", "migplug", "down")

	a.FlagLong("dir", ".env.PLUGIN_MIGRATOR_PATH", app.TYPE_STRING, "migrators directory", "migplug", "show")
	a.FlagLong("driver", ".env.DB_DRIVER", app.TYPE_STRING, "database driver", "migplug", "show")
	a.FlagLong("from", ".env.DB_*", app.TYPE_STRING, "from database dsn", "migplug", "show")
	a.Flag("v", nil, app.TYPE_STRING, "migrator version", "migplug", "show")

	a.FlagLong("dir", ".env.PLUGIN_MIGRATOR_PATH", app.TYPE_STRING, "migrators directory", "migplug", "make")
	a.FlagLong("driver", ".env.DB_DRIVER", app.TYPE_STRING, "database driver", "migplug", "make")
	a.FlagLong("from", ".env.DB_*", app.TYPE_STRING, "from database dsn", "migplug", "make")
	a.FlagLong("fromdb", ".env.DB_NAME", app.TYPE_STRING, "from database name", "migplug", "make")
	a.Flag("n", nil, app.TYPE_STRING, "migrator name", "migplug", "make")
	a.Flag("v", nil, app.TYPE_STRING, "migrator version", "migplug", "make")

	a.FlagLong("dir", ".env.PLUGIN_MIGRATOR_PATH", app.TYPE_STRING, "migrators directory", "migplug", "build")
	a.Flag("v", nil, app.TYPE_STRING, "migrator version", "migplug", "build")

	a.FlagLong("dir", ".env.PLUGIN_MIGRATOR_PATH", app.TYPE_STRING, "migrators directory", "orm")
	a.FlagLong("driver", ".env.DB_DRIVER", app.TYPE_STRING, "database driver", "orm")
	a.FlagLong("from", ".env.DB_*", app.TYPE_STRING, "from database dsn", "orm")
	a.FlagLong("fromdb", ".env.DB_NAME", app.TYPE_STRING, "from database name", "orm")

	a.FlagNonValue("c", "create config", "config")
	a.FlagNonValue("e", "edit config", "config")
	a.FlagNonValue("l", "cat config", "config")
	return nil
}

func (s *serv) Init(app.AppInterface) error {
	return nil
}

func (s *serv) checkPlugin(plugin, version string) (string, error) {
	filePath := fmt.Sprintf("%s%s%s%smigrate.so", plugin, dir.Sep(), version, dir.Sep())
	stat, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("%s is not build, please use `ksql migplug build -v %s` build %s migrators", version, version, version)
		}

		return "", err
	}

	if stat.IsDir() {
		return "", fmt.Errorf("%s is not build, please use `ksql migplug build -v %s` build %s migrators", version, version, version)
	}

	return filePath, nil
}

func (s *serv) migplug(a app.AppInterface) error {
	method, err := a.Arg(1, app.TYPE_STRING)
	if err != nil {
		return err
	}

	switch method.String() {
	case "up":
		version, _ := a.Get("migplug", "up", "v")
		if !version.IsInput() {
			return fmt.Errorf("-v is empty")
		}

		var fromDsn = ""
		if from, _ := a.Get("migplug", "up", "from"); from.IsInput() {
			fromDsn = from.String()
		} else {
			fromDsn = s.getFromDsn()
		}

		var driverName = ""
		if driver, _ := a.Get("migplug", "up", "driver"); driver.IsInput() {
			driverName = driver.String()
		} else {
			driverName = os.Getenv("DB_DRIVER")
		}

		var plugin = ""
		if p, _ := a.Get("migplug", "up", "dir"); p.IsInput() {
			plugin = p.String()
		} else {
			plugin = os.Getenv("PLUGIN_MIGRATOR_PATH")
			if plugin == "" {
				return fmt.Errorf("plugin is empty")
			}
		}

		filePath, err := s.checkPlugin(plugin, version.String())
		if err != nil {
			return err
		}

		return core.LoadPlugin(driverName, fromDsn, filePath, core.Type_Up)
	case "down":
		version, _ := a.Get("migplug", "down", "v")
		if !version.IsInput() {
			return fmt.Errorf("-v is empty")
		}
		var fromDsn = ""
		if from, _ := a.Get("migplug", "down", "from"); from.IsInput() {
			fromDsn = from.String()
		} else {
			fromDsn = s.getFromDsn()
		}
		var driverName = ""
		if driver, _ := a.Get("migplug", "down", "driver"); driver.IsInput() {
			driverName = driver.String()
		} else {
			driverName = os.Getenv("DB_DRIVER")
		}
		var plugin = ""
		if p, _ := a.Get("migplug", "up", "dir"); p.IsInput() {
			plugin = p.String()
		} else {
			plugin = os.Getenv("PLUGIN_MIGRATOR_PATH")
			if plugin == "" {
				return fmt.Errorf("plugin is empty")
			}
		}
		filePath, err := s.checkPlugin(plugin, version.String())
		if err != nil {
			return err
		}
		return core.LoadPlugin(driverName, fromDsn, filePath, core.Type_Down)
	case "show":
		version, _ := a.Get("migplug", "show", "v")
		if !version.IsInput() {
			return fmt.Errorf("-v is empty")
		}
		var fromDsn = ""
		if from, _ := a.Get("migplug", "show", "from"); from.IsInput() {
			fromDsn = from.String()
		} else {
			fromDsn = s.getFromDsn()
		}
		var driverName = ""
		if driver, _ := a.Get("migplug", "show", "driver"); driver.IsInput() {
			driverName = driver.String()
		} else {
			driverName = os.Getenv("DB_DRIVER")
		}
		var plugin = ""
		if p, _ := a.Get("migplug", "show", "plugin"); p.IsInput() {
			plugin = p.String()
		} else {
			plugin = os.Getenv("PLUGIN_MIGRATOR_PATH")
			if plugin == "" {
				return fmt.Errorf("plugin is empty")
			}
		}
		filePath, err := s.checkPlugin(plugin, version.String())
		if err != nil {
			return err
		}
		return core.Show(driverName, fromDsn, filePath)
	case "make":
		return s._make(a)
	case "build":
		return s.build(a)
	default:
		return nil
	}
}

func (s *serv) diff(a app.AppInterface) error {
	driverName := ""
	if driver, _ := a.Get("diff", "driver"); driver.IsInput() {
		driverName = driver.String()
	}
	if driverName != "mysql" {
		return fmt.Errorf("driver[%s] is not mysql", driverName)
	}

	var dirVal = ""
	if dir, _ := a.Get("diff", "dir"); dir.IsInput() {
		dirVal = dir.String()
	} else {
		dirVal = os.Getenv("DIFF_SQL_PATH")
	}
	if dirVal == "" {
		return fmt.Errorf("dir is empty")
	}
	fromDsn := ""
	toDsn := ""
	fromDbName := ""
	toDbName := ""
	if from, _ := a.Get("diff", "from"); from.IsInput() {
		fromDsn = from.String()
	} else {
		fromDsn = s.getFromDsn()
	}
	if to, _ := a.Get("diff", "to"); to.IsInput() {
		toDsn = to.String()
	} else {
		toDsn = s.getToDsn()
	}

	if fromdb, _ := a.Get("diff", "fromdb"); fromdb.IsInput() {
		fromDbName = fromdb.String()
	} else {
		fromDbName = os.Getenv("DB_NAME")
	}
	if todb, _ := a.Get("diff", "todb"); todb.IsInput() {
		toDbName = todb.String()
	} else {
		toDbName = os.Getenv("TO_DB_NAME")
	}

	if err := s.mkdir(dirVal); err != nil {
		return err
	}

	ops, err := diff.Diff(context.Background(), driverName, fromDsn, toDsn, fromDbName, toDbName)
	if err != nil {
		return err
	}

	file, err := os.Create(dirVal + fmt.Sprintf("%smigrate_%d.sql", dir.Sep(), time.Now().UnixNano()))
	if err != nil {
		return err
	}

	defer file.Close()
	buff := bufio.NewWriter(file)
	buff.WriteString(fmt.Sprintf("-- from database: %s\n", fromDbName))
	buff.WriteString(fmt.Sprintf("-- to database:   %s\n", toDbName))
	buff.WriteString(fmt.Sprintf("-- tool version:  %s\n", version.Version()))
	buff.WriteString(fmt.Sprintf("-- create time:   %s", time.Now().Format(time.DateTime)))
	for _, op := range ops {
		if _, err := buff.WriteString("\n"); err != nil {
			debug.Erro("write error: %s", err)
		}

		if _, err := buff.WriteString(op.Prepare()); err != nil {
			debug.Erro("write error: %s", err)
		}
		if _, err := buff.WriteString(";"); err != nil {
			debug.Erro("write error: %s", err)
		}
	}

	return buff.Flush()
}

func (s *serv) migrate(a app.AppInterface) error {
	driverName := ""
	if driver, _ := a.Get("migrate", "driver"); driver.IsInput() {
		driverName = driver.String()
	} else {
		driverName = os.Getenv("DB_DRIVER")
	}
	if driverName != "mysql" {
		return fmt.Errorf("driver[%s] is not mysql", driverName)
	}

	var dirVal = ""
	if dir, _ := a.Get("migrate", "dir"); dir.IsInput() {
		dirVal = dir.String()
	} else {
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

	toDsn := ""
	toDbName := ""
	if to, _ := a.Get("migrate", "to"); to.IsInput() {
		toDsn = to.String()
	} else {
		toDsn = s.getToDsn()
	}

	if todb, _ := a.Get("migrate", "todb"); todb.IsInput() {
		toDbName = todb.String()
	} else {
		toDbName = os.Getenv("TO_DB_NAME")
	}

	return diff.Migrate(driverName, toDsn, toDbName, dirVal)
}

func (s *serv) _make(a app.AppInterface) error {
	name, _ := a.Get("migplug", "make", "n")
	if !name.IsInput() {
		return fmt.Errorf("-n is empty")
	}

	version, _ := a.Get("migplug", "make", "v")
	if !version.IsInput() {
		return fmt.Errorf("-v is empty")
	}
	var fromDsn = ""
	if from, _ := a.Get("migplug", "make", "from"); from.IsInput() {
		fromDsn = from.String()
	} else {
		fromDsn = s.getFromDsn()
	}
	var dirVal = ""
	if dir, _ := a.Get("migplug", "make", "dir"); dir.IsInput() {
		dirVal = dir.String()
	}
	if dirVal == "" {
		dirVal = os.Getenv("PLUGIN_MIGRATOR_PATH")
		if dirVal == "" {
			return fmt.Errorf("dir is empty")
		}
	}
	var driverName = ""
	if driver, _ := a.Get("migplug", "make", "driver"); driver.IsInput() {
		driverName = driver.String()
	} else {
		driverName = os.Getenv("DB_DRIVER")
	}
	return mk.Make(name.String(), version.String(), dirVal, fromDsn, driverName)
}

func (s *serv) orm(a app.AppInterface) error {
	fromDsn := ""
	if from, _ := a.Get("orm", "from"); from.IsInput() {
		fromDsn = from.String()
	} else {
		fromDsn = s.getFromDsn()
	}
	var dirVal = ""
	if dir, _ := a.Get("orm", "dir"); dir.IsInput() {
		dirVal = dir.String()
	} else {
		dirVal = os.Getenv("MODELS_PATH")
		if dirVal == "" {
			return fmt.Errorf("dir is empty")
		}
	}

	driverName := ""
	if d, _ := a.Get("orm", "driver"); d.IsInput() {
		driverName = d.String()
	} else {
		driverName = os.Getenv("DB_DRIVER")
	}
	dbName := ""
	if db, _ := a.Get("orm", "fromdb"); db.IsInput() {
		dbName = db.String()
	} else {
		dbName = os.Getenv("DB_NAME")
	}
	return orm.Orm(driverName, fromDsn, dirVal, dbName)
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
	default:
		s.Usage()
		return nil
	}
}

func (s *serv) config(a app.AppInterface) error {
	if flag, err := a.Get("config", "c"); err == nil && flag.IsInput() {
		_, err := os.Stat(".env")
		if !os.IsNotExist(err) {
			return nil
		}

		return os.WriteFile(".env", []byte(config_tpl), 0644)
	}

	if flag, err := a.Get("config", "e"); err == nil && flag.IsInput() {
		commands := []string{"vim", "vi", "nvim", "emacs", "gedit"}
		for _, command := range commands {
			if _, err := exec.LookPath(command); err == nil {
				cmd := exec.Command("vim", ".env")
				cmd.Stdout = os.Stdout
				cmd.Stdin = os.Stdin
				return cmd.Run()
			}
		}

		return fmt.Errorf("editor not install, please install one of [%s]", strings.Join(commands, ","))
	}

	if flag, err := a.Get("config", "l"); err == nil && flag.IsInput() {
		content, err := os.ReadFile(".env")
		if err != nil {
			fmt.Println("no configs")
			return nil
		}

		fmt.Println(string(content))
		return nil
	}

	return nil
}

func (s *serv) build(a app.AppInterface) error {
	version, _ := a.Get("migplug", "build", "v")
	if !version.IsInput() {
		return fmt.Errorf("-v is empty")
	}

	var dirVal = ""
	if dir, err := a.Get("migplug", "build", "dir"); err == nil && dir.IsInput() {
		dirVal = dir.String()
	} else {
		dirVal = os.Getenv("PLUGIN_MIGRATOR_PATH")
		if dirVal == "" {
			return fmt.Errorf("dir is empty")
		}
	}

	path := fmt.Sprintf("%s%s%s", dirVal, dir.Sep(), version)
	stat, err := os.Stat(path + dir.Sep() + "migrate.go")
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s migrators not created, please use `ksql migplug make -v %s -n xxx` create migrator", version, version)
		}
		return err
	}

	if stat.IsDir() {
		return fmt.Errorf("%s is not file", path)
	}

	cmd := exec.Command("go", "build", "-C", path, "-buildmode=plugin", "-o", "migrate.so", "migrate.go")
	debug.Info("%s migrator build begin, please wait...", version)
	defer debug.Info("%s migrator build end.", version)
	return cmd.Run()
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
