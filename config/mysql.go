package config

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
