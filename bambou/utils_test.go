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

/*
   Fake Exposed
*/

var fakeIdentity = Identity{"fake", "fakes"}

type fakeObjectsList []*fakeObject
type fakeObject struct{ ExposedObject }

func (o *fakeObject) Save() *Error   { return nil }
func (o *fakeObject) Delete() *Error { return nil }
func (o *fakeObject) Fetch() *Error  { return nil }

/*
   Fake Rootable
*/
var fakeRootdentity = Identity{"root", "root"}

type fakeRootObject struct{ fakeObject }

func (o *fakeRootObject) GetAPIKey() string    { return "api-key" }
func (o *fakeRootObject) SetAPIKey(key string) {}

/*
   Fake Unfetchable Rootable
*/

type fakeUnfetchableRootObject struct{ fakeRootObject }

func (o *fakeUnfetchableRootObject) Fetch() *Error { return NewError(500, "Error") }
