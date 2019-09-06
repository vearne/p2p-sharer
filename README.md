# p2p-sharer

标签（空格分隔）： p2p download

---

### 1.前言
这是一个非常简单的p2p文件分发系统
它存在的目的，主要是验证p2p下载的主要思想，所以并不是完全按照BT的公开协议来实现的，目前已经经过了初步的测试，但并没有在生产环境，进行大规模的使用，因此请谨慎的使用。



### 2.组件
#### 2.1 制作种子文件
```
./p2p-sharer gen --tracker 192.168.10.200:35330 --filePath /tmp/Motrix-1.4.1.dmg --seedPath ./
```
`tracker` tracker 地址
`filePath` 待分发的文件路径
`seedPath` 生成的种子文件的存放路径

#### 2.2 tracker
1. 保存(供其他node查询)文件分片在整个集群的分布情况
2. 维护节点的上下线情况(通过心跳)

#### 2.3 node
1. 接收并执行下载任务
2. 与其它node共享文件分片

### 3. 使用
```
git clone git@github.com:vearne/p2p-sharer.git
```
#### 3.1 build
##### 3.1.1 mac
```
env GOOS=darwin GOARCH=amd64 go build ./
```
##### 3.1.2 linux
```
env GOOS=linux GOARCH=amd64 go build ./
```

#### 3.2 配置文件
请将配置文件放置在
```
./config/
```
或
```
/etc/p2p-sharer/
```
#### 3.3 启动
##### 3.3.1 启动tracker
```
./p2p-sharer tracker
```
对应配置文件 `config.tracker.yaml`
##### 3.3.2 启动node
```
./p2p-sharer node
```
对应配置文件 `config.node.yaml`


### 4. 完整的测试流程
假定tracker地址为
`192.168.10.200:35330`
node地址为
```
192.168.10.201:35331
192.168.10.202:35331
192.168.10.203:35331
```
在任何一个p2p文件分发过程中，首先需要有一个源，他拥有完整的文件(待分发的文件)

1) 针对要分发的文件制作种子文件
```
./p2p-sharer gen --tracker 192.168.10.200:35330 --filePath /tmp/Motrix-1.4.1.dmg --seedPath ./
```
会得到种子文件
`Motrix-1.4.1.dmg.seed`
2) 将
`Motrix-1.4.1.dmg.seed`   
`Motrix-1.4.1.dmg`
放置到192.168.10.201的下载目录(默认为`/tmp`)
在/tmp 中创建一个空文件
```
touch /tmp/Motrix-1.4.1.dmg.ok
```
这个文件的目的是为了标识Motrix-1.4.1.dmg是完整的

3) 重启192.168.10.201上的node服务
这时候node启动时，会扫描自己的下载文件路径，当发现
Motrix-1.4.1.dmg、Motrix-1.4.1.dmg.seed、Motrix-1.4.1.dmg.ok
会主动将该文件的分片信息发送给tracker
此时在可以调用    
查看tracker拥有的全局文件分片信息
```
192.168.10.200:35330/v1/probe
```
查看node拥有的文件分片信息
```
192.168.10.201:35331/v1/probe
```
4) 提交文件下载任务
```
curl -XPOST -d'{
	"seedFile": "http://192.168.10.199/download/Motrix-1.4.1.dmg.seed"
}' http://192.168.10.202:35331/v1/task
```

```
curl -XPOST -d'{
	"seedFile": "http://192.168.10.199:28080/download/Motrix-1.4.1.dmg.seed"
}' http://192.168.10.203:35331/v1/task
```
`seedFile`是种子文件的URL地址，请自行准备

node在下载文件的过程中，每完成一个分片都会向tracker汇报1次。
并且在上报的内容中，还包含有下载进度信息。观察tracker日志即可看到。


### 5. 也许你还能做得更好？
- [ ] 控制并发下载的线程数量(下载worker请求外部节点的并发线程数量）
- [ ] 控制并发下载的线程数量(对其它节点提供服务）
- [ ] 将下载目录和下载完成目录分开
- [ ] 允许通过指令触发对下载目录的扫描
- [ ] tracker的高可用改造
- [ ] 内网穿透
- [ ] nodeList接口 节点优先返回与请求节点本区域的节点
