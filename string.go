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

package cast

import (
	"encoding/json"
	"fmt"
	"html/template"
	"reflect"
	"strconv"
	"time"
)

// ToString casts an any to a string.
// When type is clear, it is recommended to use standard library functions.
func ToString(i any) string {
	switch s := i.(type) {
	case nil:
		return ""
	case int:
		return strconv.Itoa(s)
	case int8:
		return strconv.FormatInt(int64(s), 10)
	case int16:
		return strconv.FormatInt(int64(s), 10)
	case int32:
		return strconv.Itoa(int(s))
	case int64:
		return strconv.FormatInt(s, 10)
	case *int:
		return strconv.Itoa(*s)
	case *int8:
		return strconv.FormatInt(int64(*s), 10)
	case *int16:
		return strconv.FormatInt(int64(*s), 10)
	case *int32:
		return strconv.Itoa(int(*s))
	case *int64:
		return strconv.FormatInt(*s, 10)
	case uint:
		return strconv.FormatUint(uint64(s), 10)
	case uint8:
		return strconv.FormatUint(uint64(s), 10)
	case uint16:
		return strconv.FormatUint(uint64(s), 10)
	case uint32:
		return strconv.FormatUint(uint64(s), 10)
	case uint64:
		return strconv.FormatUint(s, 10)
	case *uint:
		return strconv.FormatUint(uint64(*s), 10)
	case *uint8:
		return strconv.FormatUint(uint64(*s), 10)
	case *uint16:
		return strconv.FormatUint(uint64(*s), 10)
	case *uint32:
		return strconv.FormatUint(uint64(*s), 10)
	case *uint64:
		return strconv.FormatUint(*s, 10)
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64)
	case *float32:
		return strconv.FormatFloat(float64(*s), 'f', -1, 32)
	case *float64:
		return strconv.FormatFloat(*s, 'f', -1, 64)
	case string:
		return s
	case *string:
		return *s
	case bool:
		return strconv.FormatBool(s)
	case *bool:
		return strconv.FormatBool(*s)
	case []byte:
		return string(s)
	case template.HTML:
		return string(s)
	case template.URL:
		return string(s)
	case template.JS:
		return string(s)
	case template.CSS:
		return string(s)
	case template.HTMLAttr:
		return string(s)
	case *time.Time:
		if s == nil {
			return ""
		}
		return s.String()
	case fmt.Stringer:
		return s.String()
	case error:
		return s.Error()
	default:
		rv := reflect.ValueOf(s)
		kind := rv.Kind()

		switch kind {
		case reflect.Chan, reflect.Map, reflect.Slice, reflect.Func,
			reflect.Ptr, reflect.Interface, reflect.UnsafePointer:
			if rv.IsNil() {
				return ""
			}
			if kind == reflect.Ptr {
				return ToString(rv.Elem().Interface())
			}
		case reflect.String:
			return rv.String()
		}

		if jb, err := json.Marshal(s); err == nil {
			return string(jb)
		}
		return fmt.Sprint(s)
	}
}
