package build

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"reflect"
	"path/filepath"
	ts "pregame/thelpers"
	"fmt"
)


func createMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["here"] = 1
	return m
}

func readBook(file string) Tasks {
	yaml := filepath.Join(ts.Pwd(), file)
	book := ReadYaml(yaml).Load().YamlPlays()
	return book
}

func TestValue(t *testing.T) {

	Convey("steps method", t, func() {
		tasks := readBook(".files/1st.yaml")
		steps := tasks[0].Series
		fmt.Println(tasks[0].Name)
		So(steps, ShouldNotBeNil)
		So(len(steps), ShouldEqual, 2)
	})

	Convey("book steps has specific type", t, func() {
		tasks := readBook(".files/1st.yaml")
		steps := tasks[0].Parallel
		So(steps, ShouldNotBeNil)
		So(len(steps), ShouldEqual, 0)
	})

	Convey("walking path produces value", t, func() {
		a := make(map[string]interface{})
		a["stuff"] = "where"
		b := make(map[string]interface{})
		b["here"] = a

		u := NewValue(b)
		u = u.To("here.stuff")

		So(u.ToString(), ShouldEqual, "where")
	})

	Convey("should get the 'here' value from the map", t, func() {
		m := make(map[string]string)
		m["here"] = "there"
		t := NewValue(m)
		t = t.Get("here")
		So(t.HasValue(), ShouldBeTrue)
		So(t.IsMap(), ShouldBeFalse)
		So(t.ToString(), ShouldEqual, "there")
	})

	Convey("using a path to a value that EXISTS should produce that value", t, func() {
		m := createMap()
		t := NewValue(m)
		t = t.Get("here")
		So(t.HasValue(), ShouldBeTrue)
		So(t.IsMap(), ShouldBeFalse)
		So(t.ToInt(), ShouldEqual, 1)
	})

	Convey("converting with just key", t, func() {
		m := make(map[string]string)
		m["here"] = "stuff"
		h := reflect.ValueOf("here")
		actual := reflect.ValueOf(m).MapIndex(h)
		val := actual.Interface()
		So(val, ShouldEqual, "stuff")
	})

	Convey("should be true that Value with 1 has a Value and that value is not a map", t, func() {
		t := NewValue(1)
		So(t.HasValue(), ShouldBeTrue)
		So(t.IsMap(), ShouldBeFalse)
	})

	Convey("should be true that Value with 1 has a Value and that value is not a map", t, func() {
		types := []interface{}{
			make(map[string]string),
			make(map[string]int),
			make(map[string]interface{}),
			make(map[interface{}]interface{}),
		}
		for _,k := range types {
			t := NewValue(k)
			So(t.HasValue(), ShouldBeTrue)
			isMap := t.IsMap()
			if !isMap {
				dumpType(k)
			}
			So(isMap, ShouldBeTrue)
		}
	})

	Convey("using a path to a value that DOESN'T exist should produce nil value", t, func() {
		m := createMap()
		t := NewValue(m)
		t = t.Get("nope")
		So(t.HasValue(), ShouldBeFalse)
		So(t.IsMap(), ShouldBeFalse)
	})

	Convey("a map backed Value should return true for isMap", t, func() {
		m := createMap()
		t := NewValue(m)
		So(t.IsMap(), ShouldBeTrue)
	})

	Convey("a map backed Value should have a 'value'", t, func() {
		m := createMap()
		t := NewValue(m)
		So(t.values, ShouldNotBeNil)
		So(t.HasValue(), ShouldBeTrue)
	})

	Convey("new value should not have backing value(s)", t, func() {
		t := &Value{}
		So(t.values, ShouldBeNil)
		So(t.HasValue(), ShouldBeFalse)
	})
}



