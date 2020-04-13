package app

import (
	"testing"

	"github.com/Ankr-network/dccn-tools/patch-pvc/app/pvc"
)

var app App

func TestInit(t *testing.T) {
	app = CreateApp("ns-90f4416b-584e-4a15-8fd8-39d3b442bda8", "app-b1ac239e-582f-4db3-9890-6b6b19e5e17f-lto-network-75fd97w7cn")
}
func TestGetPvcAndPv(t *testing.T) {
	err := app.GetPvcAndPv()
	if err != nil {
		t.Error(err)
	}
}

func TestDownLoadPvc(t *testing.T) {
	err := app.DownLoadPvc()
	if err != nil {
		t.Error(err)
	}
}

func TestHandPvc(t *testing.T) {
	pvc := pvc.Pvc1{}
	err := app.HandPvc(&pvc)
	if err != nil {
		t.Error(err)
	}
}

// func Test_ApplyPvc(t *testing.T) {}
