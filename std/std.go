package std

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

func TypedSlice[T any](v any) []T {
	items, ok := v.([]any)
	if !ok {
		panic(fmt.Errorf("TypedSlice failed: argument is of type %s, not []any", reflect.TypeOf(v)))
	}
	result := make([]T, len(items))
	for i := range items {
		result[i] = items[i].(T)
	}
	return result
}

func Map[TIn any, TOut any](items []TIn, mapFn func(v TIn) TOut) []TOut {
	result := make([]TOut, len(items))
	for i := range items {
		result[i] = mapFn(items[i])
	}
	return result
}

func MapP[TIn any, TOut any](items []TIn, mapFn func(v *TIn) TOut) []TOut {
	result := make([]TOut, len(items))
	for i := range items {
		result[i] = mapFn(&items[i])
	}
	return result
}

func EnsureSuffix(str, suffix string) string {
	if strings.HasSuffix(str, suffix) {
		return str
	}
	return str + suffix
}

// Fdump writes the structure's fields to the given writer
func Fdump(w io.Writer, v any) {
	printFieldsRecursive(w, reflect.ValueOf(v), 0, true)
}

// Dump writes the structure's fields to stdout
func Dump(v any) {
	Fdump(os.Stdout, v)
}

// Recursive function to print all public fields of a struct
func printFieldsRecursive(w io.Writer, val reflect.Value, indentLevel int, shouldIndent bool) {
	isPtr := val.Kind() == reflect.Ptr
	if isPtr {
		val = val.Elem()
	}
	indent := strings.Repeat("  ", indentLevel)
	// Handle different kinds
	switch val.Kind() {
	case reflect.Struct:
		fmt.Fprintf(w, "%s\n", niceType(val.Type(), isPtr))
		typ := val.Type()
		for i := 0; i < val.NumField(); i++ {
			fieldVal := val.Field(i)
			fieldType := typ.Field(i)

			// Only print public fields (those starting with uppercase letters)
			if fieldType.PkgPath == "" {
				fmt.Fprintf(w, "%s%s: ", indent, fieldType.Name)
				printFieldsRecursive(w, fieldVal, indentLevel+1, false)
			}
		}
	case reflect.Slice, reflect.Array:
		fmt.Fprintf(w, "[\n")
		for i := 0; i < val.Len(); i++ {
			printFieldsRecursive(w, val.Index(i), indentLevel+1, true)
		}
		fmt.Fprintf(w, "%s]\n", indent)
	case reflect.Map:
		fmt.Fprintf(w, "%s{\n", niceType(val.Type(), isPtr))
		for _, key := range val.MapKeys() {
			indent := strings.Repeat("  ", indentLevel+1)
			fmt.Fprintf(w, "%s%s: ", indent, key)
			printFieldsRecursive(w, val.MapIndex(key), indentLevel+1, shouldIndent)
		}
		fmt.Fprintf(w, "%s}\n", indent)
	default:
		// Base case: print the value of leaf nodes
		indent := indent
		if !shouldIndent {
			indent = ""
		}
		fmt.Fprintf(w, "%s%s %v\n", indent, niceType(val.Type(), isPtr), val.Interface())
	}
}

func niceType(typ reflect.Type, isPtr bool) string {
	p := ""
	if isPtr {
		p = "*"
	}
	return "｢" + p + strings.ReplaceAll(typ.String(), "interface {}", "any") + "｣"
}
