package phpcrudapi

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Unmarshal decodes the efficient columns/records encoding into
// an array of go types
func Unmarshal(data []byte, v interface{}) error {
	st, table, err := getTypeAndTableSlice(v)
	if err != nil {
		return err
	}

	// parse json
	var rs resultSet
	return rs.unmarshal(data, st, table, v)
}

type resultSet struct {
	Data        map[string]tableResultSet
	objectCache map[string]reflect.Value
}

type tableResultSet struct {
	Columns []string            `json:"columns"`
	Records [][]json.RawMessage `json:"records"`
}

func (rs *resultSet) unmarshal(data []byte, st reflect.Type, table string, v interface{}) error {
	err := json.Unmarshal(data, &rs.Data)
	if err != nil {
		return err
	}
	rs.objectCache = make(map[string]reflect.Value)
	return rs.unmarshalSlice(st, table, v)
}

func (rs *resultSet) unmarshalSlice(st reflect.Type, table string, v interface{}) error {
	// check that table is in response
	trs, ok := rs.Data[table]
	if !ok {
		return fmt.Errorf("Specified table %q not in response, found: %#v", table, rs.Data)
	}

	// build column index
	index := make(map[string]int)
	for i, column := range trs.Columns {
		index[column] = i
	}

	// create new array
	size := len(trs.Records)
	arr := reflect.MakeSlice(st, size, size)
	reflect.ValueOf(v).Elem().Set(arr)

	// set array values
	for i, row := range trs.Records {
		elem := arr.Index(i)
		elemT := elem.Type().Elem()
		item := reflect.New(elemT)
		elem.Set(item)

		for j := 0; j < elemT.NumField(); j++ {
			tag := item.Type().Elem().Field(j).Tag
			column := tag.Get("json")
			columnI := index[column]
			field := item.Elem().Field(j).Addr().Interface()

			// unmarshal relations
			if rel := tag.Get("relation"); rel != "" {
				// unmarshal reference
				var ref interface{}
				err := json.Unmarshal(row[columnI], &ref)
				if err != nil {
					return err
				}

				// get type
				switch rel {
				case "belongs-to":
					t, table, err := getTypeAndTable(field)
					if err != nil {
						return err
					}
					val, err := rs.findObject(t, table, ref)
					if err != nil {
						return err
					}
					item.Elem().Field(j).Set(val)
				case "one-to-many":
					t, table, err := getTypeAndTableSlice(field)
					if err != nil {
						return err
					}
					val, err := rs.getFromCache(t.Elem(), table)
					if err != nil {
						return err
					}
					// TODO filter elememts of the list using foreign key
					item.Elem().Field(j).Set(val.Elem())
				case "many-to-many":
					// TODO
				default:
					panic(fmt.Errorf("unknown relation %q for type %s", rel, st))
				}

				// skip regular unmarshal
				continue
			}

			// unmarshal simple values to the fields
			err := json.Unmarshal(row[columnI], field)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (rs *resultSet) findObject(t reflect.Type, table string, ref interface{}) (reflect.Value, error) {
	arrPtr, err := rs.getFromCache(t, table)
	if err != nil {
		return reflect.Value{}, err
	}
	arr := arrPtr.Elem()

	for i := 0; i < arr.Len(); i++ {
		bi := arr.Index(i).Elem()
		if bi.FieldByName("ID").Interface() == ref {
			// found matching element
			val := arr.Index(i)
			return val, nil
		}
	}

	return reflect.Value{}, fmt.Errorf("ID %v not found in table %q data (possibly not using float64 or string for ID types)!", ref, table)
}

func (rs *resultSet) getFromCache(t reflect.Type, table string) (reflect.Value, error) {
	arrPtr, ok := rs.objectCache[table]
	if !ok {
		slice := reflect.SliceOf(t)
		arrPtr = reflect.New(slice)
		err := rs.unmarshalSlice(slice, table, arrPtr.Interface())
		if err != nil {
			return reflect.Value{}, err
		}
		rs.objectCache[table] = arrPtr
	}

	return arrPtr, nil
}
