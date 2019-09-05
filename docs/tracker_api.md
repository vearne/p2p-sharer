## tracker

### 1. 接收心跳

### Request

#### method
POST

#### path

/v1/heartBeat

#### params

|名称|位置|参数类型|是否必须|说明|备注|
|:---|:---|:---|:---|:---|:---|
|nodeId|body|string|是|节点ID||
|addr|body|string|是|节点地址|形如 192.168.1.100:3933|

#### body
```
{
	"nodeId": "D526D2F5-920F-42A7-B5B9-BC643AC6B71B",
	"addr": "192.168.1.100:3933"
}
```

### Response

#### body
http status code 200
```
{
	"code": "E000",
	"msg": "success"
}
```
### 2. 获取piece对应的node列表

### Request

#### method
GET

#### path

/v1/nodeList

#### params

|名称|位置|参数类型|是否必须|说明|备注|
|:---|:---|:---|:---|:---|:---|
|pieceID|query|string|是|文件分片的ID，通过SHA1生成||


### Response

#### body
200
```
{
	"nodes": [
		"192.168.100.1:3433",
		"192.168.101.1:3433",
		"192.168.102.1:3433"
	]
}
```



### 3. 上报分片信息

#### 说明
节点每下载完成一个分片就上报一次

### Request

#### method
GET

#### path

/v1/report

#### params

|名称|位置|参数类型|是否必须|说明|备注|
|:---|:---|:---|:---|:---|:---|
|pieceID|body|string|是|文件分片的ID，通过SHA1生成||
|nodeID|body|string|是|节点ID||
|progress|body|float|是|进度||
|file|body|string|是|文件名||


#### body
```
{

	"nodeID": "D526D2F5-920F-42A7-B5B9-BC643AC6B71B",
    "file": "Motrix-1.4.1.dmg",
	"pieceID": "92465aae88444d799891ad730877f0f8593f77be",
    "progress": 0.75
}
```

### Response

#### body
http status code 200
```
{
	"code": "E000",
	"msg": "success"
}
```