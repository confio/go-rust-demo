package api

import (
	"testing"
)

func TestAdd(t *testing.T) {
	res := Add(5, 7)
	if res != 12 {
		t.Fatalf("Unexpected result: %d", res)
	}
}

func TestGreet(t *testing.T) {
	res := string(Greet([]byte("Fred")))
	if res != "Hello, Fred" {
		t.Fatalf("Bad greet: %s", res)
	}

	res = string(Greet(nil))
	if res != "Hello, <nil>" {
		t.Fatalf("Bad greet: %s", res)
	}
}

func TestDivide(t *testing.T) {
	res, err := Divide(15, 3)
	if err != nil {
		t.Fatalf("Expected no error, got %s", err)
	}
	if res != 5 {
		t.Fatalf("Unexpected result: %d", res)
	}

	res, err = Divide(6, 0)
	if err == nil {
		t.Fatalf("Expected error, but got none")
	}
	errMsg := err.Error()
	if errMsg != "Cannot divide by zero" {
		t.Fatalf("Unexpected error msg: %s", errMsg)
	}
	if res != 0 {
		t.Fatalf("Unexpected result: %d", res)
	}
}

func TestRandomMessage(t *testing.T) {
	// this should pass
	res, err := RandomMessage(123)
	if err != nil {
		t.Fatalf("Expected no err, got %s", err)
	}
	if res != "You are a winner!" {
		t.Fatalf("Unexpected result: %s", res)
	}

	// this should error (normal)
	res, err = RandomMessage(-20)
	if err == nil {
		t.Fatalf("Expected error, but got none")
	}
	if errMsg := err.Error(); errMsg != "Too low" {
		t.Fatalf("Unexpected error msg: %s", errMsg)
	}
	if res != "" {
		t.Fatalf("Unexpected result: %s", res)
	}

	// this should panic
	res, err = RandomMessage(0)
	if err == nil {
		t.Fatalf("Expected error, but got none")
	}
	if errMsg := err.Error(); errMsg != "Caught panic" {
		t.Fatalf("Unexpected error msg: %s", errMsg)
	}
	if res != "" {
		t.Fatalf("Unexpected result: %s", res)
	}

	// this should pass (again)
	res, err = RandomMessage(789)
	if err != nil {
		t.Fatalf("Expected no err, got %s", err)
	}
	if res != "You are a winner!" {
		t.Fatalf("Unexpected result: %s", res)
	}
}

/** test helpers **/

type Lookup struct {
	data map[string]string
}

func NewLookup() *Lookup {
	return &Lookup{data: make(map[string]string)}
}

func (l *Lookup) Get(key []byte) []byte {
	val := l.data[string(key)]
	return []byte(val)
}

func (l *Lookup) Set(key, value []byte) {
	l.data[string(key)] = string(value)
}

func TestDemoDBAccess(t *testing.T) {
	l := NewLookup()
	l.Set([]byte("foo"), []byte("long text that fills the buffer"))
	l.Set([]byte("bar"), []byte("short"))

	// long
	res := DemoDBAccess(l, []byte("foo"))
	if string(res) != "Got value: long text that fills the buffer" {
		t.Errorf("Unexpected result (long): %s", string(res))
	}

	// short
	res = DemoDBAccess(l, []byte("bar"))
	if string(res) != "Got value: short" {
		t.Errorf("Unexpected result (short): %s", string(res))
	}

	// missing
	res = DemoDBAccess(l, []byte("missing"))
	if string(res) != "Got value: <nil>" {
		t.Errorf("Unexpected result (missing): %s", string(res))
	}
}
