package funs

import "testing"

func TestExec(t *testing.T) {
	_, err := Exec("ls")
	if err != nil {
		t.Error("exec err:", err)
	}
}
