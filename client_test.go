package mailofly

import "testing"

func TestNewRequiresAPIKey(t *testing.T) {
	_, err := New(Options{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestErrorMessage(t *testing.T) {
	e := &Error{Status: 401, Err: "unauthorized", DetailMessage: "Invalid key"}
	if e.Error() != "unauthorized: Invalid key" {
		t.Fatalf("got %q", e.Error())
	}
}
