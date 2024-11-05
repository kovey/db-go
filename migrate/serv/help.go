package serv

import "fmt"

func migrateHelp() {
	fmt.Println(`
Usage:
	ksql-tool migrate [-dir] [-driver] [-todb] [-to]
		-dir     sql directory
		-driver  database driver(mysql)
		-todb    database name
		-to      database dsn
		`)
}

func diffHelp() {
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
}

func migplugHelp() {
	// "p", "mt", "to", "d"
	fmt.Println(`
Usage:
	ksql-tool diff [-plugin] [-driver] [-m] [-to]
		-plugin  plugin directory
		-driver  database driver(mysql)
		-m       plugin method(up|down|show|make)
		-to      to database dsn
		`)
}

func ormHelp() {
	// "to", "dir", "d", "todb"
	fmt.Println(`
Usage:
	ksql-tool orm [-driver] [-to] [-todb] [-dir]
		-driver  database driver(mysql)
		-to      to database dsn
		-todb    to database name
		-dir     orm model directory
		`)
}
