/*
 * Copyright 2023 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cast_test

import (
	"bytes"
	"errors"
	"github.com/lvan100/cast"
	"github.com/lvan100/cast/internal/assert"
	"html/template"
	"strconv"
	"testing"
	"time"
)

func BenchmarkToString(b *testing.B) {
	//int/strconv-8    419501868 2.87 ns/op
	//int/cast-8       60869038  18.2 ns/op
	b.Run("int", func(b *testing.B) {
		v := 10
		b.Run("strconv", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = strconv.Itoa(v)
			}
		})
		b.Run("cast", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = cast.ToString(v)
			}
		})
	})

}

func TestToString(t *testing.T) {

	assert.Equal(t, cast.ToString(nil), "")

	assert.Equal(t, cast.ToString(int(3)), "3")
	assert.Equal(t, cast.ToString(int8(3)), "3")
	assert.Equal(t, cast.ToString(int16(3)), "3")
	assert.Equal(t, cast.ToString(int32(3)), "3")
	assert.Equal(t, cast.ToString(int64(3)), "3")
	assert.Equal(t, cast.ToString(cast.IntPtr(3)), "3")
	assert.Equal(t, cast.ToString(cast.Int8Ptr(3)), "3")
	assert.Equal(t, cast.ToString(cast.Int16Ptr(3)), "3")
	assert.Equal(t, cast.ToString(cast.Int32Ptr(3)), "3")
	assert.Equal(t, cast.ToString(cast.Int64Ptr(3)), "3")

	assert.Equal(t, cast.ToString(uint(3)), "3")
	assert.Equal(t, cast.ToString(uint8(3)), "3")
	assert.Equal(t, cast.ToString(uint16(3)), "3")
	assert.Equal(t, cast.ToString(uint32(3)), "3")
	assert.Equal(t, cast.ToString(uint64(3)), "3")
	assert.Equal(t, cast.ToString(cast.UintPtr(3)), "3")
	assert.Equal(t, cast.ToString(cast.Uint8Ptr(3)), "3")
	assert.Equal(t, cast.ToString(cast.Uint16Ptr(3)), "3")
	assert.Equal(t, cast.ToString(cast.Uint32Ptr(3)), "3")
	assert.Equal(t, cast.ToString(cast.Uint64Ptr(3)), "3")

	assert.Equal(t, cast.ToString(float32(3)), "3")
	assert.Equal(t, cast.ToString(float64(3)), "3")
	assert.Equal(t, cast.ToString(cast.Float32Ptr(3)), "3")
	assert.Equal(t, cast.ToString(cast.Float64Ptr(3)), "3")

	assert.Equal(t, cast.ToString("3"), "3")
	assert.Equal(t, cast.ToString(cast.StringPtr("3")), "3")

	assert.Equal(t, cast.ToString(true), "true")
	assert.Equal(t, cast.ToString(false), "false")
	assert.Equal(t, cast.ToString(cast.BoolPtr(true)), "true")
	assert.Equal(t, cast.ToString(cast.BoolPtr(false)), "false")

	assert.Equal(t, cast.ToString([]byte("3")), "3")
	assert.Equal(t, cast.ToString(template.HTML("3")), "3")
	assert.Equal(t, cast.ToString(template.URL("3")), "3")
	assert.Equal(t, cast.ToString(template.JS("3")), "3")
	assert.Equal(t, cast.ToString(template.CSS("3")), "3")
	assert.Equal(t, cast.ToString(template.HTMLAttr("3")), "3")
	assert.Equal(t, cast.ToString(bytes.NewBufferString("abc")), "abc")
	assert.Equal(t, cast.ToString(errors.New("abc")), "abc")

	type String string
	assert.Equal(t, cast.ToString(String("abc")), "abc")

	var a []string = nil
	assert.Equal(t, cast.ToString(a), "")
	a = append(a, "1.2x")
	assert.Equal(t, cast.ToString(a), "[\"1.2x\"]")
	assert.Equal(t, cast.ToString([]int{2}), "[2]")

	var time1 time.Time
	assert.Equal(t, cast.ToString(time1), "0001-01-01 00:00:00 +0000 UTC")
	var time2 *time.Time
	assert.Equal(t, cast.ToString(time2), "")

	type Stu struct {
		Name string
	}
	var stu *Stu
	assert.Equal(t, cast.ToString(stu), "")
	stu = &Stu{Name: "test"}
	assert.Equal(t, cast.ToString(stu), "{\"Name\":\"test\"}")

	var f = func(a, b int) {}
	_ = cast.ToString(&f)

}
