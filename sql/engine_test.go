package sql

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type e_data struct {
	meta int
	name string
}

func (e *e_data) String() string {
	return fmt.Sprintf("%d-%s", e.meta, e.name)
}

func TestEngineFormat(t *testing.T) {
	e := DefaultEngine()
	dTime, _ := time.Parse(time.DateTime, "2025-04-02 13:40:42")
	in := NewInsert().Table("user")
	in.Add("num_int8", int8(8))
	in.Add("num_int16", int16(16))
	in.Add("num_int32", int32(32))
	in.Add("num_int64", int64(64))
	in.Add("num_uint8", uint8(8))
	in.Add("num_uint16", uint16(16))
	in.Add("num_uint32", uint32(32))
	in.Add("num_uint64", uint64(64))
	in.Add("num_uint", uint(642))
	in.Add("num_int", int(642))
	in.Add("bool_val", true).Add("str", "kovey").Add("date", dTime)
	in.Add("num_float32", float32(11.11)).Add("num_float64", float64(12.12))
	var (
		num_int     = int(1)
		num_int8    = int8(8)
		num_int16   = int16(16)
		num_int32   = int32(32)
		num_int64   = int64(64)
		num_uint    = uint(1)
		num_uint8   = uint8(8)
		num_uint16  = uint16(16)
		num_uint32  = uint32(32)
		num_uint64  = uint64(64)
		bool_val    = false
		str_val     = "kkk"
		num_float32 = float32(13.13)
		num_float64 = float64(14.13)
	)

	in.Add("p_num_int", &num_int)
	in.Add("p_num_int8", &num_int8)
	in.Add("p_num_int16", &num_int16)
	in.Add("p_num_int32", &num_int32)
	in.Add("p_num_int64", &num_int64)
	in.Add("p_num_uint", &num_uint)
	in.Add("p_num_uint8", &num_uint8)
	in.Add("p_num_uint16", &num_uint16)
	in.Add("p_num_uint32", &num_uint32)
	in.Add("p_num_uint64", &num_uint64)
	in.Add("p_bool_val", &bool_val)
	in.Add("p_str_val", &str_val)
	in.Add("p_date", &dTime)
	in.Add("p_num_float32", &num_float32)
	in.Add("p_num_float64", &num_float64)
	in.Add("e_data", &e_data{meta: 100, name: "aaaa"})
	assert.Equal(t, "INSERT INTO `user` (`num_int8`, `num_int16`, `num_int32`, `num_int64`, `num_uint8`, `num_uint16`, `num_uint32`, `num_uint64`, `num_uint`, `num_int`, `bool_val`, `str`, `date`, `num_float32`, `num_float64`, `p_num_int`, `p_num_int8`, `p_num_int16`, `p_num_int32`, `p_num_int64`, `p_num_uint`, `p_num_uint8`, `p_num_uint16`, `p_num_uint32`, `p_num_uint64`, `p_bool_val`, `p_str_val`, `p_date`, `p_num_float32`, `p_num_float64`, `e_data`) VALUES (8, 16, 32, 64, 8, 16, 32, 64, 642, 642, true, 'kovey', '2025-04-02 13:40:42', 11.110000, 12.120000, 1, 8, 16, 32, 64, 1, 8, 16, 32, 64, false, 'kkk', '2025-04-02 13:40:42', 13.130000, 14.130000, '100-aaaa')", e.Format(in))
	assert.Equal(t, "INSERT INTO `user` (`num_int8`, `num_int16`, `num_int32`, `num_int64`, `num_uint8`, `num_uint16`, `num_uint32`, `num_uint64`, `num_uint`, `num_int`, `bool_val`, `str`, `date`, `num_float32`, `num_float64`, `p_num_int`, `p_num_int8`, `p_num_int16`, `p_num_int32`, `p_num_int64`, `p_num_uint`, `p_num_uint8`, `p_num_uint16`, `p_num_uint32`, `p_num_uint64`, `p_bool_val`, `p_str_val`, `p_date`, `p_num_float32`, `p_num_float64`, `e_data`) VALUES (8, 16, 32, 64, 8, 16, 32, 64, 642, 642, true, kovey, 2025-04-02 13:40:42, 11.110000, 12.120000, 1, 8, 16, 32, 64, 1, 8, 16, 32, 64, false, kkk, 2025-04-02 13:40:42, 13.130000, 14.130000, 100-aaaa)", e.formatOriginal(in))
}

func TestEngineFormatRaw(t *testing.T) {
	e := DefaultEngine()
	raw := Raw("select * from user where id = ? and name like ? between ? and ? limit ?", 1, "%test%", 100, 1000, 10)
	assert.Equal(t, "select * from user where id = 1 and name like '%test%' between 100 and 1000 limit 10", e.FormatRaw(raw))
	assert.Equal(t, "select * from user where id = 1 and name like %test% between 100 and 1000 limit 10", e.formatOriginalRaw(raw))
}
