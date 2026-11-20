package xinterp

import (
	"go/types"
	"reflect"
	"testing"
)

func TestNativeValueDoesNotExposeStorageCells(t *testing.T) {
	intSlice := &RuntimeSlice{
		ElemType: types.Typ[types.Int],
		Elements: []*Cell{
			{Value: 1},
			{Value: 2},
			{Value: 3},
		},
	}

	got := nativeValue(intSlice, reflect.TypeOf((*interface{})(nil)).Elem()).Interface()
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("nativeValue() = %#v (%T), want %#v (%T)", got, got, want, want)
	}
}

func TestNativeValueMaterializesVariadicSlice(t *testing.T) {
	intSlice := &RuntimeSlice{
		ElemType: types.Typ[types.Int],
		Elements: []*Cell{
			{Value: 1},
			{Value: 2},
		},
	}
	variadicArgs := &RuntimeSlice{
		ElemType: types.NewInterfaceType(nil, nil).Complete(),
		Elements: []*Cell{
			{Value: intSlice},
		},
	}

	target := reflect.TypeOf([]interface{}{})
	got := nativeValue(variadicArgs, target).Interface().([]interface{})
	if len(got) != 1 || !reflect.DeepEqual(got[0], []int{1, 2}) {
		t.Fatalf("nativeValue() = %#v, want []interface{}{[]int{1, 2}}", got)
	}
}

func TestNativeValuePreservesTypedNilSlice(t *testing.T) {
	nilSlice := &RuntimeSlice{ElemType: types.Typ[types.Int], Nil: true}
	got := nativeValue(nilSlice, reflect.TypeOf((*interface{})(nil)).Elem()).Interface()

	value, ok := got.([]int)
	if !ok || value != nil {
		t.Fatalf("nativeValue() = %#v (%T), want nil []int", got, got)
	}
}
