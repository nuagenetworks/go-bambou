//+build !test

package bambou

/*
   Fake Exposed
*/

var fakeIdentity = Identity{"fake", "fakes"}

type fakeObjectsList []*fakeObject
type fakeObject struct {
	RemoteObject
	Name string `json:"name"`
}

/*
   Fake Rootable
*/
var fakeRootIdentity = Identity{"root", "root"}

type fakeRootObject struct {
	fakeObject
	Token string `json:"APIKey,omitempty"`
}

func (o *fakeRootObject) APIKey() string       { return o.Token }
func (o *fakeRootObject) SetAPIKey(key string) { o.Token = key }
