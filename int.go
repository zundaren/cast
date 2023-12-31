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
	"fmt"
	"strconv"
)

// ToInt casts an any to an int.
// When type is clear, it is recommended to use standard library functions.
func ToInt(i any) int {
	v, _ := ToInt64E(i)
	return int(v)
}

// ToInt8 casts an any to an int8.
// When type is clear, it is recommended to use standard library functions.
func ToInt8(i any) int8 {
	v, _ := ToInt64E(i)
	return int8(v)
}

// ToInt16 casts an any to an int16.
// When type is clear, it is recommended to use standard library functions.
func ToInt16(i any) int16 {
	v, _ := ToInt64E(i)
	return int16(v)
}

// ToInt32 casts an any to an int32.
// When type is clear, it is recommended to use standard library functions.
func ToInt32(i any) int32 {
	v, _ := ToInt64E(i)
	return int32(v)
}

// ToInt64 casts an any to an int64.
// When type is clear, it is recommended to use standard library functions.
func ToInt64(i any) int64 {
	v, _ := ToInt64E(i)
	return v
}

// ToInt64E casts an any to an int64.
// When type is clear, it is recommended to use standard library functions.
func ToInt64E(i any) (int64, error) {
	switch s := i.(type) {
	case nil:
		return 0, nil
	case int:
		return int64(s), nil
	case int8:
		return int64(s), nil
	case int16:
		return int64(s), nil
	case int32:
		return int64(s), nil
	case int64:
		return s, nil
	case *int:
		return int64(*s), nil
	case *int8:
		return int64(*s), nil
	case *int16:
		return int64(*s), nil
	case *int32:
		return int64(*s), nil
	case *int64:
		return *s, nil
	case uint:
		return int64(s), nil
	case uint8:
		return int64(s), nil
	case uint16:
		return int64(s), nil
	case uint32:
		return int64(s), nil
	case uint64:
		return int64(s), nil
	case *uint:
		return int64(*s), nil
	case *uint8:
		return int64(*s), nil
	case *uint16:
		return int64(*s), nil
	case *uint32:
		return int64(*s), nil
	case *uint64:
		return int64(*s), nil
	case float32:
		return int64(s), nil
	case float64:
		return int64(s), nil
	case *float32:
		return int64(*s), nil
	case *float64:
		return int64(*s), nil
	case string:
		return strconv.ParseInt(s, 0, 0)
	case *string:
		return strconv.ParseInt(*s, 0, 0)
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case *bool:
		if *s {
			return 1, nil
		}
		return 0, nil
	}
	return 0, fmt.Errorf("unable to cast type (%T) to int64", i)
}
