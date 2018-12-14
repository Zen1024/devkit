package mapstruct

import (
	"testing"
)

type Person struct {
	Age     int       `json:"age"`
	Name    string    `json:"name"`
	Friends []*Person `json:"friends"`
	Gf      *Person   `json:"gf"`
	Extra   string    `json:"extra"`
	Extra2  string    `json:""`
	Extra3  string
}

func TestMapstruct(t *testing.T) {
	p := &Person{
		Age:    10,
		Name:   "gilf",
		Extra:  "extra",
		Extra2: "123",
		Gf: &Person{
			Age:  10,
			Name: "gilf",
		},
		Friends: []*Person{&Person{
			Age:  10,
			Name: "gilf",
			Friends: []*Person{&Person{
				Age:  10,
				Name: "gilf",
			}}},
			&Person{
				Age:  10,
				Name: "gilf",
				Friends: []*Person{&Person{
					Age:  10,
					Name: "gilf",
				}}}},
	}
	mapval := mapstruct(p, "json")
	t.Logf("got mapval:%+v", mapval)
	var p2 = new(Person)
	if err := scanstruct(mapval, "json", p2); err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("got scanned val:%+v,gf:%+v", p2, p2.Gf)

}
