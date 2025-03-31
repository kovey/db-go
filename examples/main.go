package main

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
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
}
