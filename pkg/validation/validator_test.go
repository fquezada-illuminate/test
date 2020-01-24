package validation

import "testing"

func TestSingleton(t *testing.T) {
	v := Singleton()

	if _, ok := (*v).(Validator); !ok {
		t.Error("Singleton did not return a Validator type")
	}

	v2 := Singleton()

	if v != v2 {
		t.Error("Singleton did not return the same address of the validator")
	}
}
