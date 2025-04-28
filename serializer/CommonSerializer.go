package serializer

// 基础序列化
type Response struct {
	Status int         `json:"code"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
}

// 携带Token响应
type TokenData struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}

func BuildListResponse(items interface{}, total uint) Response {
	return Response{
		Status: 200,
		Data: DataList{
			Item:  items,
			Total: total,
		},
	}
}
