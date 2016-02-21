package bambou

import "testing"

func TestResponse_NewResponse(t *testing.T) {

	r := NewResponse()

	if r.Headers == nil {
		t.Errorf("Headers should not be nil")
	}
}

func TestResponse_SetGetHeader(t *testing.T) {

	r := NewResponse()
	r.SetHeader("header", "value")

	v := r.GetHeader("header")
	if v != "value" {
		t.Errorf("GetHeader(\"value\") should be 'value' but is '%s'", v)
	}
}
