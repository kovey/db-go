package ksql

const (
	CURRENT_TIMESTAMP                             = "CURRENT_TIMESTAMP"
	CURRENT_TIMESTAMP_ON_UPDATE_CURRENT_TIMESTAMP = "CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"
	NULL                                          = "NULL"
)

func IsDefaultKeyword(val string) bool {
	return val == CURRENT_TIMESTAMP || val == CURRENT_TIMESTAMP_ON_UPDATE_CURRENT_TIMESTAMP || val == NULL
}
