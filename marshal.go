package phpcrudapi

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type tableResultSet struct {
	Columns []string            `json:"columns"`
	Records [][]json.RawMessage `json:"records"`
}

// Unmarshal decodes the efficient columns/records encoding into
// an array of go types
func Unmarshal(data []byte, v interface{}) error {
	st, table, err := getTypeAndTableSlice(v)
	if err != nil {
		return err
	}

	// parse json
	var inter map[string]tableResultSet
	err = json.Unmarshal(data, &inter)
	if err != nil {
		return err
	}

	return unmarshalSlice(st, table, inter, v)
}

func unmarshalSlice(st reflect.Type, table string, inter map[string]tableResultSet, v interface{}) error {
	// check that table is in response
	rs, ok := inter[table]
	if !ok {
		return fmt.Errorf("Specified table %q not in response, found: %#v", table, inter)
	}

	// build column index
	index := make(map[string]int)
	for i, column := range rs.Columns {
		index[column] = i
	}

	// create new array
	size := len(rs.Records)
	arr := reflect.MakeSlice(st, size, size)
	reflect.ValueOf(v).Elem().Set(arr)

	// set array values
	for i, row := range rs.Records {
		elem := arr.Index(i)
		elemT := elem.Type().Elem()
		item := reflect.New(elemT)
		elem.Set(item)

		for j := 0; j < elemT.NumField(); j++ {
			tag := item.Type().Elem().Field(j).Tag
			column := tag.Get("json")
			columnI := index[column]
			f := item.Elem().Field(j).Addr().Interface()

			if rel := tag.Get("relation"); rel != "" {
				// unmarshal reference
				var ref int64
				err := json.Unmarshal(row[columnI], &ref)
				if err != nil {
					return err
				}

				// get type
				switch rel {
				case "belongs-to":
					otmst, otmtable, err := getTypeAndTable(f)
					if err != nil {
						return err
					}
					otmstSlice := reflect.SliceOf(otmst)
					data := reflect.New(otmstSlice)
					err = unmarshalSlice(otmstSlice, otmtable, inter, data.Interface())
					if err != nil {
						return err
					}
					arr := data.Elem()
					for i := 0; i < arr.Len(); i++ {
						bi := arr.Index(i).Elem()
						if bi.FieldByName("ID").Int() == ref {
							// found matching element
							item.Elem().Field(j).Set(arr.Index(i))
						}
					}
					// TODO: cache & filter
				case "one-to-many":
					otmst, otmtable, err := getTypeAndTableSlice(f)
					if err != nil {
						return err
					}
					err = unmarshalSlice(otmst, otmtable, inter, f)
					if err != nil {
						return err
					}
					// TODO: cache & filter
				case "many-to-many":

				}

				continue
			}

			err := json.Unmarshal(row[columnI], f)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
