package config

import (
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Mysql struct {
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	Dbname        string `yaml:"dbname"`
	Charset       string `yaml:"charset"`
	ActiveMax     int    `yaml:"active_max"`
	ConnectionMax int    `yaml:"connection_max"`
	LifeTime      int    `yaml:"life_time"`
	Dev           string `yaml:"dev"`
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
