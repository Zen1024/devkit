package mapstruct

import (
	"fmt"
	"reflect"
	"strings"
)

func Mapstruct(obj interface{}, tag string) map[string]interface{} {
	return mapstruct(obj, tag)
}

func mapstruct(obj interface{}, tag string) map[string]interface{} {
	objtyp := reflect.TypeOf(obj)
	objval := reflect.ValueOf(obj)
	objkind := objtyp.Kind()
	if objkind != reflect.Ptr && objkind != reflect.Struct {
		panic("invalid kind:" + objkind.String())
	}
	valelem := objval
	typelem := objtyp
	if objkind == reflect.Ptr {
		valelem = objval.Elem()
		typelem = objtyp.Elem()
	}

	fieldlen := valelem.NumField()
	result := map[string]interface{}{}

	for i := 0; i < fieldlen; i++ {
		valfield := valelem.Field(i)
		typfield := typelem.Field(i)
		key := typfield.Name
		tagval := typfield.Tag.Get(tag)

		if tagval != "" {
			valarr := strings.Split(tagval, ",")
			key = valarr[0]
		} else {
			continue
		}
		val := valfield.Interface()
		kind := valfield.Kind()
		switch kind {
		case reflect.Struct:
			val = mapstruct(val, tag)
		case reflect.Ptr:
			if valfield.IsNil() {
				continue
			}
			val = mapstruct(val, tag)
		case reflect.Array:
			//@TODO:
		case reflect.Slice:
			if valfield.IsNil() {
				continue
			}
			sl := valfield.Len()
			arr := make([]interface{}, sl)
			for i := 0; i < sl; i++ {
				elem := valfield.Index(i)
				ekind := elem.Kind()
				if ekind != reflect.Ptr && ekind != reflect.Struct {
					arr[i] = elem.Interface()
				} else {
					arr[i] = mapstruct(elem.Interface(), tag)
				}
			}
			val = arr

		}
		result[key] = val
	}
	return result

}

func ScanStruct(val map[string]interface{}, tag string, dest interface{}) error {
	return scanstruct(val, tag, dest)
}

func scanstruct(val map[string]interface{}, tag string, dest interface{}) error {
	if len(val) == 0 {
		return nil
	}
	typ := reflect.TypeOf(dest)
	typelem := typ
	value := reflect.ValueOf(dest)
	valelem := value
	typkind := typ.Kind()
	if typkind != reflect.Ptr {
		return fmt.Errorf("invalid scan obj type:%s", typkind.String())
	}

	if !valelem.CanAddr() {
		return nil
	}

	fieldlen := valelem.NumField()
	for i := 0; i < fieldlen; i++ {
		field := valelem.Field(i)
		typfield := typelem.Field(i)
		key := typfield.Name
		tagval := typfield.Tag.Get(tag)

		if tagval != "" {
			valarr := strings.Split(tagval, ",")
			key = valarr[0]
		} else {
			continue
		}
		if field.CanSet() {
			obj, ok := val[key]
			if !ok {
				continue
			}
			kind := field.Kind()
			if isbasekind(kind) {
				field.Set(reflect.ValueOf(obj))
			} else {
				if kind == reflect.Ptr {
					if mpval, ok := obj.(map[string]interface{}); ok {
						if err := scanstruct(mpval, tag, field.Interface()); err != nil {
							return err
						}
					}
				} else if kind == reflect.Slice {
					//@TODO
				}
			}
		}
	}

	return nil
}

func isbasekind(kind reflect.Kind) bool {
	return kind <= reflect.Float64
}
