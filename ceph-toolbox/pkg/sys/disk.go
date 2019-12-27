package sys

import (
	"bufio"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const splitChar = "##"

func LookUpValidDisk(diskInfo string) (string, int) {
	scan := bufio.NewScanner(strings.NewReader(diskInfo))
	// ignore head line
	scan.Scan()
	scan.Text()
	var (
		ds  Disks
		err error
	)
	for scan.Scan() {
		re := regexp.MustCompile(`[ ]+`)
		fs := strings.Split(re.ReplaceAllString(scan.Text(), splitChar), splitChar)
		d := DiskInfo{}
		d.size, err = strconv.Atoi(fs[3])
		if err != nil {
			continue
		}
		d.path = fs[5]
		ds = append(ds, d)
	}
	if len(ds) > 0 {
		sort.Sort(ds)
		return ds[0].path, ds[0].size
	}
	return "", 0
}

type DiskInfo struct {
	path string
	size int
}

type Disks []DiskInfo

func (d Disks) Len() int {
	return len(d)
}

func (d Disks) Less(i, j int) bool {
	if d[i].size > d[j].size {
		return true
	} else if len(d[i].path) < len(d[j].path) {
		return true
	}
	return false
}

func (d Disks) Swap(i, j int) {
	t := d[i]
	d[i] = d[j]
	d[j] = t
}
