# eDB

### 基本介绍

-  源于阿里adb客户端的实现逻辑，自行封装的批处理工具。
-  批处理写库的go客户端，目的是为了实现数据库写操作的批量化，结构化，提高效率，规范统一代码结构.


### 使用方法
```
go get -u github.com/et-zone/eDB

```

### 说明
- 目前不支持并发批处理，并发批处理本身存在安全隐患，可能会导致复写，或者写冲突问题。
- 针对批处理业务不走并发更安全
- 后期如有更好的方案会更新项目

### 数据表配置json文件（必须配置）
```
{
    "tbname1":["field1","field2","field3"],
    "tbname":[]
    
}
```

### example
````
package main

import (
	"github.com/et-zone/eDB"
)

func main() {
	cfg := &eDB.EConfig{
		UserName:  "root",
		PassWord:  "mysql",
		Addr:      "127.0.0.1",
		Port:      3306,
		DB:        "test",
		TableFile: "./tbinfo.json",
	}
	row := eDB.NewRow()
	row.SetColumn(0, 0)
	row.SetColumn(1, nil)
	row.SetColumn(2, "e空间")
	// fmt.Println(row.String())

	client := eDB.InitClient(cfg)
	client.AddRow("sk", row)
	// fmt.Println(client.GettableNanme("aaa"))
	// client.FlushTx("sk")
	client.FlushAll()
}

````