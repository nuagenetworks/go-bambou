package bambou

import "testing"

var r *TestRoot

func setup() {
	r = &TestRoot{
		ExposedObject: ExposedObject{
			Identity: TestRooIdentity,
		},
	}
}

func teardown() {
	r = nil
}

func TestSession_NewSession(t *testing.T) {
	setup()
	defer teardown()

	s := NewSession("username", "password", "organization", "url", r)

	if w := s.Username; "username" != w {
		t.Error("UserName should be 'username' but is %s", w)
	}

	if w := s.Password; "password" != w {
		t.Error("Password should be 'password' but is %s", w)
	}

	if w := s.Organization; "organization" != w {
		t.Error("Organization should be 'organization' but is %s", w)
	}

	if w := s.URL; "url" != w {
		t.Error("URL should be 'url' but is %s", w)
	}

	if s.Root == nil {
		t.Error("Root should be not be nil")
	}
}

func TestSession_MakeAuthorizationHeaders(t *testing.T) {
	setup()
	defer teardown()

	s := NewSession("username", "password", "organization", "url", r)

	if w := s.MakeAuthorizationHeaders(); "XREST dXNlcm5hbWU6cGFzc3dvcmQ=" != w {
		t.Error("Authorization Headers should be 'XREST dXNlcm5hbWU6cGFzc3dvcmQ=' but is %s", w)
	}
}

func TestSession_StartStopSession(t *testing.T) {
	setup()
	defer teardown()

	if CurrentSession() != nil {
		t.Error("CurrentSession should be nil")
	}

	s := NewSession("username", "password", "organization", "url", r)
	s.Start()

	if CurrentSession() != s {
		t.Error("CurrentSession should be equal to r")
	}

	if w := s.APIKey; "api-key" != w {
		t.Error("session APIKey should be 'api-key' but is %s", w)
	}

	if w := r.APIKey; "api-key" != w {
		t.Error("root object APIKey should be 'api-key' but is %s", r)
	}

	s.Reset()

	if CurrentSession() != nil {
		t.Error("CurrentSession should be nil")
	}

	if w := s.APIKey; "" != w {
		t.Error("session APIKey should be '' but is %s", w)
	}

	if w := r.APIKey; "" != w {
		t.Error("root object APIKey should be '' but is %s", w)
	}

}

var TestRooIdentity = Identity{
	RESTName:     "root",
	ResourceName: "root",
}

type TestRoot struct {
	ExposedObject

	UserName     string `json:"userName,omitempty"`
	Password     string `json:"password,omitempty"`
	APIKey       string `json:"APIKey,omitempty"`
	Organization string `json:"enterprise,omitempty"`
}

func (o *TestRoot) GetAPIKey() string                           { return o.APIKey }
func (o *TestRoot) SetAPIKey(key string)                        { o.APIKey = key }
func (o *TestRoot) GetURL() string                              { return CurrentSession().URL + "/" + o.Identity.ResourceName }
func (o *TestRoot) Save() *Error                                { return nil }
func (o *TestRoot) Delete() *Error                              { return nil }
func (o *TestRoot) GetURLForChildrenIdentity(i Identity) string { return "" }

func (o *TestRoot) Fetch() *Error {
	o.APIKey = "api-key"
	return nil
}
