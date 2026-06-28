# db-go

A high-performance, zero-reflection SQL toolkit for Go — designed for MySQL with a fluent, type-safe query builder, model-based CRUD, sharding, and transaction support.

## Philosophy

Unlike typical Go ORMs, **db-go avoids reflection entirely**. You define your own `Columns()` and `Values()` methods, giving you full control over column-to-field mapping while eliminating the runtime cost and opacity of struct-tag parsing. The result is predictable performance, straightforward debugging, and no surprises.

## Installation

```bash
go get -u github.com/kovey/db-go/v3
```

This library targets **MySQL**. Import your preferred MySQL driver (e.g. `github.com/go-sql-driver/mysql`) in your `main` package.

## Quick Start

### 1. Initialize the connection

```go
import (
    _ "github.com/go-sql-driver/mysql"
    "github.com/kovey/db-go/v3/db"
)

func init() {
    if err := db.Init(db.Config{
        DriverName:     "mysql",
        DataSourceName: "user:password@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=true",
        MaxIdleTime:    60 * time.Second,
        MaxLifeTime:    120 * time.Second,
        MaxIdleConns:   10,
        MaxOpenConns:   50,
        LogOpened:      true,  // enable SQL logging
        LogMax:          2048, // log buffer size
    }); err != nil {
        panic(err)
    }
    defer db.Close()
}
```

### 2. Define a model

```go
type User struct {
    *model.Model
    Id         int64     `json:"id"`
    Account    string    `json:"account"`
    Nickname   string    `json:"nickname"`
    Status     int       `json:"status"`
    CreateTime int64     `json:"create_time"`
    UpdateTime int64     `json:"update_time"`
}

func NewUser() *User {
    return &User{Model: model.NewModel("user", "id", model.Type_Int)}
}

// Clone returns a fresh instance (used internally for row scanning).
func (u *User) Clone() ksql.RowInterface {
    return NewUser()
}

// Columns declares the database column names in the same order as Values.
func (u *User) Columns() []string {
    return []string{"id", "account", "nickname", "status", "create_time", "update_time"}
}

// Values returns pointers to the struct fields for scanning.
func (u *User) Values() []any {
    return []any{&u.Id, &u.Account, &u.Nickname, &u.Status, &u.CreateTime, &u.UpdateTime}
}

// Save persists the model (insert or update).
func (u *User) Save(ctx context.Context) error {
    return u.Model.SaveBy(ctx, u)
}

// Delete removes the model from the database.
func (u *User) Delete(ctx context.Context) error {
    return u.Model.DeleteBy(ctx, u)
}
```

> **Primary key types**: use `model.Type_Int` for integer keys or `model.Type_Str` for string keys. Integer auto-increment is the default; call `m.NoAutoInc()` to disable it.

## CRUD Operations

### Create

```go
u := NewUser()
u.Account = "alice"
u.Nickname = "Alice"
u.Status = 0
u.CreateTime = time.Now().Unix()
u.UpdateTime = u.CreateTime

if err := u.Save(context.Background()); err != nil {
    // handle error
}
// u.Id is now set to the auto-increment value
```

Non-auto-increment keys are also supported:

```go
u := NewUser()
u.NoAutoInc()
u.Id = 42
u.Account = "bob"
// ...
u.Save(ctx)
```

### Read (single row)

```go
u := NewUser()
if err := db.Model(u).Where("id", ksql.Eq, 1).First(ctx); err != nil {
    // handle error
}
fmt.Println(u.Account)
```

Use `Find` when querying by primary key:

```go
u := NewUser()
db.Find(ctx, u, 1)  // SELECT ... WHERE id = 1
```

Locking rows inside a transaction:

```go
db.Lock(ctx, conn, u, 1)       // FOR UPDATE
db.LockShare(ctx, conn, u, 1)  // FOR SHARE
```

### Read (multiple rows)

```go
var users []*User
if err := db.Models(&users).Where("status", ksql.Eq, 0).All(ctx); err != nil {
    // handle error
}
```

### Read (raw structs without Model)

```go
type SimpleRow struct {
    *db.Row
    Id   int64
    Name string
}
// implement Clone, Columns, Values...
var rows []*SimpleRow
db.Rows(&rows).Table("user").Columns("id", "name").All(ctx)
```

### Update

```go
// Fetch then modify
u := NewUser()
db.Model(u).Where("id", ksql.Eq, 1).First(ctx)
u.Nickname = "NewNick"
u.UpdateTime = time.Now().Unix()
u.Save(ctx) // issues UPDATE
```

Lifecycle hooks are available for custom logic:

```go
func (u *User) OnUpdateBefore(conn ksql.ConnectionInterface) error { /* before UPDATE */ }
func (u *User) OnUpdateAfter(conn ksql.ConnectionInterface) error  { /* after UPDATE */ }
func (u *User) OnCreateBefore(conn ksql.ConnectionInterface) error { /* before INSERT */ }
func (u *User) OnCreateAfter(conn ksql.ConnectionInterface) error  { /* after INSERT */ }
func (u *User) OnDeleteBefore(conn ksql.ConnectionInterface) error { /* before DELETE */ }
func (u *User) OnDeleteAfter(conn ksql.ConnectionInterface) error  { /* after DELETE */ }
```

### Delete

```go
u := NewUser()
db.Model(u).Where("id", ksql.Eq, 1).First(ctx)
u.Delete(ctx)
```

### Convenience helpers

```go
data := db.NewData()
data.Set("nickname", "NewName")
data.Set("update_time", time.Now().Unix())

db.Insert(ctx, "user", data)
db.Update(ctx, "user", data, db.NewWhere().Where("id", ksql.Eq, 1))
db.Delete(ctx, "user", db.NewWhere().Where("id", ksql.Eq, 1))
```

## Query Builder

The query builder supports a fluent, chainable API with full SQL coverage.

### Operators

```go
ksql.Eq   // =
ksql.Neq  // <>
ksql.Gt   // >
ksql.Ge   // >=
ksql.Lt   // <
ksql.Le   // <=
ksql.Like // LIKE
```

### Filtering

```go
db.Models(&users).
    Where("status", ksql.Eq, 0).
    WhereIsNull("deleted_at").
    WhereIsNotNull("email").
    WhereIn("role", []any{"admin", "moderator"}).
    WhereNotIn("id", []any{1, 2, 3}).
    Between("created_at", startTime, endTime).
    NotBetween("score", 0, 50).
    All(ctx)
```

### Subquery filtering

```go
sub := db.NewQuery().Table("orders").Columns("user_id").Where("amount", ksql.Gt, 100)

db.Models(&users).
    WhereInBy("id", sub).     // WHERE id IN (SELECT ...)
    WhereNotInBy("id", sub).  // WHERE id NOT IN (SELECT ...)
    All(ctx)
```

### AND / OR grouping

```go
db.Models(&users).
    Where("status", ksql.Eq, 0).
    AndWhere(func(w ksql.WhereInterface) {
        w.Where("age", ksql.Ge, 18).Where("age", ksql.Le, 65)
    }).
    OrWhere(func(w ksql.WhereInterface) {
        w.Where("vip", ksql.Eq, 1).Where("score", ksql.Gt, 1000)
    }).
    All(ctx)
// WHERE status = 0 AND (age >= 18 AND age <= 65) OR (vip = 1 AND score > 1000)
```

### Raw WHERE expressions

```go
db.Models(&users).WhereExpress(
    db.Raw("JSON_EXTRACT(meta, ?) = ?", "$.verified", true),
).All(ctx)
```

### HAVING

```go
db.Models(&users).
    Group("role").
    Having("cnt", ksql.Gt, 10).
    HavingIsNull("deleted_at").
    HavingBetween("avg_score", 60, 100).
    All(ctx)
```

### Subquery in FROM

```go
sub := db.NewQuery().Table("orders").Columns("user_id", "COUNT(1) as order_count").Group("user_id")
var rows []*UserOrderCount
db.Rows(&rows).TableBy(sub, "o").Columns("u.id", "u.name", "o.order_count").
    LeftJoin("user").As("u").On("u.id", "=", "o.user_id").
    All(ctx)
```

### Joins

```go
db.Rows(&rows).Table("user").As("u").
    Join("profile").As("p").On("p.user_id", "=", "u.id").            // INNER JOIN
    LeftJoin("settings").As("s").On("s.user_id", "=", "u.id").       // LEFT JOIN
    RightJoin("audit_log").As("a").On("a.user_id", "=", "u.id").     // RIGHT JOIN
    JoinExpress(db.Raw("LEFT JOIN extra AS e ON e.id = u.id")).       // raw join
    All(ctx)
```

### Ordering & grouping

```go
db.Models(&users).
    Order("created_at").          // ASC
    OrderDesc("score").            // DESC
    Group("role", "status").
    GroupWithRollUp().
    All(ctx)
```

### Pagination

```go
pageInfo, err := db.Models(&[]*User{}).
    Where("status", ksql.Eq, 0).
    Pagination(ctx, 1, 20) // page 1, 20 per page

fmt.Println(pageInfo.TotalCount()) // 157
fmt.Println(pageInfo.TotalPage())  // 8
for _, u := range pageInfo.List() {
    // ...
}
```

### Aggregate functions

```go
// COUNT
count, _ := db.Models(&[]*User{}).Where("status", ksql.Eq, 0).Count(ctx)

// SUM
total, _ := db.Rows(&users).Table("orders").SumFloat(ctx, "amount") // float64
total, _ := db.Rows(&users).Table("orders").SumInt(ctx, "amount")   // uint64

// MAX / MIN
u := NewUser()
db.Model(u).Where("status", ksql.Eq, 0).Max(ctx, "score")
db.Model(u).Where("status", ksql.Eq, 0).Min(ctx, "score")
```

### EXISTS

```go
exists, _ := db.Models(&[]*User{}).Where("account", ksql.Eq, "alice").Exist(ctx)
```

### DISTINCT

```go
db.Models(&users).Distinct().Columns("role").All(ctx)
// SELECT DISTINCT role FROM user
```

### Locking clauses

```go
db.Model(u).Where("id", ksql.Eq, 1).For().Update().First(ctx)   // FOR UPDATE
db.Model(u).Where("id", ksql.Eq, 1).For().Share().First(ctx)    // FOR SHARE
db.Model(u).Where("id", ksql.Eq, 1).For().Update().NoWait().First(ctx)   // FOR UPDATE NOWAIT
db.Model(u).Where("id", ksql.Eq, 1).For().Update().SkipLocked().First(ctx)
db.Model(u).ForUpdate().Where("id", ksql.Eq, 1).First(ctx)      // shorthand
```

### Query modifiers

```go
db.Models(&users).HighPriority().StraightJoin().SqlSmallResult().SqlCalcFoundRows().All(ctx)
db.Models(&users).SqlBigResult().SqlBufferResult().SqlNoCache().All(ctx)
db.Models(&users).DistinctRow().All(ctx)
```

## Table Management

### Create a table

```go
err := db.Create(ctx, "user", func(table ksql.TableInterface) {
    table.AddBigInt("id").AutoIncrement().Unsigned().Comment("Primary key")
    table.AddString("account", 63).NotNullable().Default("").Comment("Account")
    table.AddString("nickname", 63).Default("").Comment("Nickname")
    table.AddTinyInt("status").Default("0").Comment("0 - active, 1 - banned")
    table.AddBigInt("create_time").Unsigned().Default("0")
    table.AddBigInt("update_time").Unsigned().Default("0")
    table.AddUnique("idx_account", "account")
    table.AddIndex("idx_status").Column("status", 0, ksql.Order_None)
    table.Engine("InnoDB").Charset("utf8mb4").Collate("utf8mb4_0900_ai_ci").Comment("User table")
})
```

### Alter a table

```go
err := db.Table(ctx, "user", func(table ksql.TableInterface) {
    table.AddString("email", 127).Default("").Comment("Email")
    table.ChangeColumn("nickname", "nick", "varchar", 31, 0).Default("").Comment("New nickname")
    table.DropColumn("old_column")
    table.DropColumnIfExists("optional_column")
    table.DropIndex("idx_old")
    table.AddUnique("idx_email", "email")
    table.Engine("InnoDB")
})
```

### Drop a table

```go
db.DropTable(ctx, "user")           // DROP TABLE user
db.DropTableIfExists(ctx, "user")  // DROP TABLE IF EXISTS user
```

### Check existence

```go
db.HasTable(ctx, "user")
db.HasColumn(ctx, "user", "account")
db.HasIndex(ctx, "user", "idx_account")
```

### Show DDL

```go
ddl, err := db.ShowDDL(ctx, "user")
```

### Column type helpers

| Method | SQL Type |
|---|---|
| `AddInt` | `INT` |
| `AddTinyInt` | `TINYINT` |
| `AddSmallInt` | `SMALLINT` |
| `AddBigInt` | `BIGINT` |
| `AddString(length)` | `VARCHAR(length)` |
| `AddChar(length)` | `CHAR(length)` |
| `AddText` | `TEXT` |
| `AddBlob` | `BLOB` |
| `AddDecimal(length, scale)` | `DECIMAL(length, scale)` |
| `AddDouble(length, scale)` | `DOUBLE(length, scale)` |
| `AddFloat(length, scale)` | `FLOAT(length, scale)` |
| `AddDate` | `DATE` |
| `AddDateTime` | `DATETIME` |
| `AddTimestamp` | `TIMESTAMP` |
| `AddEnum(options)` | `ENUM(...)` |
| `AddSet(sets)` | `SET(...)` |
| `AddBinary(length)` | `BINARY(length)` |
| `AddGeoMetry` | `GEOMETRY` |
| `AddPoint` | `POINT` |
| `AddPolygon` | `POLYGON` |
| `AddLineString` | `LINESTRING` |

## Raw SQL

For queries that don't fit the builder pattern:

```go
// Execute
result, err := db.ExecRaw(ctx, db.Raw("TRUNCATE TABLE user"))

// Query rows
var users []*User
db.QueryRaw(ctx, db.Raw("SELECT * FROM user WHERE status = ?", 0), &users)

// Query single row
u := NewUser()
db.QueryRowRaw(ctx, db.Raw("SELECT * FROM user WHERE id = ?", 1), u)

// Insert / Update / Delete (with type-checking)
db.InsertRaw(ctx, db.Raw("INSERT INTO user (name) VALUES (?)", "alice"))
db.UpdateRaw(ctx, db.Raw("UPDATE user SET status = ? WHERE id = ?", 1, 42))
db.DeleteRaw(ctx, db.Raw("DELETE FROM user WHERE status = ?", 2))

// Scan scalars
var name string
db.Scan(ctx, db.Raw("SELECT name FROM user WHERE id = ?", 1), &name)
```

## Transactions

```go
err := db.Transaction(ctx, func(ctx context.Context, conn ksql.ConnectionInterface) error {
    // Use conn for all operations within the transaction
    u := NewUser()
    u.Account = "alice"
    if err := u.Save(ctx); err != nil {
        return err // triggers rollback
    }

    // Lock a row
    if err := db.Lock(ctx, conn, u, u.Id); err != nil {
        return err
    }

    u.Status = 1
    return u.Save(ctx) // triggers commit
})

if err != nil {
    txErr := err.(ksql.TxError)
    fmt.Println("begin:", txErr.Begin())
    fmt.Println("call:", txErr.Call())
    fmt.Println("commit:", txErr.Commit())
    fmt.Println("rollback:", txErr.Rollback())
}
```

With custom transaction options:

```go
db.TransactionBy(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable}, func(ctx context.Context, conn ksql.ConnectionInterface) error {
    // ...
})
```

### Savepoints

Savepoints are automatically used for nested transactions. Drivers that support savepoints (e.g. MySQL) create savepoints on each nested `Begin` call and release/rollback them when committed/rolled back.

## Sharding

db-go supports hash-based database sharding out of the box.

### Initialize sharded connections

```go
import "github.com/kovey/db-go/v3/sharding"

sharding.Init([]db.Config{
    {DriverName: "mysql", DataSourceName: "root:pass@tcp(127.0.0.1:3306)/db0?...", MaxIdleConns: 10, MaxOpenConns: 50},
    {DriverName: "mysql", DataSourceName: "root:pass@tcp(127.0.0.1:3307)/db1?...", MaxIdleConns: 10, MaxOpenConns: 50},
})
```

### Define a sharded model

```go
type Order struct {
    *sharding.Model
    Id     int64
    UserId int64
    Amount float64
}

func NewOrder() *Order {
    return &Order{Model: sharding.NewModel("order", "id", model.Type_Int)}
}

func (o *Order) Clone() ksql.RowInterface { return NewOrder() }
func (o *Order) Columns() []string        { return []string{"id", "user_id", "amount"} }
func (o *Order) Values() []any            { return []any{&o.Id, &o.UserId, &o.Amount} }
func (o *Order) Save(ctx context.Context) error { return o.Model.SaveBy(ctx, o) }
func (o *Order) Delete(ctx context.Context) error { return o.Model.DeleteBy(ctx, o) }
```

### Query with sharding key

```go
// Routes to the correct shard based on userId
o := NewOrder()
sharding.Row(userId, o).Where("id", ksql.Eq, 1).First(ctx)

var orders []*Order
sharding.Rows(userId, &orders).Where("amount", ksql.Gt, 100).All(ctx)
```

### Sharded table auto-creation

```go
mgr := model.NewTableManager([]*model.Template{
    {
        Table:         "order",
        Keep:          30,                   // keep 30 days
        Type:          ksql.Sharding_Day,    // table_order_20260101, ...
        TemplateTable: "order_template",
    },
})
mgr.Create(ctx)  // creates tables for ±3 days
mgr.Delete(ctx)  // drops tables older than Keep days
```

## Logging & Tracing

### SQL logging

```go
db.Init(db.Config{
    // ...
    LogOpened: true,
    LogMax:    2048, // channel buffer size
})

// Optionally write logs to file
db.LogUseFile("/var/log/myapp")
```

Log output format (JSON per line):

```json
{
    "delay": "12.345ms",
    "trace_id": "t_1234567890",
    "begin_time": "2026-01-15 10:30:00.000",
    "end_time": "2026-01-15 10:30:00.012",
    "sql": "SELECT `id`, `name` FROM `user` WHERE `id` = ?",
    "span_id": "ABCDEFG-HIJKLMN"
}
```

### Tracing

```go
ctx := db.NewContext(context.Background()).WithTraceId(requestId)
// All SQL executed with this context will carry the trace_id
db.Model(u).Where("id", ksql.Eq, 1).First(ctx)
```

## Architecture

```
ksql        — Interfaces (QueryInterface, WhereInterface, TableInterface, etc.)
├── db      — Connection management, global helpers, Builder (model-based queries)
├── sql     — SQL generation (SELECT, INSERT, UPDATE, DELETE, DDL, etc.)
├── model   — Base Model with CRUD + lifecycle hooks
├── sharding — Hash-based multi-database sharding support
├── logger  — SQL log capture (stdout or file)
├── express — Raw SQL statement wrapper with type detection
├── schema  — Database schema helpers
├── ksql    — CLI entry point
└── korm    — Code generation entry point
```

## Comparison with traditional ORMs

| | db-go | GORM | sqlx |
|---|---|---|---|
| **Reflection** | None | Heavy | Moderate |
| **Performance** | Fast (no reflection overhead) | Slower | Fast |
| **Type safety** | Compile-time (explicit mapping) | Runtime | Runtime |
| **SQL control** | Full | Abstracted | Full |
| **Learning curve** | Moderate (explicit mapping) | Low | Low |
| **Code generation** | Via korm | Via gen | None built-in |

## License

MIT — see [LICENSE](LICENSE).
