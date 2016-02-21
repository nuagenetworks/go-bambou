package bambou

import (
	"encoding/json"
	"testing"
)

func TestNotification_NewNotification(t *testing.T) {

	n := NewNotification()

	if n.Events == nil {
		t.Error("Events should not be nil")
	}

	if n.UUID != "" {
		t.Error("UUID should be '' but it '%s'", n.UUID)
	}
}

func TestNotification_FromJSON(t *testing.T) {

	n := NewNotification()
	d := "{\"uuid\": \"007\", \"events\": [{\"entityType\": \"cat\", \"type\": \"UPDATE\", \"updateMechanism\": \"useless\", \"entities\":[{\"name\": \"hello\"}]}]}"

	json.Unmarshal([]byte(d), &n)

	if w := n.UUID; w != "007" {
		t.Error("UUID should be '007' but it '%s'", w)
	}

	if w := len(n.Events); w != 1 {
		t.Error("Len of Events should be '1' but it '%d'", w)
	}

	if w := n.Events[0].EntityType; w != "cat" {
		t.Error("EntityType should be 'cat' but it '%s'", w)
	}

	if w := n.Events[0].Type; w != "UPDATE" {
		t.Error("Type should be 'UPDATE' but it '%s'", w)
	}

	if w := n.Events[0].UpdateMechanism; w != "useless" {
		t.Error("UpdateMechanism should be 'useless' but it '%s'", w)
	}

	if w := len(n.Events[0].DataMap); w != 1 {
		t.Error("Len of DataMap should be '1' but it '%d'", w)
	}

	if w := n.Events[0].DataMap[0]["name"]; w != "hello" {
		t.Error("name should be 'hello' but it '%s'", w)
	}
}
