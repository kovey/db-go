package meta

import "fmt"

type Func string

const (
	Func_None              Func = ""
	Func_ASCII             Func = "ASCII"
	Func_CHAR_LENGTH       Func = "CHAR_LENGTH"
	Func_CHARACTER_LENGTH  Func = "CHARACTER_LENGTH"
	Func_CONCAT            Func = "CONCAT"
	Func_CONCAT_WS         Func = "CONCAT_WS"
	Func_FIELD             Func = "FIELD"
	Func_FIND_IN_SET       Func = "FIND_IN_SET"
	Func_FORMAT            Func = "FORMAT"
	Func_INSERT            Func = "INSERT"
	Func_LOCATE            Func = "LOCATE"
	Func_LCASE             Func = "LCASE"
	Func_LEFT              Func = "LEFT"
	Func_LOWER             Func = "LOWER"
	Func_LPAD              Func = "LPAD"
	Func_LTRIM             Func = "LTRIM"
	Func_MID               Func = "MID"
	Func_POSITION          Func = "POSITION"
	Func_REPEAT            Func = "REPEAT"
	Func_REPLACE           Func = "REPLACE"
	Func_REVERSE           Func = "REVERSE"
	Func_RIGHT             Func = "RIGHT"
	Func_RPAD              Func = "RPAD"
	Func_RTRIM             Func = "RTRIM"
	Func_SPACE             Func = "SPACE"
	Func_STRCMP            Func = "STRCMP"
	Func_SUBSTR            Func = "SUBSTR"
	Func_SUBSTRING         Func = "SUBSTRING"
	Func_SUBSTRING_INDEX   Func = "SUBSTRING_INDEX"
	Func_TRIM              Func = "TRIM"
	Func_UCASE             Func = "UCASE"
	Func_UPPER             Func = "UPPER"
	Func_ABS               Func = "ABS"
	Func_ACOS              Func = "ACOS"
	Func_ASIN              Func = "ASIN"
	Func_ATAN              Func = "ATAN"
	Func_ATAN2             Func = "ATAN2"
	Func_AVG               Func = "AVG"
	Func_CEIL              Func = "CEIL"
	Func_CEILING           Func = "CEILING"
	Func_COS               Func = "COS"
	Func_COT               Func = "COT"
	Func_COUNT             Func = "COUNT"
	Func_DEGREES           Func = "DEGREES"
	Func_DIV               Func = "DIV"
	Func_EXP               Func = "EXP"
	Func_FLOOR             Func = "FLOOR"
	Func_GREATEST          Func = "GREATEST"
	Func_LEAST             Func = "LEAST"
	Func_LN                Func = "LN"
	Func_LOG               Func = "LOG"
	Func_LOG10             Func = "LOG10"
	Func_LOG2              Func = "LOG2"
	Func_MAX               Func = "MAX"
	Func_MIN               Func = "MIN"
	Func_MOD               Func = "MOD"
	Func_PI                Func = "PI"
	Func_POW               Func = "POW"
	Func_POWER             Func = "POWER"
	Func_RADIANS           Func = "RADIANS"
	Func_RAND              Func = "RAND"
	Func_ROUND             Func = "ROUND"
	Func_SIGN              Func = "SIGN"
	Func_SIN               Func = "SIN"
	Func_SQRT              Func = "SQRT"
	Func_SUM               Func = "SUM"
	Func_TAN               Func = "TAN"
	Func_TRUNCATE          Func = "TRUNCATE"
	Func_ADDDATE           Func = "ADDDATE"
	Func_ADDTIME           Func = "ADDTIME"
	Func_CURDATE           Func = "CURDATE"
	Func_CURRENT_DATE      Func = "CURRENT_DATE"
	Func_CURRENT_TIME      Func = "CURRENT_TIME"
	Func_CURRENT_TIMESTAMP Func = "CURRENT_TIMESTAMP"
	Func_CURTIME           Func = "CURTIME"
	Func_DATE              Func = "DATE"
	Func_DATEDIFF          Func = "DATEDIFF"
	Func_DATE_ADD          Func = "DATE_ADD"
	Func_DATE_FORMAT       Func = "DATE_FORMAT"
	Func_DATE_SUB          Func = "DATE_SUB"
	Func_DAY               Func = "DAY"
	Func_DAYNAME           Func = "DAYNAME"
	Func_DAYOFMONTH        Func = "DAYOFMONTH"
	Func_DAYOFWEEK         Func = "DAYOFWEEK"
	Func_DAYOFYEAR         Func = "DAYOFYEAR"
	Func_EXTRACT           Func = "EXTRACT"
	Func_FROM_DAYS         Func = "FROM_DAYS"
	Func_HOUR              Func = "HOUR"
	Func_LAST_DAY          Func = "LAST_DAY"
	Func_LOCALTIME         Func = "LOCALTIME"
	Func_LOCALTIMESTAMP    Func = "LOCALTIMESTAMP"
	Func_MAKEDATE          Func = "MAKEDATE"
	Func_MAKETIME          Func = "MAKETIME"
	Func_MICROSECOND       Func = "MICROSECOND"
	Func_MINUTE            Func = "MINUTE"
	Func_MONTHNAME         Func = "MONTHNAME"
	Func_MONTH             Func = "MONTH"
	Func_NOW               Func = "NOW"
	Func_PERIOD_ADD        Func = "PERIOD_ADD"
	Func_PERIOD_DIFF       Func = "PERIOD_DIFF"
	Func_QUARTER           Func = "QUARTER"
	Func_SECOND            Func = "SECOND"
	Func_SEC_TO_TIME       Func = "SEC_TO_TIME"
	Func_STR_TO_DATE       Func = "STR_TO_DATE"
	Func_SUBDATE           Func = "SUBDATE"
	Func_SUBTIME           Func = "SUBTIME"
	Func_SYSDATE           Func = "SYSDATE"
	Func_TIME              Func = "TIME"
	Func_TIME_FORMAT       Func = "TIME_FORMAT"
	Func_TIME_TO_SEC       Func = "TIME_TO_SEC"
	Func_TIMEDIFF          Func = "TIMEDIFF"
	Func_TIMESTAMP         Func = "TIMESTAMP"
	Func_TO_DAYS           Func = "TO_DAYS"
	Func_WEEK              Func = "WEEK"
	Func_WEEKDAY           Func = "WEEKDAY"
	Func_WEEKOFYEAR        Func = "WEEKOFYEAR"
	Func_YEAR              Func = "YEAR"
	Func_YEARWEEK          Func = "YEARWEEK"
	Func_BIN               Func = "BIN"
	Func_BINARY            Func = "BINARY"
	Func_CAST              Func = "CAST"
	Func_COALESCE          Func = "COALESCE"
	Func_CONNECTION_ID     Func = "CONNECTION_ID"
	Func_CONV              Func = "CONV"
	Func_CONVERT           Func = "CONVERT"
	Func_CURRENT_USER      Func = "CURRENT_USER"
	Func_DATABASE          Func = "DATABASE"
	Func_IF                Func = "IF"
	Func_IFNULL            Func = "IFNULL"
	Func_ISNULL            Func = "ISNULL"
	Func_LAST_INSERT_ID    Func = "LAST_INSERT_ID"
	Func_NULLIF            Func = "NULLIF"
	Func_SESSION_USER      Func = "SESSION_USER"
	Func_SYSTEM_USER       Func = "SYSTEM_USER"
	Func_USER              Func = "USER"
	Func_VERSION           Func = "VERSION"
)

const (
	columnFormat = "%s AS `%s`"
	columnFunc   = "%s(%s) AS `%s`"
	divFunc      = "%s %s %s AS `%s`"
)

type Field struct {
	Name                string
	Table               string
	IsConstOrExpression bool
}

func NewField(name, table string, isConstOrExpression bool) *Field {
	return &Field{Name: name, Table: table, IsConstOrExpression: isConstOrExpression}
}

func (f *Field) String() string {
	if f == nil {
		return ""
	}

	if f.IsConstOrExpression {
		return f.Name
	}

	if f.Table == "" {
		return fmt.Sprintf("`%s`", f.Name)
	}

	return fmt.Sprintf("`%s`.`%s`", f.Table, f.Name)
}

type Column struct {
	Name  *Field
	Ext   *Field
	Alias string
	Func  Func
}

func NewColumn(name, alias string) *Column {
	return NewFuncComlumn(NewField(name, "", false), alias, Func_None, nil)
}

func NewFuncComlumn(name *Field, alias string, f Func, ext *Field) *Column {
	return &Column{Name: name, Alias: alias, Func: f, Ext: ext}
}

func (c *Column) SetTable(table string) {
	if c.Name != nil {
		c.Name.Table = table
	}

	if c.Ext != nil {
		c.Ext.Table = table
	}
}

func (c *Column) String() string {
	if c.Func == Func_None {
		return fmt.Sprintf(columnFormat, c.Name, c.Alias)
	}

	switch c.Func {
	case Func_DIV:
		return fmt.Sprintf(divFunc, c.Name, c.Func, c.Ext, c.Alias)
	case Func_CAST:
		return fmt.Sprintf(columnFunc, c.Func, fmt.Sprintf("%s AS %s", c.Name, c.Ext), c.Alias)
	case Func_CONVERT:
		return fmt.Sprintf(columnFunc, c.Func, fmt.Sprintf("%s USING %s", c.Name, c.Ext), c.Alias)
	default:
		return fmt.Sprintf(columnFunc, c.Func, c.Name, c.Alias)
	}
}
