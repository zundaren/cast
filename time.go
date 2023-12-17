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
	"time"
)

// ToTime casts an any to a time.Time.
// When type is clear, it is recommended to use standard library functions.
func ToTime(i any, opts ...Option) time.Time {
	v, _ := ToTimeE(i, opts...)
	return v
}

// ToTimeE casts an any to a time.Time.
// When type is clear, it is recommended to use standard library functions.
func ToTimeE(i any, opts ...Option) (time.Time, error) {
	switch v := i.(type) {
	case nil:
		return time.Time{}, nil
	case int:
		return parseTimestamp(int64(v), opts...), nil
	case int8:
		return parseTimestamp(int64(v), opts...), nil
	case int16:
		return parseTimestamp(int64(v), opts...), nil
	case int32:
		return parseTimestamp(int64(v), opts...), nil
	case int64:
		return parseTimestamp(v, opts...), nil
	case *int:
		return parseTimestamp(int64(*v), opts...), nil
	case *int8:
		return parseTimestamp(int64(*v), opts...), nil
	case *int16:
		return parseTimestamp(int64(*v), opts...), nil
	case *int32:
		return parseTimestamp(int64(*v), opts...), nil
	case *int64:
		return parseTimestamp(*v, opts...), nil
	case uint:
		return parseTimestamp(int64(v), opts...), nil
	case uint8:
		return parseTimestamp(int64(v), opts...), nil
	case uint16:
		return parseTimestamp(int64(v), opts...), nil
	case uint32:
		return parseTimestamp(int64(v), opts...), nil
	case uint64:
		return parseTimestamp(int64(v), opts...), nil
	case *uint:
		return parseTimestamp(int64(*v), opts...), nil
	case *uint8:
		return parseTimestamp(int64(*v), opts...), nil
	case *uint16:
		return parseTimestamp(int64(*v), opts...), nil
	case *uint32:
		return parseTimestamp(int64(*v), opts...), nil
	case *uint64:
		return parseTimestamp(int64(*v), opts...), nil
	case float32:
		return parseTimestamp(float64(v), opts...), nil
	case float64:
		return parseTimestamp(v, opts...), nil
	case *float32:
		return parseTimestamp(float64(*v), opts...), nil
	case *float64:
		return parseTimestamp(*v, opts...), nil
	case string:
		return parseFormatTime(v, opts...)
	case *string:
		return parseFormatTime(*v, opts...)
	case time.Time:
		return v, nil
	case *time.Time:
		return *v, nil
	default:
		return time.Time{}, fmt.Errorf("unable to cast type (%T) to Time", i)
	}
}

func parseTimestamp[T int64 | float64](v T, opts ...Option) time.Time {
	arg := OptionArg{
		TimeFormat: "ns",
	}
	for _, opt := range opts {
		opt(&arg)
	}
	unitN, _ := unitMap[arg.TimeFormat]
	i := int64(float64(v) * float64(unitN))
	return time.Unix(i/int64(time.Second), i%int64(time.Second))
}

func parseFormatTime(v string, opts ...Option) (time.Time, error) {
	if d, err := time.ParseDuration(v); err == nil {
		return time.Unix(int64(d/time.Second), int64(d%time.Second)), nil
	}
	arg := OptionArg{
		TimeFormat: "2006-01-02 15:04:05 -0700",
	}
	for _, opt := range opts {
		opt(&arg)
	}
	return time.Parse(arg.TimeFormat, v)
}
