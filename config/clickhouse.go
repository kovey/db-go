package config

type ClickHouse struct {
	Username     string  `yaml:"username"`
	Password     string  `yaml:"password"`
	Dbname       string  `yaml:"dbname"`
	Debug        bool    `yaml:"debug"`
	OpenStrategy string  `yaml:"open_strategy"`
	BlockSize    int     `yaml:"block_size"`
	PoolSize     int     `yaml:"pool_size"`
	Compress     int     `yaml:"compress"`
	Timeout      Timeout `yaml:"timeout"`
	Cluster      Cluster `yaml:"cluster"`
	Server       Addr    `yaml:"server"`
}

type Timeout struct {
	Read  int `yaml:"read"`
	Write int `yaml:"write"`
}

type Cluster struct {
	Open    string `yaml:"open"`
	Servers []Addr `yaml:"servers"`
}

type Addr struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
