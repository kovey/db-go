package config

import "fmt"

type ClickHouse struct {
	Username      string  `yaml:"username" json:"username"`
	Password      string  `yaml:"password" json:"password"`
	Dbname        string  `yaml:"dbname" json:"dbname"`
	Debug         bool    `yaml:"debug" json:"debug"`
	BlockSize     int     `yaml:"block_size" json:"block_size"`
	Compress      int     `yaml:"compress" json:"compress"`
	Timeout       Timeout `yaml:"timeout" json:"timeout"`
	Cluster       Cluster `yaml:"cluster" json:"cluster"`
	Server        Addr    `yaml:"server" json:"server"`
	ActiveMax     int     `yaml:"active_max" json:"active_max"`
	ConnectionMax int     `yaml:"connection_max" json:"connection_max"`
	LifeTime      int     `yaml:"life_time" json:"life_time"`
}

type Timeout struct {
	Read int `yaml:"read" json:"read"`
	Exec int `yaml:"exec" json:"exec"`
	Dial int `yaml:"dial" json:"dial"`
}

type Cluster struct {
	Open    string `yaml:"open" json:"open"`
	Servers []Addr `yaml:"servers" json:"servers"`
}

type Addr struct {
	Host string `yaml:"host" json:"host"`
	Port int    `yaml:"port" json:"port"`
}

func (a Addr) Info() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}
