package image

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/Ankr-network/dccn-tools/rm-docker-image/funs"
)

// DelImage 删除镜像的条件
type DelImage struct {
	date int     // 删除大于多少天的
	size float64 // 超过 size 继续删除
}

// Run 执行删除 条件达到的镜像
func (i *DelImage) Run() {
	for {
		if i.date < 0 {
			fmt.Println("DelImage.date 已小于0天，无可再删了的")
			return
		}
		if i.getSize() > i.size {
			i.delete()
			i.date--
		} else {
			return
		}
		time.Sleep(time.Second)
	}
}

func (i *DelImage) getSize() float64 {
	str, err := funs.Exec("docker", "system", "df")
	if err != nil {
		fmt.Printf("docker system df 获取出错：%s\n", err)
		return 0
	}
	bf := bufio.NewReader(strings.NewReader(str))
	bf.ReadString('\n')
	imageStr, err := bf.ReadString('\n')

	if err != nil {
		fmt.Printf("readString 没有获取数据，err:%s\n", err)
		return 0
	}
	reclaimable := strings.Trim(imageStr, "\r\n ")
	reclaimables := strings.Fields(reclaimable)
	b := strings.Contains(reclaimables[3], "B")
	var (
		sizenum float64
		numtype string
	)
	if b {
		str := reclaimables[3]
		strbyte := str[:len(reclaimables[3])-2]
		sizenum, _ = strconv.ParseFloat(string(strbyte), 64)
		numtype = string(str[len(reclaimables[3])-2:])
	} else {
		sizenum, _ = strconv.ParseFloat(reclaimables[3], 64)
		numtype = reclaimables[4]
	}

	switch numtype {
	case "MB":
		return sizenum / 1024
	case "GB":
		return sizenum
	case "KB":
		return sizenum / 1024 / 1024
	case "TB":
		return sizenum * 1024
	default:
		fmt.Printf("未识别到images所占的内存总量：%s   %s\n", reclaimables[3], reclaimables[4])
	}
	return 0
}

func (i *DelImage) delete() {
	images := make(images, 0)
	images.getCreated()
	var ids []string
	ids = append(ids, "rmi")
	for _, v := range images {
		if v.created >= i.date {
			ids = append(ids, v.ID)
		}
	}
	if len(ids) == 1 {
		fmt.Printf("delete %d day: 无images可删 \n", i.date)
		return
	}

	funs.Exec("docker", ids...)

	fmt.Printf("delete %d day: docker %v \n", i.date, strings.Join(ids, " "))
}

// CreateDelImage 创建镜像
func CreateDelImage(date int, size float64) DelImage {
	return DelImage{date, size}
}

type image struct {
	ID      string
	created int // 创建镜像的天数
}

// image 应该有 delete getSize 等功能呢

type images []image

// images 应该有获取总size,总删除，筛选的功能
func (i *images) getCreated() {
	data, err := funs.Exec("docker", "images", "--format", `{{.ID}}\t{{.CreatedAt}}`)
	if err != nil {
		fmt.Printf("docker images --format 获取数据出错：%s\n", err)
		return
	}
	bf := bufio.NewReader(strings.NewReader(data))
	now := time.Now()
	timeLayout := "2006-01-02 15:04:05" //转化所需模板
	for {
		str, err := bf.ReadString('\n')

		if err != nil {
			if err == io.EOF { //文件已经结束
				break
			}
			fmt.Println("getCreated readeSteing err = ", err)
		}
		strs := strings.Fields(str)
		tmp, _ := time.Parse(timeLayout, strs[1]+" "+strs[2])
		sub := now.Sub(tmp)
		*i = append(*i, image{strs[0], int(sub.Hours() / 24)})
	}
}
