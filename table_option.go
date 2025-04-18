package ksql

type RowFormat string

const (
	Row_Format_Default    RowFormat = "DEFAULT"
	Row_Format_Dynamic    RowFormat = "DYNAMIC"
	Row_Format_Fixed      RowFormat = "FIXED"
	Row_Format_Compressed RowFormat = "COMPRESSED"
	Row_Format_Redundant  RowFormat = "REDUNDANT"
	Row_Format_Compact    RowFormat = "COMPACT"
)

type InsertMethod string

const (
	Insert_Method_No    InsertMethod = "NO"
	Insert_Method_Last  InsertMethod = "LAST"
	Insert_Method_First InsertMethod = "FIRST"
)

type TableOptVal string

const (
	Table_Opt_Val_0       TableOptVal = "0"
	Table_Opt_Val_1       TableOptVal = "1"
	Table_Opt_Val_Default TableOptVal = "DEFAULT"
	Table_Opt_Val_Y       TableOptVal = "Y"
	Table_Opt_Val_N       TableOptVal = "N"
	Table_Opt_Val_ZLIB    TableOptVal = "ZLIB"
	Table_Opt_Val_LZ4     TableOptVal = "LZ4"
	Table_Opt_Val_NONE    TableOptVal = "NONE"
)

type TableOptKey string

const (
	Table_Opt_Key_Autoextend_Size            TableOptKey = "AUTOEXTEND_SIZE"
	Table_Opt_Key_Auto_Increment             TableOptKey = "AUTO_INCREMENT"
	Table_Opt_Key_Avg_Row_Length             TableOptKey = "AVG_ROW_LENGTH"
	Table_Opt_Key_Character_Set              TableOptKey = "CHARACTER SET"
	Table_Opt_Key_Checksum                   TableOptKey = "CHECKSUM"
	Table_Opt_Key_Collate                    TableOptKey = "COLLATE"
	Table_Opt_Key_Comment                    TableOptKey = "COMMENT"
	Table_Opt_Key_Compression                TableOptKey = "COMPRESSION"
	Table_Opt_Key_Connection                 TableOptKey = "CONNECTION"
	Table_Opt_Key_Directory                  TableOptKey = "DIRECTORY"
	Table_Opt_Key_Delay_Key_Write            TableOptKey = "DELAY_KEY_WRITE"
	Table_Opt_Key_Encryption                 TableOptKey = "ENCRYPTION"
	Table_Opt_Key_Engine                     TableOptKey = "ENGINE"
	Table_Opt_Key_Engine_Attribute           TableOptKey = "ENGINE_ATTRIBUTE"
	Table_Opt_Key_Insert_Method              TableOptKey = "INSERT_METHOD"
	Table_Opt_Key_Key_Block_Size             TableOptKey = "KEY_BLOCK_SIZE"
	Table_Opt_Key_Max_Rows                   TableOptKey = "MAX_ROWS"
	Table_Opt_Key_Min_Rows                   TableOptKey = "MIN_ROWS"
	Table_Opt_Key_Pack_Keys                  TableOptKey = "PACK_KEYS"
	Table_Opt_Key_Password                   TableOptKey = "PASSWORD"
	Table_Opt_Key_Secondary_Engine_Attribute TableOptKey = "SECONDARY_ENGINE_ATTRIBUTE"
	Table_Opt_Key_Stats_Auto_Recalc          TableOptKey = "STATS_AUTO_RECALC"
	Table_Opt_Key_Stats_Persistent           TableOptKey = "STATS_PERSISTENT"
	Table_Opt_Key_Stats_Sample_Pages         TableOptKey = "STATS_SAMPLE_PAGES"
)

func (t TableOptKey) IsStr() bool {
	switch t {
	case Table_Opt_Key_Comment, Table_Opt_Key_Compression, Table_Opt_Key_Connection, Table_Opt_Key_Directory, Table_Opt_Key_Encryption,
		Table_Opt_Key_Engine_Attribute, Table_Opt_Key_Password, Table_Opt_Key_Secondary_Engine_Attribute:
		return true
	default:
		return false
	}
}

type TableOptPrefix string

const (
	Table_Opt_Prefix_Default TableOptPrefix = "DEFAULT"
	Table_Opt_Prefix_Data    TableOptPrefix = "DATA"
	Table_Opt_Prefix_Index   TableOptPrefix = "INDEX"
)

type PartitionDefinitionOptKey string

const (
	Partition_Definition_Opt_Key_Engine          = "STORAGE ENGINE"
	Partition_Definition_Opt_Key_Comment         = "COMMENT"
	Partition_Definition_Opt_Key_Data_Directory  = "DATA DIRECTORY"
	Partition_Definition_Opt_Key_Index_Directory = "INDEX DIRECTORY"
	Partition_Definition_Opt_Key_Max_Rows        = "MAX_ROWS"
	Partition_Definition_Opt_Key_Min_Rows        = "MIN_ROWS"
	Partition_Definition_Opt_Key_Tablespace      = "TABLESPACE"
)

func (t PartitionDefinitionOptKey) IsStr() bool {
	switch t {
	case Partition_Definition_Opt_Key_Comment, Partition_Definition_Opt_Key_Data_Directory, Partition_Definition_Opt_Key_Index_Directory:
		return true
	default:
		return false
	}
}
