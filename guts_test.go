package guts

import (
	"encoding/json"
	"testing"
)

func TestBase(t *testing.T) {
	a := 3
	g := Gut(a)
	gp := Gut(&a)

	if g.Id != 0 {
		t.Fatal("value should not have address")
	}

	if g.Kind != "int" {
		t.Fatal("int should have kind int")
	}

	if g.Type != "int" {
		t.Fatal("int should have type int")
	}

	if g.Data != a {
		t.Fatal("data should have value of input")
	}

	if gp.Id == 0 {
		t.Fatal("Pointer should have address")
	}

	if gp.Kind != "int" {
		t.Fatal("int pointer should have kind int")
	}

	if gp.Type != "*int" {
		t.Fatal("int pointer should have type *int")
	}

	if gp.Data != &a {
		t.Fatal("data should have value of input value")
	}

}

func TestSlice(t *testing.T) {
	a := 3
	b := []*int{&a, &a, nil}
	g := Gut(b)

	s, ok := g.Data.([]Guts)
	if !ok {
		t.Fatal("Result of gutting a slice should be []Guts")
	}

	for _, e := range s {
		if e.Type != "*int" {
			t.Fatal("all elements in slice should be of type *int")
		}
	}

	if s[0].Id == 0 {
		t.Fatal("int pointer should have an address")
	}

	if s[1].Id == 0 {
		t.Fatal("int pointer should have an address")
	}

	if s[2].Id != 0 {
		t.Fatal("nil pointer should not have an address")
	}

	if s[0].Kind != "int" {
		t.Fatal("int pointer should be of kind int")
	}

	if s[1].Kind != "int" {
		t.Fatal("int pointer should be of kind int")
	}

	if s[2].Kind != "invalid" {
		t.Fatal("nil pointer should be of kind invalid")
	}

	if _, err := json.Marshal(g); err != nil {
		t.Fatal("slice should always marshal to json")
	}

}

func TestMap(t *testing.T) {
	v := map[int]string{
		3: "Yes",
		5: "no",
	}

	g := Gut(v)

	if g.Id == 0 {
		t.Fatal("Map is a reference type")
	}
	if g.Kind != "map" {
		t.Fatal("kind of map should be map")
	}

	if _, err := json.Marshal(g); err != nil {
		t.Fatal("map should marshal to json")
	}

	m, ok := g.Data.(map[string]Guts)
	if !ok {
		t.Fatal("Result of gutting a map should be map[string]Guts")
	}

	a, ok := m["3"]
	if !ok {
		t.Fatal("Guts map should have string keys based on sprint")
	}

	if a.Data != v[3] {
		t.Fatal("Guts nested data does not match input")
	}

}

type T struct {
	A, B int
	C    *int
	D    string
	e    string
}

func TestStruct(t *testing.T) {
	a := 3
	v := T{1, 3, &a, "Hello", "Hidden"}
	g := Gut(&v)

	if _, err := json.Marshal(g); err != nil {
		t.Fatal("struct should always marshall to json")
	}

	m, ok := g.Data.(map[string]Guts)
	if !ok {
		t.Fatal("Result of gutting a struct should be map[string]Guts")
	}

	if _, ok := m["A"]; !ok {
		t.Fatal("Guts map should have key string for struct fields")
	}

	if _, ok := m["e"]; ok {
		t.Fatal("Unexported fields should not be present in result")
	}

}
