package build

import (
	"github.com/spf13/cast"
	"reflect"
	"strings"
	"github.com/spf13/viper"
	"fmt"
)

type any interface{}
func dumpType(a any) {
	fmt.Println(reflect.TypeOf(a))
}

type Value struct {
	values interface{}
}

func (v *Value) To(s string) *Value {
	parts := strings.Split(s, ".")
	u := v
	for _,p := range parts {
		u = u.Get(p)
	}
	return u
}

func (v *Value) Get(s string) *Value {
	if v.IsMap() {
		h := reflect.ValueOf(s)
		a := reflect.ValueOf(v.values)
		actual := a.MapIndex(h)
		if actual.IsValid() {
			return NewValue(actual.Interface())
		} else {
			return NewValue(nil)
		}
	} else if v.IsViper() {
		if vip,ok := v.values.(*viper.Viper); ok {
			return NewValue(vip.Get(s))
		} else {
			return NewValue(nil)
		}
	} else {
		return NewValue(nil)
	}
}

func (v *Value) IsViper() bool {
	_,ok := v.values.(*viper.Viper)
	return ok
}

func (v *Value) HasValue() bool {
	return v.values != nil
}

func (v *Value) ToInt() int {
	return cast.ToInt(v.values)
}

func (v *Value) String(name string) string {
	return cast.ToString(v.Get(name).values)
}

func (v *Value) ToString() string {
	return cast.ToString(v.values)
}

func (v *Value) ToSlice() []interface{} {
	return cast.ToSlice(v.values)
}

func (v *Value) IsMap() bool {
	t := reflect.TypeOf(v.values)
	return t != nil && t.Kind() == reflect.Map
}

func (v *Value) In(i int) *Value {
	ar := v.ToSlice()
	n := len(ar)
	if 0 <= i && i < n {
		return NewValue(ar[i])
	} else {
		return NewValue(nil)
	}
}

func NewValue(m interface{}) *Value {
	return &Value{
		values: m,
	}
}


