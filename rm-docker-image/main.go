package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/Ankr-network/dccn-tools/rm-docker-image/image"
)

var (
	h    bool
	date int
	size float64
)

func init() {
	flag.BoolVar(&h, "h", false, "帮助文档")
	flag.IntVar(&date, "d", 7, "大于多少天的将会被删除")
	flag.Float64Var(&size, "s", 100, "大于多少GB的将会被删除")
	flag.Parse()
}

func main() {
	if h {
		flag.Usage()
		return
	}
	if date < 0 || size < 0 {
		fmt.Println("请填写正确的信息，date >= 0 && size >= 0")
		return
	}
	fmt.Printf("------- %s start -------\n", time.Now().Format("2006-01-02"))
	defer fmt.Printf("------- %s end -------\n\n", time.Now().Format("2006-01-02"))

	dic := image.CreateDelImage(date, size)
	dic.Run()
}
