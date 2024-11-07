package serv

import "fmt"

func migrateHelp() {
	fmt.Println(`
Usage:
	ksql migrate [-dir] [-driver] [-todb] [-to]
		--dir     sql directory
		--driver  database driver(mysql)
		--todb    database name
		--to      database dsn
		`)
}

func diffHelp() {
	// "from", "to", "fromdb", "todb", "d", "dir"
	fmt.Println(`
Usage:
	ksql diff [-dir] [-driver] [-fromdb] [-from] [-todb] [-to]
		--dir     created sql directory
		--driver  database driver(mysql)
		--todb    to database name
		--to      to database dsn
		--fromdb  from database name
		--from    from database dsn
		`)
}

func migplugHelp() {
	// "p", "mt", "to", "d"
	// "n", "v", "to", "dir", "d"
	fmt.Println(`
Usage:
	ksql migplug [commands] [arguments]

The commands are:
	up	  upgrade migrator
	down  downgrade migrator
	show  show migrator status
	make  create upgrade and downgrade file 
	help  show help info
	
		`)
}

func makeHelp() {
	fmt.Println(`
Usage:
	ksql migplug make [-dir] [-driver] [-to] [-n] [-v]
		--dir     make migrator to directory
		--driver  database driver(mysql)
		--to      to database dsn
		-n        migrator name
		-v        migrator version
	`)
}

func upHelp() {
	fmt.Println(`
Usage:
	ksql migplug up [-plugin] [-driver] [-to]
		--plugin  plugin directory
		--driver  database driver(mysql)
		--to      to database dsn
	`)
}

func downHelp() {
	fmt.Println(`
Usage:
	ksql migplug down [-plugin] [-driver] [-to]
		--plugin  plugin directory
		--driver  database driver(mysql)
		--to      to database dsn
	`)
}

func showHelp() {
	fmt.Println(`
Usage:
	ksql migplug show [-plugin] [-driver] [-to]
		--plugin  plugin directory
		--driver  database driver(mysql)
		--to      to database dsn
	`)
}

func ormHelp() {
	// "to", "dir", "d", "todb"
	fmt.Println(`
Usage:
	ksql orm [-driver] [-to] [-todb] [-dir]
		--driver  database driver(mysql)
		--to      to database dsn
		--todb    to database name
		--dir     orm model directory
		`)
}
