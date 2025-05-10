package serializer

type AdminImageItem struct {
	URL         string `json:"url"`
	Name        string `json:"image_name"`
	Description string `json:"description"`
	UploaderID  uint   `json:"uploader_id"` // 新增上传者ID
	// UploaderName string `json:"uploader_name"` // 如果你查了用户名
}

type AdminImagesSerializer struct {
	ID   uint             `json:"id"`
	URLS []AdminImageItem `json:"urls"`
}
