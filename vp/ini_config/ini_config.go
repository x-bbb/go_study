package iniconfig

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Marshal(v interface{}) (data []byte, err error) {
	rt := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)
	if rt.Kind() != reflect.Struct {
		err = errors.New("please pass a struct")
		return
	}

	for i := 0; i < rt.NumField(); i++ {
		sectionType := rt.Field(i).Type
		if sectionType.Kind() != reflect.Struct {
			err = fmt.Errorf("section %s must  be a struct", rt.Field(i).Name)
			return
		}

		sectionName := sectionType.Name()
		section := fmt.Sprintf("[%s]\n", sectionName)
		data = append(data, []byte(section)...)

		for j := 0; j < sectionType.NumField(); j++ {
			// 下层struct的tag
			tagName := sectionType.Field(j).Tag.Get("ini")
			if len(tagName) == 0 {
				tagName = sectionType.Field(j).Name
			}

			tagVal := rv.Field(i).Field(j)
			item := fmt.Sprintf("%s = %v\n", tagName, tagVal.Interface())
			data = append(data, []byte(item)...)
		}

	}

	fmt.Println(string(data))

	return
}

func UnMarshal(data []byte, v interface{}) (err error) {

	rt := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)
	if rt.Kind() != reflect.Ptr || rv.IsNil() {
		err = errors.New("please pass valid pointer")
		return
	}
	if rv.Elem().Kind() != reflect.Struct {
		err = errors.New("please pass struct")
		return
	}

	lineArr := strings.Split(string(data), "\n")

	var lastSectionName string
	for i, line := range lineArr {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if line[0] == ';' || line[0] == '#' {
			continue
		}

		if line[0] == '[' {
			if len(line) <= 2 || line[len(line)-1] != ']' {
				err = fmt.Errorf("syntax error, invalid section:%s, lineNo:%d\n", line, i+1)
				return
			}
		}

		if line[0] == '[' && line[len(line)-1] == ']' {
			lastSectionName, err = ParseSection(line, rt)
			if err != nil {
				return
			}
			continue
		}

		if strings.Contains(line, "=") {
			err = ParseLine(line, lastSectionName, rv, i)
			if err != nil {
				return
			}
		} else {
			err = fmt.Errorf("syntax error, invalid section:%s, lineNo:%d\n", line, i+1)
			return
		}

	}

	return
}

func ParseSection(line string, rt reflect.Type) (lastSectionName string, err error) {
	sectionName := strings.TrimSpace(line[1 : len(line)-1])
	if len(sectionName) == 0 {
		err = fmt.Errorf("syntax error, invalid section:%s", line)
		return
	}

	for i := 0; i < rt.Elem().NumField(); i++ {
		filed := rt.Elem().Field(i)
		tagVal := filed.Tag.Get("ini")
		if tagVal == sectionName {
			lastSectionName = filed.Name
			return
		}
	}

	return
}

func ParseLine(line, lastSectionName string, rv reflect.Value, i int) (err error) {
	if len(line) <= 3 {
		err = fmt.Errorf("syntax error, invalid,line:%d", i+1)
		return
	}

	index := strings.Index(line, "=")
	lineName := strings.TrimSpace(line[:index])
	lineVal := strings.TrimSpace(line[index+1:])

	if len(lineName) == 0 || len(lineVal) == 0 {
		err = fmt.Errorf("syntax error, invalid,line:%d", i+1)
		return
	}

	if len(lastSectionName) == 0 {
		err = fmt.Errorf("syntax error, invalid,line:%d, mast pass an section", i+1)
		return
	}

	sectionType := rv.Elem().FieldByName(lastSectionName).Type()
	if sectionType.Kind() != reflect.Struct {
		err = fmt.Errorf("section %s must be a struct ", sectionType.Name())
		return
	}
	for i := 0; i < sectionType.NumField(); i++ {
		filed := sectionType.Field(i)
		tagVal := filed.Tag.Get("ini")
		if tagVal == lineName {
			// TODO:添加其他类型
			switch sectionType.Field(i).Type.Kind() {
			case reflect.String:
				rv.Elem().FieldByName(lastSectionName).Field(i).SetString(lineVal)
			case reflect.Int:
				id, err := strconv.Atoi(lineVal)
				if err != nil {
					return err
				}
				rv.Elem().FieldByName(lastSectionName).Field(i).SetInt(int64(id))
			default:
				err = errors.New("Unsupport type")
				return
			}
			//rv.Elem().FieldByName(lastSectionName).Field(i).SetString(lineVal)
		}
	}

	return
}
