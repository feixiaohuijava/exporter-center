 
打包操作:
    set GOARCH=amd64
    set GOOS=linux
    go build 


    
    
关于前后端交流
type Response struct {
	// 代码
	Code int `json:"code" example:"200"`
	// 数据集
	Data interface{} `json:"data"`
	// 消息
	Msg string `json:"msg"`
}
