package db

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func _testGetData() *Data {
	dt := NewData()
	tmpStr1 := "koveychange2"
	str1 := &tmpStr1
	dt._set("key1", "koveychange")
	dt._set("key2", str1)
	dt._set("key3", &str1)
	tmpInt1 := 2
	int1 := &tmpInt1
	dt._set("key4", int1)
	dt._set("key5", &int1)
	tmpInt2 := int8(16)
	int2 := &tmpInt2
	dt._set("key6", int2)
	dt._set("key7", &int2)
	tmpInt3 := int16(32)
	int3 := &tmpInt3
	dt._set("key8", int3)
	dt._set("key9", &int3)
	tmpInt4 := int32(64)
	int4 := &tmpInt4
	dt._set("key10", int4)
	dt._set("key11", &int4)
	tmpInt5 := int64(128)
	int5 := &tmpInt5
	dt._set("key12", int5)
	dt._set("key13", &int5)
	tmpuInt1 := uint(2)
	uint1 := &tmpuInt1
	dt._set("key14", uint1)
	dt._set("key15", &uint1)
	tmpuInt2 := uint8(16)
	uint2 := &tmpuInt2
	dt._set("key16", uint2)
	dt._set("key17", &uint2)
	tmpuInt3 := uint16(32)
	uint3 := &tmpuInt3
	dt._set("key18", uint3)
	dt._set("key19", &uint3)
	tmpuInt4 := uint32(64)
	uint4 := &tmpuInt4
	dt._set("key20", uint4)
	dt._set("key21", &uint4)
	tmpuInt5 := uint64(128)
	uint5 := &tmpuInt5
	dt._set("key22", uint5)
	dt._set("key23", &uint5)
	tmpBool := false
	bl := &tmpBool
	dt._set("key24", bl)
	dt._set("key25", &bl)
	tmpFloat := float32(64.64)
	float := &tmpFloat
	dt._set("key26", float)
	dt._set("key27", &float)
	tmpFloat1 := float64(128.128)
	float1 := &tmpFloat1
	dt._set("key28", float1)
	dt._set("key29", &float1)
	tmpNow, _ := time.Parse(time.DateTime, "2025-04-11 11:11:11")
	now := &tmpNow
	dt._set("key30", now)
	dt._set("key31", &now)
	dt._set("key32", &sql.NullBool{Valid: true, Bool: true})
	dt._set("key33", &sql.NullByte{Valid: true, Byte: 1})
	dt._set("key34", &sql.NullFloat64{Valid: true, Float64: 3.14})
	dt._set("key35", &sql.NullInt16{Valid: true, Int16: 16})
	dt._set("key36", &sql.NullInt32{Valid: true, Int32: 32})
	dt._set("key37", &sql.NullInt64{Valid: true, Int64: 64})
	dt._set("key38", &sql.NullString{Valid: true, String: "koveychangestr"})
	dt._set("key39", &sql.NullTime{Valid: true, Time: tmpNow})

	return dt
}

func TestDataSet(t *testing.T) {
	dt := NewData()
	tmpStr1 := "kovey2"
	str1 := &tmpStr1
	dt.Set("key1", "kovey")
	dt.Set("key2", str1)
	dt.Set("key3", &str1)
	tmpInt1 := 1
	int1 := &tmpInt1
	dt.Set("key4", int1)
	dt.Set("key5", &int1)
	tmpInt2 := int8(8)
	int2 := &tmpInt2
	dt.Set("key6", int2)
	dt.Set("key7", &int2)
	tmpInt3 := int16(16)
	int3 := &tmpInt3
	dt.Set("key8", int3)
	dt.Set("key9", &int3)
	tmpInt4 := int32(32)
	int4 := &tmpInt4
	dt.Set("key10", int4)
	dt.Set("key11", &int4)
	tmpInt5 := int64(64)
	int5 := &tmpInt5
	dt.Set("key12", int5)
	dt.Set("key13", &int5)
	tmpuInt1 := uint(1)
	uint1 := &tmpuInt1
	dt.Set("key14", uint1)
	dt.Set("key15", &uint1)
	tmpuInt2 := uint8(8)
	uint2 := &tmpuInt2
	dt.Set("key16", uint2)
	dt.Set("key17", &uint2)
	tmpuInt3 := uint16(16)
	uint3 := &tmpuInt3
	dt.Set("key18", uint3)
	dt.Set("key19", &uint3)
	tmpuInt4 := uint32(32)
	uint4 := &tmpuInt4
	dt.Set("key20", uint4)
	dt.Set("key21", &uint4)
	tmpuInt5 := uint64(64)
	uint5 := &tmpuInt5
	dt.Set("key22", uint5)
	dt.Set("key23", &uint5)
	tmpBool := true
	bl := &tmpBool
	dt.Set("key24", bl)
	dt.Set("key25", &bl)
	tmpFloat := float32(32.32)
	float := &tmpFloat
	dt.Set("key26", float)
	dt.Set("key27", &float)
	tmpFloat1 := float64(64.64)
	float1 := &tmpFloat1
	dt.Set("key28", float1)
	dt.Set("key29", &float1)
	tmpNow := time.Now()
	now := &tmpNow
	dt.Set("key30", now)
	dt.Set("key31", &now)
	dt.Set("key32", &sql.NullBool{Valid: true, Bool: false})
	dt.Set("key33", &sql.NullByte{Valid: true, Byte: 8})
	dt.Set("key34", &sql.NullFloat64{Valid: true, Float64: 4.14})
	dt.Set("key35", &sql.NullInt16{Valid: true, Int16: 32})
	dt.Set("key36", &sql.NullInt32{Valid: true, Int32: 64})
	dt.Set("key37", &sql.NullInt64{Valid: true, Int64: 128})
	dt.Set("key38", &sql.NullString{Valid: true, String: "kovey"})
	dt.Set("key39", &sql.NullTime{Valid: true, Time: tmpNow})

	assert.Equal(t, "kovey", dt.Get("key1"))
	assert.Equal(t, "kovey2", dt.Get("key2"))
	assert.Equal(t, str1, dt.Get("key3"))
	assert.Equal(t, 1, dt.Get("key4"))
	assert.Equal(t, int1, dt.Get("key5"))
	assert.Equal(t, int8(8), dt.Get("key6"))
	assert.Equal(t, int2, dt.Get("key7"))
	assert.Equal(t, int16(16), dt.Get("key8"))
	assert.Equal(t, int3, dt.Get("key9"))
	assert.Equal(t, int32(32), dt.Get("key10"))
	assert.Equal(t, int4, dt.Get("key11"))
	assert.Equal(t, int64(64), dt.Get("key12"))
	assert.Equal(t, int5, dt.Get("key13"))
	assert.Equal(t, uint(1), dt.Get("key14"))
	assert.Equal(t, uint1, dt.Get("key15"))
	assert.Equal(t, uint8(8), dt.Get("key16"))
	assert.Equal(t, uint2, dt.Get("key17"))
	assert.Equal(t, uint16(16), dt.Get("key18"))
	assert.Equal(t, uint3, dt.Get("key19"))
	assert.Equal(t, uint32(32), dt.Get("key20"))
	assert.Equal(t, uint4, dt.Get("key21"))
	assert.Equal(t, uint64(64), dt.Get("key22"))
	assert.Equal(t, uint5, dt.Get("key23"))
	assert.Equal(t, true, dt.Get("key24"))
	assert.Equal(t, bl, dt.Get("key25"))
	assert.Equal(t, float32(32.32), dt.Get("key26"))
	assert.Equal(t, float, dt.Get("key27"))
	assert.Equal(t, float64(64.64), dt.Get("key28"))
	assert.Equal(t, float1, dt.Get("key29"))
	assert.Equal(t, *now, dt.Get("key30"))
	assert.Equal(t, now, dt.Get("key31"))
	assert.Equal(t, sql.NullBool{Valid: true, Bool: false}, dt.Get("key32"))
	assert.Equal(t, sql.NullByte{Valid: true, Byte: 8}, dt.Get("key33"))
	assert.Equal(t, sql.NullFloat64{Valid: true, Float64: 4.14}, dt.Get("key34"))
	assert.Equal(t, sql.NullInt16{Valid: true, Int16: 32}, dt.Get("key35"))
	assert.Equal(t, sql.NullInt32{Valid: true, Int32: 64}, dt.Get("key36"))
	assert.Equal(t, sql.NullInt64{Valid: true, Int64: 128}, dt.Get("key37"))
	assert.Equal(t, sql.NullString{Valid: true, String: "kovey"}, dt.Get("key38"))
	assert.Equal(t, sql.NullTime{Valid: true, Time: tmpNow}, dt.Get("key39"))

	old := _testGetData()
	old.Range(func(key string, val any) {
		assert.True(t, dt.Changed(key, val))
	})
}

func TestDataFrom(t *testing.T) {
	old := NewData()
	old.Set("key1", 1)
	old.Set("key2", "kovey")

	n := NewData()
	n.From(old)
	assert.Equal(t, 1, n.Get("key1"))
	assert.Equal(t, "kovey", n.Get("key2"))
	assert.Equal(t, []string{"key1", "key2"}, n.Keys())
	var keys []string
	var vals []any
	n.Range(func(key string, val any) {
		keys = append(keys, key)
		vals = append(vals, val)
	})
	assert.Equal(t, []string{"key1", "key2"}, keys)
	assert.Equal(t, []any{1, "kovey"}, vals)
	assert.False(t, n.Empty())
}
