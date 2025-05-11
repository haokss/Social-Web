package serializer

type AdminImageItem struct {
	Id          int    `json:"id"`
	URL         string `json:"url"`
	Name        string `json:"image_name"`
	Description string `json:"description"`
	UploaderID  uint   `json:"uid"`
	IsChecked   int    `json:"is_checked"`
	UserName    string `json:"user_name"`
}

type AdminImagesSerializer struct {
	Total int64            `json:"total"`
	URLS  []AdminImageItem `json:"urls"`
}
