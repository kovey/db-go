package serv

import "fmt"

func migrateHelp() {
	fmt.Println(`
Usage:
	ksql migrate [--dir] [--driver] [--todb] [--to]
		--dir     sql directory, if non, use .env.DIFF_SQL_PATH
		--driver  database driver, if non, use .env.DB_DRIVER
		--todb    database name, if non, use .env.TO_DB_NAME
		--to      database dsn, if non, use .env.TO_DB_*
		`)
}

func diffHelp() {
	fmt.Println(`
Usage:
	ksql diff [--dir] [--driver] [--fromdb] [--from] [--todb] [--to]
		--dir     created sql directory, if non, use .env.DIFF_SQL_PATH
		--driver  database driver, if non, use .env.DB_DRIVER
		--todb    database name, if non, use .env.TO_DB_NAME
		--to      database dsn, if non, use .env.TO_DB_*
		--fromdb  from database name, if non, use .env.DB_NAME
		--from    from database dsn, if non, use .env.DB_*
		`)
}

func migplugHelp() {
	fmt.Println(`
Usage:
	ksql migplug [commands] [arguments]

The commands are:
	up    upgrade migrator
	down  downgrade migrator
	show  show migrator status
	make  create upgrade and downgrade file 
	build build migrator to ksql plugins

Use "ksql migplug help <command>" for more information about a command.
		`)
}

func buildHelp() {
	fmt.Println(`
Usage:
	ksql migplug build [--dir] [-v]
		--dir     migrators directory, if non, use .env.PLUGIN_MIGRATOR_PATH
		-v        migrator version
	`)
}

func makeHelp() {
	fmt.Println(`
Usage:
	ksql migplug make [--dir] [--driver] [--from] [--fromdb] [-n] [-v]
		--dir     make migrator to directory, if non, use .env.PLUGIN_MIGRATOR_PATH
		--driver  database driver, if non, use .env.DB_DRIVER
		--from    database dsn, if non, use .env.DB_*
		--fromdb  database name, if non, use .env.DB_NAME
		-n        migrator name
		-v        migrator version
	`)
}

func upHelp() {
	fmt.Println(`
Usage:
	ksql migplug up [--dir] [--driver] [--from] [-v]
		--dir     migrators directory, if non, use .env.PLUGIN_MIGRATOR_PATH
		--driver  database driver, if non, use .env.DB_DRIVER
		--from    database dsn, if non, use .env.DB_*
		-v        migrator version
	`)
}

func downHelp() {
	fmt.Println(`
Usage:
	ksql migplug down [--dir] [--driver] [--from] [-v]
		--dir     migrators directory, if non, use .env.PLUGIN_MIGRATOR_PATH
		--driver  database driver, if non, use .env.DB_DRIVER
		--from    database dsn, if non, use .env.DB_*
		-v        migrator version
	`)
}

func showHelp() {
	fmt.Println(`
Usage:
	ksql migplug show [--dir] [--driver] [--from] [-v]
		--dir     migrators directory, if non, use .env.PLUGIN_MIGRATOR_PATH
		--driver  database driver, if non, use .env.DB_DRIVER
		--from    database dsn, if non, use .env.DB_*
		-v        migrator version
	`)
}

func ormHelp() {
	fmt.Println(`
Usage:
	ksql orm [--driver] [--from] [--fromdb] [--dir]
		--driver  database driver, if non, use .env.DB_DRIVER
		--from    database dsn, if non, use .env.DB_*
		--fromdb  database name, if non, use .env.DB_NAME
		--dir     orm model directory, if non, use .env.MODELS_PATH
		`)
}

func configHelp() {
	fmt.Println(`
Usage:
	ksql config [options] 

options:
		-c  create .env config file
		-e  edit .env config file
		-l  view .env config file
		`)
}
