## node

### 1. 接收下载任务

### Request

#### method
POST

#### path

/v1/task

#### params

|名称|位置|参数类型|是否必须|说明|备注|
|:---|:---|:---|:---|:---|:---|
|seedFile|body|string|是|种子文件地址||


#### body
```
{
	"seedFile": "http://dev3:28080/download/Motrix-1.4.1.dmg_p2p"
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
### 2. 获取分片数据(供其他node调用)

### Request

#### method
GET

#### path

/v1/pieceData

#### params

|名称|位置|参数类型|是否必须|说明|备注|
|:---|:---|:---|:---|:---|:---|
|pieceID|query|string|是|文件分片的ID，通过SHA1生成||


### Response

#### body
200
```
二进制数据流
```


