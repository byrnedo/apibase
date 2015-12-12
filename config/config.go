package config

import (
	"errors"
	"fmt"
	"github.com/byrnedo/apibase/helpers/stringhelp"
	"github.com/byrnedo/typesafe-config/parse"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

func ParseFile(path string) (*parse.Tree, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("Failed to read config file")
	}
	tree, err := Parse(bytes)
	return tree, err
}

func Parse(configFileData []byte) (tree *parse.Tree, err error) {
	tree, err = parse.Parse("config", string(configFileData))
	return
}

func Populate(targetPtr interface{}, conf *parse.Config) {
	fmt.Println("Populate:", targetPtr)
	setValue(reflect.ValueOf(targetPtr), conf, "", "")
}

func configFieldNamer(field reflect.StructField, prefix string) (name string, defaultVal string) {
	t := field.Tag.Get("config")
	tArr := strings.Split(t, ",")

	name = stringhelp.ToDashCase(field.Name)
	fmt.Println(name)

	if len(tArr) > 0 && len(tArr[0]) > 0 {
		switch tArr[0] {
		case "-":
			return "", ""
		default:
			name = tArr[0]
		}
	}

	if len(tArr) > 1 {
		defaultVal = tArr[1]
	}

	if len(prefix) > 0 {
		prefix = prefix + "."
	}
	fmt.Println("Field name: " + prefix + name)
	return prefix + name, defaultVal
}

func setValue(field reflect.Value, conf *parse.Config, configName string, defaultVal string) {
	if field.Kind() != reflect.Ptr {
		panic("Not a pointer value " + field.Kind().String())
	}

	field = reflect.Indirect(field)
	if !field.CanSet() {
		return
	}
	switch field.Kind() {
	case reflect.Struct:
		itemType := reflect.TypeOf(field.Interface())

		for i := 0; i < field.NumField(); i++ {
			configFieldName, defaultVal := configFieldNamer(itemType.Field(i), configName)

			setValue(field.Field(i).Addr(), conf, configFieldName, defaultVal)
		}
	case reflect.Ptr:
	case reflect.Bool:
		defaultBool, _ := strconv.ParseBool(defaultVal)
		boolVal := conf.GetDefaultBool(configName, defaultBool)
		field.SetBool(boolVal)
	case reflect.String:
		strVal := conf.GetDefaultString(configName, defaultVal)
		field.SetString(strVal)
	case reflect.Int:
		defaultInt, _ := strconv.Atoi(defaultVal)
		intVal := conf.GetDefaultInt(configName, int64(defaultInt))
		field.SetInt(intVal)
	case reflect.Int8:
		defaultInt, _ := strconv.ParseInt(defaultVal, 10, 8)
		intVal := conf.GetDefaultInt(configName, defaultInt)
		field.SetInt(intVal)
	case reflect.Int16:
		defaultInt, _ := strconv.ParseInt(defaultVal, 10, 16)
		intVal := conf.GetDefaultInt(configName, defaultInt)
		field.SetInt(intVal)
	case reflect.Int32:
		defaultInt, _ := strconv.ParseInt(defaultVal, 10, 32)
		intVal := conf.GetDefaultInt(configName, defaultInt)
		field.SetInt(intVal)
	case reflect.Int64:
		defaultInt, _ := strconv.ParseInt(defaultVal, 10, 64)
		intVal := conf.GetDefaultInt(configName, defaultInt)
		field.SetInt(intVal)
	case reflect.Uint:
		defaultInt, _ := strconv.Atoi(defaultVal)
		intVal := conf.GetDefaultUInt(configName, uint64(defaultInt))
		field.SetUint(intVal)
	case reflect.Uint8:
		defaultInt, _ := strconv.ParseUint(defaultVal, 10, 8)
		intVal := conf.GetDefaultUInt(configName, defaultInt)
		field.SetUint(intVal)
	case reflect.Uint16:
		defaultInt, _ := strconv.ParseUint(defaultVal, 10, 16)
		intVal := conf.GetDefaultUInt(configName, defaultInt)
		field.SetUint(intVal)
	case reflect.Uint32:
		defaultInt, _ := strconv.ParseUint(defaultVal, 10, 32)
		intVal := conf.GetDefaultUInt(configName, defaultInt)
		field.SetUint(intVal)
	case reflect.Uint64:
		defaultInt, _ := strconv.ParseUint(defaultVal, 10, 64)
		intVal := conf.GetDefaultUInt(configName, defaultInt)
		field.SetUint(intVal)
	case reflect.Float32:
		defaultInt, _ := strconv.ParseFloat(defaultVal, 32)
		intVal := conf.GetDefaultFloat(configName, defaultInt)
		field.SetFloat(intVal)
	case reflect.Float64:
		defaultInt, _ := strconv.ParseFloat(defaultVal, 64)
		intVal := conf.GetDefaultFloat(configName, defaultInt)
		field.SetFloat(intVal)
	default:
	}

}
