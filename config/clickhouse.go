package config

import "fmt"

type ClickHouse struct {
	Username      string  `yaml:"username"`
	Password      string  `yaml:"password"`
	Dbname        string  `yaml:"dbname"`
	Debug         bool    `yaml:"debug"`
	BlockSize     int     `yaml:"block_size"`
	Compress      int     `yaml:"compress"`
	Timeout       Timeout `yaml:"timeout"`
	Cluster       Cluster `yaml:"cluster"`
	Server        Addr    `yaml:"server"`
	ActiveMax     int     `yaml:"active_max"`
	ConnectionMax int     `yaml:"connection_max"`
	LifeTime      int     `yaml:"life_time"`
}

type Timeout struct {
	Read int `yaml:"read"`
	Exec int `yaml:"exec"`
	Dial int `yaml:"dial"`
}

type Cluster struct {
	Open    string `yaml:"open"`
	Servers []Addr `yaml:"servers"`
}

type Addr struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (a Addr) Info() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}
