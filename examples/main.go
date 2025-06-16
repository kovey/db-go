package main

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/db"
	"github.com/kovey/db-go/v3/examples/models"
)

func main() {
	conf := db.Config{
		DriverName:     "mysql",
		DataSourceName: "root:some34QA@123@tcp(127.0.0.1:3306)/test_dev?charset=utf8mb4&parseTime=true",
		MaxIdleTime:    time.Second * 60,
		MaxLifeTime:    time.Second * 120,
		MaxIdleConns:   10,
		MaxOpenConns:   50,
		LogOpened:      true,
		LogMax:         1024,
	}

	if err := db.Init(conf); err != nil {
		panic(err)
	}

	defer db.Close()

	ctx := db.NewContext(context.Background()).WithTraceId(fmt.Sprintf("t_%d", time.Now().UnixNano()))
	u := models.NewUser()
	if err := db.Model(u).Where("id", "=", 2).First(ctx); err != nil {
		panic(err)
	}

	fmt.Println("user: ", u)
	var users []*models.User
	if err := db.Models(&users).Where("id", ">", 0).All(ctx); err != nil {
		panic(err)
	}

	for _, u := range users {
		fmt.Printf("uu: %+v\n", u)
	}

	pageInfo, err := db.Models(&[]*models.User{}).Where("id", ">", 0).Pagination(ctx, 1, 10)
	if err != nil {
		panic(err)
	}

	fmt.Println("total page:", pageInfo.TotalPage(), "count:", pageInfo.TotalCount())
	for _, u := range pageInfo.List() {
		fmt.Printf("uu: %+v\n", u)
	}

	var uus []*models.User
	if err := db.Rows(&uus).Where("id", ">", 0).Columns(u.Columns()...).Table(u.Table()).All(ctx); err != nil {
		panic(err)
	}
	for _, u := range uus {
		fmt.Printf("uu: %+v\n", u)
	}

	var uut []*models.UserTest
	if err := db.Rows(&uut).Where("id", ">", 0).Columns(u.Columns()...).Table(u.Table()).All(ctx); err != nil {
		panic(err)
	}
	for _, u := range uus {
		fmt.Printf("uut: %+v\n", u)
	}

	leftJoin(ctx)
}

func leftJoin(ctx context.Context) {
	var rows []*models.UserExtJoin
	builder := db.Rows(&rows).Where("u.id", ksql.Gt, 1).Columns(models.NewUserExtJoin().Columns()...).Table("user").As("u")
	builder.LeftJoin("user_ext").As("e").On("u.id", "=", "e.id")
	if err := builder.All(ctx); err != nil {
		panic(err)
	}

	for _, u := range rows {
		fmt.Printf("uue: %+v\n", u)
	}
}
