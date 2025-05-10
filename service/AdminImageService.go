package service

import (
	"todo_list/model"
	"todo_list/serializer"
)

func (service *UploadService) GetAllImages() serializer.Response {
	var images []model.Image
	if err := model.DB.Find(&images).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "查询所有图片失败",
		}
	}

	var imageItems []serializer.AdminImageItem
	for _, img := range images {
		imageItems = append(imageItems, serializer.AdminImageItem{
			URL:         img.Url,
			Name:        img.ImageName,
			Description: img.Description,
			UploaderID:  img.Uid, // 假设你的表有 Uid 字段
			// 可以扩展加 UploaderName，如果需要关联 user 表查询
		})
	}

	return serializer.Response{
		Status: 200,
		Data: serializer.AdminImagesSerializer{
			URLS: imageItems,
		},
		Msg: "获取所有图片成功！",
	}
}
