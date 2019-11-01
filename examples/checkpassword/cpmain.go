package main

import (
	"context"
	"fmt"
	"net"

	checkPassword "github.com/Ankr-network/dccn-tools/examples/checkpassword/proto"
	"google.golang.org/grpc"
)

type holder struct{}

func (h *holder) Check(ctx context.Context, req *checkPassword.Request) (*checkPassword.Response, error) {
	rsp := &checkPassword.Response{
		Rsp: &checkPassword.RspCodeMsg{
			Code: "0000",
			Msg:  "do it successfully",
		},
	}
	return rsp, nil
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:8081")
	if err != nil {
		fmt.Println(err)
		return
	}
	gs := grpc.NewServer()
	checkPassword.RegisterCheckPWDServer(gs, &holder{})
	fmt.Println("start check password service ...")
	if err := gs.Serve(l); err != nil {
		fmt.Println("service start failed: ", err)
		return
	}
}
