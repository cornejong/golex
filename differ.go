package golex

import (
	"fmt"
	"reflect"
	"strings"
)

// Diff stores the differences between two values
type Diff struct {
	Field  string
	Expect interface{}
	Got    interface{}
}

func (d Diff) String() string {
	return fmt.Sprintf("Field: %s\nExpected: %v\nGot: %v\n", d.Field, d.Expect, d.Got)
}

// Differ is the main diff engine
type Differ struct {
	Diffs []Diff
}

// Compare is the entry point for comparing two values
func (d *Differ) Compare(expected, got any) {
	d.compare(reflect.ValueOf(expected), reflect.ValueOf(got), "")
}

// compare recursively compares two reflect.Values
func (d *Differ) compare(expected, got reflect.Value, fieldPath string) {
	if !expected.IsValid() && !got.IsValid() {
		return // Both are invalid, no diff
	}

	if expected.IsValid() && !got.IsValid() {
		d.addDiff(fieldPath, expected.Interface(), nil)
		return
	}

	if !expected.IsValid() && got.IsValid() {
		d.addDiff(fieldPath, nil, got.Interface())
		return
	}

	if expected.Type() != got.Type() {
		d.addDiff(fieldPath, expected.Interface(), got.Interface())
		return
	}

	switch expected.Kind() {
	case reflect.Struct:
		d.compareStructs(expected, got, fieldPath)
	case reflect.Map:
		d.compareMaps(expected, got, fieldPath)
	case reflect.Slice, reflect.Array:
		d.compareSlices(expected, got, fieldPath)
	case reflect.Interface, reflect.Ptr:
		d.compare(expected.Elem(), got.Elem(), fieldPath)
	default:
		if !reflect.DeepEqual(expected.Interface(), got.Interface()) {
			d.addDiff(fieldPath, expected.Interface(), got.Interface())
		}
	}
}

func (d *Differ) compareStructs(expected, got reflect.Value, fieldPath string) {
	for i := 0; i < expected.NumField(); i++ {
		fieldName := expected.Type().Field(i).Name
		d.compare(expected.Field(i), got.Field(i), fieldPath+"."+fieldName)
	}
}

func (d *Differ) compareMaps(expected, got reflect.Value, fieldPath string) {
	for _, key := range expected.MapKeys() {
		expVal := expected.MapIndex(key)
		gotVal := got.MapIndex(key)
		d.compare(expVal, gotVal, fmt.Sprintf("%s[%v]", fieldPath, key))
	}
	for _, key := range got.MapKeys() {
		if expected.MapIndex(key).IsValid() {
			continue // already handled
		}
		d.addDiff(fmt.Sprintf("%s[%v]", fieldPath, key), nil, got.MapIndex(key).Interface())
	}
}

func (d *Differ) compareSlices(expected, got reflect.Value, fieldPath string) {
	for i := 0; i < expected.Len(); i++ {
		if i < got.Len() {
			d.compare(expected.Index(i), got.Index(i), fmt.Sprintf("%s[%d]", fieldPath, i))
		} else {
			d.addDiff(fmt.Sprintf("%s[%d]", fieldPath, i), expected.Index(i).Interface(), nil)
		}
	}

	for i := expected.Len(); i < got.Len(); i++ {
		d.addDiff(fmt.Sprintf("%s[%d]", fieldPath, i), nil, got.Index(i).Interface())
	}
}

func (d *Differ) addDiff(field string, expect, got any) {
	d.Diffs = append(d.Diffs, Diff{
		Field:  strings.TrimPrefix(field, "."),
		Expect: expect,
		Got:    got,
	})
}

// String returns the string representation of all diffs
func (d *Differ) String() string {
	var sb strings.Builder
	for _, diff := range d.Diffs {
		sb.WriteString(diff.String())
		sb.WriteString("\n")
	}
	return sb.String()
}

func (d *Differ) HasDifference() bool {
	return len(d.Diffs) > 0
}
