package common

import "reflect"

// IsEmptyStruct 判断结构体是否为空
func IsEmptyStruct(s interface{}) bool {
	v := reflect.ValueOf(s)

	// 只有结构体类型才能判断是否为空
	if v.Kind() != reflect.Struct {
		return false
	}

	// 遍历结构体的字段
	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)

		// 判断字段是否为零值或空值
		switch fieldValue.Kind() {
		case reflect.String:
			if fieldValue.String() != "" {
				return false
			}
		case reflect.Array, reflect.Slice, reflect.Map:
			if !fieldValue.IsNil() && fieldValue.Len() != 0 {
				return false
			}
		default:
			// 使用 IsZero 方法判断零值
			if !fieldValue.IsZero() {
				return false
			}
		}
	}

	return true
}
