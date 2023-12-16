package cast

import (
	"encoding/json"

	"github.com/spf13/cast"
	_ "github.com/spf13/cast"
)

// To 将 i 转换为 T 类型的值。
func To[T any](i interface{}) (T, error) {
	var t T
	if err := to(i, &t); err != nil {
		return t, err
	}
	return t, nil
}

func to(i any, v any) error {
	switch p := v.(type) {
	case *int:
		*p = cast.ToInt(i)
	case *string:
		*p = cast.ToString(i)
	default:
		b, err := json.Marshal(i)
		if err != nil {
			return err
		}
		err = json.Unmarshal(b, v)
		if err != nil {
			return err
		}
	}
	return nil
}
