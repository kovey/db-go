package models

import (
	"time"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/db"
)

//go:korm values,columns
type UserExtJoin struct {
	*db.Row
	Id            int        `db:"u.id" json:"id"`                           // 用户ID
	Username      string     `db:"u.username" json:"username"`               // 用户名
	Password      string     `db:"u.password" json:"password"`               // 密码
	Age           int        `db:"u.age" json:"age"`                         // 年龄
	CreatedTime   time.Time  `db:"u.created_time" json:"created_time"`       // 创建时间
	UpdatedTime   time.Time  `db:"u.updated_time" json:"updated_time"`       // 更新时间
	Email         string     `db:"u.email" json:"email"`                     // 邮箱
	UtmSource     string     `db:"u.utm_source" json:"utm_source"`           // 来源
	RegIp         string     `db:"u.reg_ip" json:"reg_ip"`                   // 注册ip
	PrevLoginDate time.Time  `db:"e.prev_login_date" json:"prev_login_date"` // 上次登录日期
	PrevLoginTime *time.Time `db:"e.prev_login_time" json:"prev_login_time"` // 上次登录时间
	PrevLoginIp   string     `db:"e.prev_login_ip" json:"prev_login_ip"`     // 上次登录IP
}

func NewUserExtJoin() *UserExtJoin {
	return &UserExtJoin{Row: &db.Row{}}
}

func (u *UserExtJoin) Clone() ksql.RowInterface {
	return NewUserExtJoin()
}
