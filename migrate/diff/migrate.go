package diff

import (
	"bufio"
	"context"
	"io"
	"os"

	"github.com/kovey/db-go/v3/sql"
	"github.com/kovey/debug-go/debug"
)

func Migrate(driverName, dsn, dbname, path string) error {
	conn, err := getConn(driverName, dsn)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	debug.Info("migrate from [%s] to [%s] begin...", path, dbname)
	defer debug.Info("migrate from [%s] to [%s] end.", path, dbname)
	ctx := context.Background()
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		f, err := os.Open(path + "/" + file.Name())
		if err != nil {
			return err
		}

		defer f.Close()
		debug.Info("migrate file[%s/%s] begin...", path, file.Name())
		defer debug.Info("migrate file[%s/%s] end.", path, file.Name())
		buf := bufio.NewReader(f)
		for {
			line, err := buf.ReadString('\n')
			if err == io.EOF {
				if _, err := conn.ExecRaw(ctx, sql.Raw(line)); err != nil {
					return err
				}

				break
			}

			if err != nil {
				break
			}

			if _, err := conn.ExecRaw(ctx, sql.Raw(line)); err != nil {
				return err
			}
		}
	}

	return nil
}