package cast

import "time"

func BoolPtr(s bool) *bool          { return &s }
func IntPtr(s int) *int             { return &s }
func Int8Ptr(s int8) *int8          { return &s }
func Int16Ptr(s int16) *int16       { return &s }
func Int32Ptr(s int32) *int32       { return &s }
func Int64Ptr(s int64) *int64       { return &s }
func UintPtr(s uint) *uint          { return &s }
func Uint8Ptr(s uint8) *uint8       { return &s }
func Uint16Ptr(s uint16) *uint16    { return &s }
func Uint32Ptr(s uint32) *uint32    { return &s }
func Uint64Ptr(s uint64) *uint64    { return &s }
func Float32Ptr(s float32) *float32 { return &s }
func Float64Ptr(s float64) *float64 { return &s }
func StringPtr(s string) *string    { return &s }

type OptionArg struct {
	TimeFormat string
}

type Option func(arg *OptionArg)

func TimeFormat(format string) Option {
	return func(arg *OptionArg) {
		arg.TimeFormat = format
	}
}

// To 将 i 转换为 T 类型的值。
func To[T any](i interface{}, opts ...Option) (T, error) {
	var t T
	if err := to(i, &t); err != nil {
		return t, err
	}
	return t, nil
}

func to(i any, v any, opts ...Option) error {
	var err error
	switch p := v.(type) {
	case *bool:
		*p, err = ToBoolE(i)
	case *int:
		var r int64
		r, err = ToInt64E(i)
		*p = int(r)
	case *int8:
		var r int64
		r, err = ToInt64E(i)
		*p = int8(r)
	case *int16:
		var r int64
		r, err = ToInt64E(i)
		*p = int16(r)
	case *int32:
		var r int64
		r, err = ToInt64E(i)
		*p = int32(r)
	case *int64:
		var r int64
		r, err = ToInt64E(i)
		*p = r
	case *uint:
		var r uint64
		r, err = ToUint64E(i)
		*p = uint(r)
	case *uint8:
		var r uint64
		r, err = ToUint64E(i)
		*p = uint8(r)
	case *uint16:
		var r uint64
		r, err = ToUint64E(i)
		*p = uint16(r)
	case *uint32:
		var r uint64
		r, err = ToUint64E(i)
		*p = uint32(r)
	case *uint64:
		var r uint64
		r, err = ToUint64E(i)
		*p = r
	case *float32:
		var r float64
		r, err = ToFloat64E(i)
		*p = float32(r)
	case *float64:
		var r float64
		r, err = ToFloat64E(i)
		*p = r
	case *string:
		*p = ToString(i)
		err = nil
	case *time.Duration:
		var r time.Duration
		r, err = ToDurationE(i, opts...)
		*p = r
	case *time.Time:
		var r time.Time
		r, err = ToTimeE(i, opts...)
		*p = r
	default:
		return JSON.Convert(i, v)
	}
	return err
}
