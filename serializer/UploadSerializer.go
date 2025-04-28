package serializer

// 头像序列化数据
type UploadAvaSerializer struct {
	ID  uint   `json:"id"`
	URL string `json:"url"`
}

// 图片序列化数据
type ImageItem struct {
	URL         string `json:"url"`
	Name        string `json:"image_name"`
	Description string `json:"description"`
}

type ImagesSerializer struct {
	ID   uint        `json:"id"`
	URLS []ImageItem `json:"urls"`
}
