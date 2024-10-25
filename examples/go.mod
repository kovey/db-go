module github.com/kovey/db-go/v3/examples

go 1.22.3

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/kovey/db-go/v3 v3.0.0
	github.com/kovey/db-go/v3/migplug v0.0.0
)

require filippo.io/edwards25519 v1.1.0 // indirect

replace github.com/kovey/db-go/v3 v3.0.0 => ../

replace github.com/kovey/db-go/v3/migplug v0.0.0 => ../migplug
