package phpcrudapi

import (
	"fmt"
	"reflect"
)

func getTypeAndTableSlice(v interface{}) (reflect.Type, string, error) {
	value := reflect.ValueOf(v).Elem()

	// get type of marshaled value
	if value.Kind() != reflect.Slice {
		return nil, "", fmt.Errorf("Only slices can be serialized")
	}

	pst := value.Type().Elem()
	if pst.Kind() != reflect.Ptr {
		return nil, "", fmt.Errorf("Only slices of pointers can be serialized")
	}

	// get collection name from type
	table, err := collectionName(pst.Elem())
	return value.Type(), table, err
}

func getTypeAndTable(v interface{}) (reflect.Type, string, error) {
	st := reflect.ValueOf(v).Elem()
	if st.Kind() != reflect.Ptr {
		return nil, "", fmt.Errorf("Only pointers can be serialized")
	}

	table, err := collectionName(st.Type().Elem())
	return st.Type(), table, err
}

func collectionName(t reflect.Type) (string, error) {
	field, ok := t.FieldByName("ID")
	if !ok {
		return "", fmt.Errorf("Field ID needs to exist")
	}
	table := field.Tag.Get("collection")
	if table == "" {
		return "", fmt.Errorf("Field ID needs to define collection tag with the name of the collection (table name)")
	}
	return table, nil
}
