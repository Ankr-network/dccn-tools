简介
---

### 做什么的
```
累计总量images磁盘>100G 时，删除7天以上的镜像。
如果删除后仍然超过100G 则天数逐渐递减删除 直至小于100G
```

### 大致流程
```
通过 docker system df 获取images的大小

再通过 docker images --format "{{.ID}}\t{{.CreateAt}}" 获取所有的id 以及镜像的创建时间

最后通过 docker rmi IDS... 删除符合条件的镜像
```

### 使用说明
```
go run main.go -d 7 -s 100
// -d 天数  默认7天
// -s images 总和 GB  默认100GB
// image总和 大于100G时，删除超过7天的image, 如果还是超过100G 则天数递减，直至小于100G 或 小于0天
```

### 为什么 image 没有被删除掉
```
本程序没有做docker rm -f 强制删除操作

原因有二：
1、有些image 的容器还在使用中
2、有些image 做了tag ,image ID 重复，导致无法删除（本程序未考虑该情况，因为集群主要是pull image 不做tag）
、
```