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

type fakeExposedList []*fakeExposed

type fakeExposed struct {
	ExposedObject
}

func (o *fakeExposed) Save() *Error {
	return nil
}
func (o *fakeExposed) Delete() *Error {
	return nil
}
func (o *fakeExposed) Fetch() *Error {
	return nil
}

/*
   Fake Rootable
*/
var testRootdentity = Identity{
	RESTName:     "root",
	ResourceName: "root",
}

type testRoot struct {
	fakeExposed

	UserName     string `json:"userName,omitempty"`
	Password     string `json:"password,omitempty"`
	APIKey       string `json:"APIKey,omitempty"`
	Organization string `json:"enterprise,omitempty"`
}

func (o *testRoot) GetAPIKey() string {
	return o.APIKey
}
func (o *testRoot) SetAPIKey(key string) {
	o.APIKey = key
}
func (o *testRoot) Fetch() *Error {
	o.APIKey = "api-key"
	return nil
}

type testFailedRoot struct {
	testRoot
}

func (o *testFailedRoot) Fetch() *Error { return NewError(500, "Error") }
