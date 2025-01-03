package models

// Code generated by ksql.
// Do'nt Edit!!!
// Do'nt Edit!!!
// Do'nt Edit!!!
// 用户表
// from database: test_dev
// table:         user
// orm version:   1.0.1
// created time:  2025-01-03 11:11:29
// ddl:
/**
CREATE TABLE `user` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '用户id',
  `username` varchar(45) NOT NULL DEFAULT '' COMMENT '用户名',
  `password` varchar(60) NOT NULL DEFAULT '' COMMENT '密码',
  `age` int NOT NULL DEFAULT '0' COMMENT '年龄',
  `created_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `email` varchar(45) NOT NULL DEFAULT '' COMMENT '邮箱',
  `utm_source` varchar(45) NOT NULL DEFAULT 'ios' COMMENT '来源',
  `reg_ip` varchar(45) NOT NULL DEFAULT '127.0.0.1' COMMENT '注册ip',
  PRIMARY KEY (`id`),
  KEY `idx_email` (`email`),
  KEY `idx_ip_utm` (`utm_source`,`reg_ip`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表'
*/

import (
	"context"

	"github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/model"
	"time"
)

const (
	Table_User             = "user"         // 用户表
	Table_User_Age         = "age"          // 年龄
	Table_User_CreatedTime = "created_time" // 创建时间
	Table_User_Email       = "email"        // 邮箱
	Table_User_Id          = "id"           // 用户id
	Table_User_Password    = "password"     // 密码
	Table_User_RegIp       = "reg_ip"       // 注册ip
	Table_User_UpdatedTime = "updated_time" // 更新时间
	Table_User_Username    = "username"     // 用户名
	Table_User_UtmSource   = "utm_source"   // 来源
)

type User struct {
	*model.Model `db:"-" json:"-"` // model
	Age          int               `db:"age" json:"age"`                   // 年龄
	CreatedTime  time.Time         `db:"created_time" json:"created_time"` // 创建时间
	Email        string            `db:"email" json:"email"`               // 邮箱
	Id           int               `db:"id" json:"id"`                     // 用户id
	Password     string            `db:"password" json:"password"`         // 密码
	RegIp        string            `db:"reg_ip" json:"reg_ip"`             // 注册ip
	UpdatedTime  time.Time         `db:"updated_time" json:"updated_time"` // 更新时间
	Username     string            `db:"username" json:"username"`         // 用户名
	UtmSource    string            `db:"utm_source" json:"utm_source"`     // 来源
}

func NewUser() *User {
	return &User{Model: model.NewModel(Table_User, "id", model.Type_Int)}
}

func (self *User) Save(ctx context.Context) error {
	return self.Model.Save(ctx, self)
}

func (self *User) Clone() ksql.RowInterface {
	return NewUser()
}

func (self *User) Values() []any {
	return []any{&self.Age, &self.CreatedTime, &self.Email, &self.Id, &self.Password, &self.RegIp, &self.UpdatedTime, &self.Username, &self.UtmSource}
}

func (self *User) Columns() []string {
	return []string{"age", "created_time", "email", "id", "password", "reg_ip", "updated_time", "username", "utm_source"}
}

func (self *User) Delete(ctx context.Context) error {
	return self.Model.Delete(ctx, self)
}

func (self *User) Query() ksql.BuilderInterface[*User] {
	return model.Query(self)
}
