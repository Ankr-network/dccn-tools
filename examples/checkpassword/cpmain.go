package main

import (
	checkPassword "checkpassword/proto"
	"context"
	"fmt"
	"net"

	public "github.com/Ankr-network/dccn-tools/examples/public"
	"google.golang.org/grpc"
)

type holder struct{}

func (h *holder) Check(ctx context.Context, req *checkPassword.Request) (*checkPassword.Response, error) {
	rsp := &checkPassword.Response{
		Rsp: &public.RspCodeMsg{
			Code: "0000",
			Msg:  "do it successfully",
		},
	}
	return rsp, nil
}

func main() {
	l, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		fmt.Println(err)
		return
	}
	gs := grpc.NewServer()
	checkPassword.RegisterCheckPWDServer(gs, &holder{})
	if err := gs.Serve(l); err != nil {
		fmt.Println("service start failed: ", err)
		return
	}
}
