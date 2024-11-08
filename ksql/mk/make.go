package mk

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"time"

	"github.com/kovey/db-go/ksql/core"
	"github.com/kovey/db-go/ksql/mk/template"
	v "github.com/kovey/db-go/ksql/version"
	"github.com/kovey/db-go/v3/db"
	"github.com/kovey/debug-go/debug"
)

var Err_Version_Exists = errors.New("Version Exists")
var Err_Parse_Error = errors.New("Parse Error")

func Make(name, version, dir, dsn, driverName string) error {
	path := dir + "/" + version + "/migrations"
	if err := os.MkdirAll(path, 0755); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	if err := db.Init(db.Config{DriverName: driverName, DataSourceName: dsn, MaxIdleTime: 120 * time.Second, MaxLifeTime: 120 * time.Second, MaxIdleConns: 10, MaxOpenConns: 10}); err != nil {
		return err
	}

	t := &template.MigrateTemplate{Name: name, Package: "migrations", Id: uint64(time.Now().UnixNano()), Version: version, CreateTime: time.Now().Format(time.DateTime), ToolVersion: v.Version()}
	if ok, err := core.Has(context.Background(), t.Id); err != nil {
		return err
	} else if ok {
		return Err_Version_Exists
	}

	mainT := getTemplate(path, version, getFullPackage(dir))
	if mainT == nil {
		return Err_Parse_Error
	}

	if mainT.Has(name) {
		return Err_Version_Exists
	}

	mainT.Migrates = append(mainT.Migrates, name)
	mainRes, err := mainT.Parse()
	if err != nil {
		return err
	}
	res, err := t.Parse()
	if err != nil {
		return err
	}

	if err := os.WriteFile(path+"/"+name+".go", res, 0644); err != nil {
		return err
	}

	return os.WriteFile(dir+"/"+version+"/migrate.go", mainRes, 0644)
}

func getFullPackage(dir string) string {
	dirInfo := strings.Split(dir, "/")
	dir += "/.."
	sub := 1
	prefix := dirInfo[len(dirInfo)-sub]
	for {
		files, err := os.ReadDir(dir)
		if err != nil {
			debug.Erro("read dir[%s] error: %s", dir, err)
			break
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			if file.Name() != "go.mod" {
				continue
			}

			pack := readFirst(dir + "/" + file.Name())
			if prefix == "" {
				return pack
			}

			return pack + "/" + prefix
		}

		sub++
		dir += "/.."
		prefix = dirInfo[len(dirInfo)-sub] + "/" + prefix
	}

	return ""
}

func readFirst(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return ""
	}

	defer file.Close()
	buf := bufio.NewReader(file)
	line, err := buf.ReadString('\n')
	if err == nil || err == io.EOF {
		return strings.ReplaceAll(strings.Trim(line, "\r\n\t "), "module ", "")
	}

	return ""
}

func getTemplate(path, version, fullPackage string) *template.MainTpl {
	files, err := os.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &template.MainTpl{}
		}

		return nil
	}

	tmp := &template.MainTpl{CreateTime: time.Now().Format(time.DateTime), Version: v.Version()}
	tmp.Imports = append(tmp.Imports, fullPackage+"/"+version+"/migrations")
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := strings.Split(file.Name(), ".")[0]
		if tmp.Has(name) {
			continue
		}

		tmp.Migrates = append(tmp.Migrates, name)
	}

	return tmp
}
