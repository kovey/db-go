package main

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kovey/db-go/v3/db"
	"github.com/kovey/db-go/v3/examples/models"
	"github.com/kovey/db-go/v3/model"
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

	ctx := context.Background()
	u := models.NewUser()
	if err := model.Query(u).Where("id", "=", 2).First(ctx, u); err != nil {
		panic(err)
	}

	fmt.Println("user: ", u)
	var users []*models.User
	if err := model.Query(models.NewUser()).Where("id", ">", 0).All(ctx, &users); err != nil {
		panic(err)
	}

	for _, u := range users {
		fmt.Printf("uu: %+v\n", u)
	}

	db.Close()
}
