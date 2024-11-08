package serv

const (
	config_tpl = `
# ksql env config
# tool name
APP_NAME = ksql

# database config
DB_DRIVER   = mysql
DB_HOST     = 127.0.0.1
DB_PORT     = 3306
DB_USER     = root
DB_PASSWORD = password
DB_NAME     = test
DB_CHARSET  = utf8mb4

# migrator plugin path
PLUGIN_MIGRATOR_PATH = path/to/migrator

# diff sql path
DIFF_SQL_PATH = path/to/diff

# models path
MODELS_PATH   = path/to/models

# to database config
TO_DB_HOST     = 127.0.0.1
TO_DB_PORT     = 3306
TO_DB_USER     = root
TO_DB_PASSWORD = password
TO_DB_NAME     = test_prod
TO_DB_CHARSET  = utf8mb4
`
)
