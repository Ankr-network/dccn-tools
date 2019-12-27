package ceph

import "testing"

func TestGetNodeInfo(t *testing.T) {
	for _, v := range GetNodeInfo("/home/mobius/.kube/config") {
		t.Logf("%+v\n", v)
	}
}

func TestGetUserName(t *testing.T) {
	exp := "mobius"
	act := GetUserName()
	if exp != act {
		t.Logf("act: %s exp: %s \n", act, exp)
	}
}
