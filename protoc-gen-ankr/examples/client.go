package main

import (
	"context"
	greeter "github.com/Ankr-network/dccn-tools/protoc-gen-ankr/examples/pb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:9090", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	cli := greeter.NewGreeterClient(conn)

	res, e := cli.World(context.Background(), &greeter.Request{Name: "alvin"})
	if e != nil {
		log.Fatal(e)
	}
	log.Println(res.Msg)

	res, e = cli.Hello(context.Background(), &greeter.Request{Name: "alvin"})
	if e != nil {
		log.Fatal(e)
	}
	log.Println(res.Msg)
}
