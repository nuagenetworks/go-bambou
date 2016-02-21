package bambou

import (
	"encoding/json"
	"testing"
)

func TestError_NewError(t *testing.T) {

	e := NewError(42, "a big error")

	if e.Code != 42 {
		t.Error("Code should be 42 but it %d", e.Code)
	}

	if e.Message != "a big error" {
		t.Error("Code should be 'a big error' but it '%s'", e.Message)
	}
}

func TestError_Error(t *testing.T) {

	e := NewError(42, "a big error")

	if w := e.Error(); w != "<Error: 42, message: a big error>" {
		t.Error("String() should be '<Error: 42, message: a big error>' but it '%s'", w)
	}
}

func TestError_FromJSON(t *testing.T) {

	e := NewError(42, "a big error")
	d := "{\"property\": \"prop\", \"type\": \"iznogood\", \"descriptions\": [{\"title\": \"oula\", \"description\": \"pas bon\"}]}"

	json.Unmarshal([]byte(d), e)

	if w := e.Message; w != "iznogood" {
		t.Error("Message should be 'iznogood' but it '%s'", w)
	}

	if w := e.Property; w != "prop" {
		t.Error("Property should be 'prop' but it '%s'", w)
	}

	if w := len(e.Descriptions); w != 1 {
		t.Error("Descriptions should be '1' but it '%d'", w)
	}

	if w := e.Descriptions[0].Title; w != "oula" {
		t.Error("Descriptions should be 'oula' but it '%s'", w)
	}

	if w := e.Descriptions[0].Description; w != "pas bon" {
		t.Error("Descriptions should be 'pas bon' but it '%s'", w)
	}

}
