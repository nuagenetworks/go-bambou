package bambou

import "testing"

func TestRequest_NewRequest(t *testing.T) {

	r := NewRequest("https://fake.com")

	if r.URL != "https://fake.com" {
		t.Errorf("Data should be 'https://fake.com' but is '%s'", r.URL)
	}

	if r.Method != RequestMethodGet {
		t.Errorf("Method should be '%s' but is '%s'", RequestMethodGet, r.Method)
	}

	if r.Headers == nil {
		t.Errorf("Headers should not be nil")
	}

	if r.Parameters == nil {
		t.Errorf("Parameters should not be nil")
	}
}

func TestRequest_SetGetHeader(t *testing.T) {

	r := NewRequest("https://fake.com")
	r.SetHeader("header", "value")

	v := r.GetHeader("header")
	if v != "value" {
		t.Errorf("GetHeader(\"value\") should be 'value' but is '%s'", v)
	}
}

func TestRequest_SetGetParameter(t *testing.T) {

	r := NewRequest("https://fake.com")
	r.SetParameter("param", "value")

	v := r.GetParameter("param")
	if v != "value" {
		t.Errorf("GetParameter(\"value\") should be 'value' but is '%s'", v)
	}
}

func TestReques_ToNative(t *testing.T) {

	r := NewRequest("https://fake.com")
	r.SetHeader("header", "value")
	r.SetParameter("param", "value")
	r.Data = []byte("hello")

	n := r.ToNative()

	if u := n.URL.String(); r.URL != u {
		t.Errorf("native request URL should be \"https://fake.com\" but is '%s'", u)
	}

	if h := r.Headers["header"]; n.Header.Get("header") != h {
		t.Errorf("native request Header for 'header' should be \"value\" but is '%s'", h)
	}

}
