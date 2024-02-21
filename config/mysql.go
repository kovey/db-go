package config

import (
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Mysql struct {
	Host          string `yaml:"host" json:"host"`
	Port          int    `yaml:"port" json:"port"`
	Username      string `yaml:"username" json:"username"`
	Password      string `yaml:"password" json:"password"`
	Dbname        string `yaml:"dbname" json:"dbname"`
	Charset       string `yaml:"charset" json:"charset"`
	ActiveMax     int    `yaml:"active_max" json:"active_max"`
	ConnectionMax int    `yaml:"connection_max" json:"connection_max"`
	LifeTime      int    `yaml:"life_time" json:"life_time"`
	Dev           string `yaml:"dev" json:"dev"`
}

func (m *Mysql) ToDSN() string {
	mc := mysql.NewConfig()
	mc.User = m.Username
	mc.Passwd = m.Password
	mc.Net = "tcp"
	mc.Addr = fmt.Sprintf("%s:%d", m.Host, m.Port)
	mc.DBName = m.Dbname
	mc.Params = map[string]string{
		"charset": m.Charset,
	}
	mc.Loc = time.Local

	return mc.FormatDSN()
}
