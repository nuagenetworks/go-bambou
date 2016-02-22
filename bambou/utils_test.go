//+build !test

package bambou

import "reflect"

type restorer func()

func (r restorer) restore() {
	r()
}

func patch(dest, value interface{}) restorer {
	destv := reflect.ValueOf(dest).Elem()
	oldv := reflect.New(destv.Type()).Elem()
	oldv.Set(destv)
	valuev := reflect.ValueOf(value)
	if !valuev.IsValid() {
		valuev = reflect.Zero(destv.Type())
	}
	destv.Set(valuev)
	return func() {
		destv.Set(oldv)
	}
}
