package pointerhelp

import "testing"

func TestBoolPtr(t *testing.T) {
	myBool := true
	myBoolPtr := BoolPtr(myBool)
	if *myBoolPtr != true {
		t.Error("Wait a minute")
	}
}

func TestSafeBool(t *testing.T) {
	var myBoolPtr *bool = nil
	mySafeBool := SafeBool(myBoolPtr)
	if mySafeBool != false {
		t.Error("Wait a minute...")
	}
}
