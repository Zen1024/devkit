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
	value := reflect.ValueOf(dest)

	typkind := typ.Kind()
	if typkind != reflect.Ptr {
		return fmt.Errorf("invalid scan obj type:%s", typkind.String())
	}
	if value.IsNil() {
		return nil
	}
	valelem := value.Elem()
	typelem := typ.Elem()
	fmt.Println(value, valelem, typ.String())
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
			valobj, ok := val[key]
			if !ok {
				continue
			}
			kind := field.Kind()
			if isbasekind(kind) {
				field.Set(reflect.ValueOf(valobj))
			} else {
				if kind == reflect.Ptr {
					if mpval, ok := valobj.(map[string]interface{}); ok {
						if err := scanstruct(mpval, tag, field.Interface()); err != nil {
							return err
						}
					}
				} else if kind == reflect.Slice {
					//@TODO
					fieldtyp := typfield.Type.Elem()
					if mpvals, ok := valobj.([]interface{}); ok {
						slen := len(mpvals)
						slice := reflect.MakeSlice(typfield.Type, slen, slen)
						for _, impval := range mpvals {
							if mpval, ok := impval.(map[string]interface{}); ok {

								obj := reflect.New(fieldtyp.Elem())
								if err := scanstruct(mpval, tag, obj.Interface()); err != nil {
									return err
								}
								reflect.Append(slice, obj)
							}
						}
						field.Set(slice)
					}
				}
			}
		}
	}

	return nil
}

func isbasekind(kind reflect.Kind) bool {
	return (kind <= reflect.Float64) || (kind == reflect.String)
}
