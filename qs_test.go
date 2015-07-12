package qs

import (
	"fmt"
	"strings"
	"testing"
	"unsafe"
)

func TestNil(t *testing.T) {
	if r := Encode(nil); len(r) != 0 {
		t.Fatalf("qs has %d values", len(r))
	}
}

type Person struct {
	Age      int     `qs:"age"`
	Name     string  `qs:"name"`
	Location string  `qs:"-"`
	Parent   *Person `qs:"parent,omitempty"`
}

func (p Person) String() string { return p.Name }

func TestNonStructs(t *testing.T) {
	for _, v := range []interface{}{
		5,                 // primitive
		new(int),          // ptr
		make(map[int]int), // map
		[]int{1},          // slice
		[1]int{1},         // array
		func() {},         // func
		make(chan int),    // chan
		complex(1, 1),     // complex
	} {
		testShouldPanic(t, v, fmt.Sprintf("qs: passed value (type: %T) is not a struct", v))
	}
}

func TestIgnore(t *testing.T) {
	testOut(t, struct {
		Int    int    `qs:"-"`
		String string `qs:"-"`
	}{10, "foo"}, "")
}

func TestStructBasic(t *testing.T) {
	testOut(t, Person{
		Age:      10,
		Name:     "foo",
		Location: "home",
	}, "age=10&name=foo")
}

func TestPtrIndirection(t *testing.T) {
	i := 10
	n := "foo"
	ip := &i
	ipp := &ip

	testOut(t, struct {
		Age  **int
		Name *string
	}{ipp, &n}, "Age=10&Name=foo")
}

func TestStructOmitEmpty_All(t *testing.T) {
	testOut(t, struct {
		Int    int    `qs:"i,omitempty"`
		String string `qs:"s,omitempty"`
		Error  error  `qs:"err,omitempty"`
	}{}, "")
}

func TestStructStringer(t *testing.T) {
	testOut(t, Person{
		Name:   "child",
		Age:    5,
		Parent: &Person{Name: "dad", Age: 30},
	},
		"age=5&name=child&parent=dad")
}

func TestEmptyStructTag(t *testing.T) {
	testOut(t, struct {
		Int    int    `qs:""`
		String string `qs:""`
	}{10, "foo"},
		"Int=10&String=foo")
}

func TestDefaultNames(t *testing.T) {
	testOut(t, struct {
		Int    int `qs:",omitempty"`
		String string
	}{10, "foo"},
		"Int=10&String=foo")
}

func TestStructSlice(t *testing.T) {
	testOut(t, struct {
		Friends []string `qs:"f"`
	}{
		[]string{"foo", "bar", "qux"},
	},
		"f=foo&f=bar&f=qux")
}

func TestStructSlicePtr(t *testing.T) {
	testOut(t, struct {
		Friends *[]string `qs:"f"`
	}{
		&[]string{"foo", "bar", "qux"},
	},
		"f=foo&f=bar&f=qux")
}

func TestAllowedFields(t *testing.T) {
	// list all reflect.Kinds here

	structs := []struct {
		shouldPanic bool
		in          interface{}
		expected    string
	}{
		{false, struct{ X bool }{true}, "X=true"},                        // Bool
		{false, struct{ X uint64 }{uint64(10)}, "X=10"},                  // [U]Int[8|16|32|64]
		{false, struct{ X float64 }{10.123}, "X=10.123"},                 // Float[32|64]
		{true, struct{ X complex128 }{complex(10, 10)}, ""},              // Complex[64|128]
		{true, struct{ X chan int }{make(chan int)}, ""},                 // Chan
		{true, struct{ X func() }{func() {}}, ""},                        // Func
		{true, struct{ X map[int]int }{map[int]int{0: 0, 1: 1}}, ""},     // Map
		{true, struct{ X unsafe.Pointer }{unsafe.Pointer(new(int))}, ""}, // UnsafePointer
	}

	for _, c := range structs {
		if c.shouldPanic {
			testShouldPanic(t, c.in, "qs: type cannot be serialized")
		} else {
			testOut(t, c.in, c.expected)
		}
	}
}

func testOut(t *testing.T, v interface{}, expected string) {
	if out := Encode(v).Encode(); out != expected {
		t.Fatalf("wrong qs: '%s' expected: '%s'", out, expected)
	}
}

func testShouldPanic(t *testing.T, v interface{}, msg string) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Didn't panic")
		}
		if msg != "" && !strings.HasPrefix(fmt.Sprintf("%s", r), msg) {
			t.Fatalf("wrong panic text: %q\nexpected: %q", r, msg)
		}
	}()
	_ = Encode(v)
	t.Fatal("Did not panic")
}
