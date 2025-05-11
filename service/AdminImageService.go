package service

import (
	"todo_list/model"
	"todo_list/serializer"

	"github.com/jinzhu/gorm"
)

type AdminBatchAuditImageItem struct {
	ImageIDs []uint `json:"image_ids" binding:"required"` // 要审核的图片ID
	Uid      uint   `json:"uid"`                          // 可选：按用户审核
}

type AdminBatchAuditImageService struct {
	AuditItems []AdminBatchAuditImageItem `json:"audit_items" binding:"required"`
	IsChecked  int                        `json:"is_checked"` // 1=通过，2=拒绝
}

func (service *UploadService) GetAllImages(page, pageSize int, keyword string) serializer.Response {
	var images []model.Image
	var total int64

	query := model.DB.Model(&model.Image{}).Preload("User")

	// 关键词搜索（模糊匹配）
	if keyword != "" {
		likePattern := "%" + keyword + "%"
		query = query.Where("image_name LIKE ? OR description LIKE ?", likePattern, likePattern)
	}

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "查询图片总数失败",
		}
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&images).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "分页查询图片失败",
		}
	}

	var imageItems []serializer.AdminImageItem
	for _, img := range images {
		imageItems = append(imageItems, serializer.AdminImageItem{
			Id:          int(img.Model.ID),
			URL:         img.Url,
			Name:        img.ImageName,
			Description: img.Description,
			UploaderID:  img.Uid,
			IsChecked:   img.IsChecked,
			UserName:    img.User.UserName,
		})
	}

	return serializer.Response{
		Status: 200,
		Data: serializer.AdminImagesSerializer{
			URLS:  imageItems,
			Total: total,
		},
		Msg: "获取所有图片成功！",
	}
}

func (service *AdminBatchAuditImageService) BatchAudit() serializer.Response {
	db := model.DB

	err := db.Transaction(func(tx *gorm.DB) error {
		for _, item := range service.AuditItems {
			if len(item.ImageIDs) == 0 {
				continue
			}

			query := tx.Model(&model.Image{}).Where("id IN (?)", item.ImageIDs)
			if item.Uid != 0 {
				query = query.Where("uid = (?)", item.Uid)
			}

			if err := query.Update("is_checked", service.IsChecked).Error; err != nil {
				return err
			}
		}
		return nil
	})

	// 发送通知给用户

	if err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "批量审核失败",
		}
	}

	return serializer.Response{
		Status: 200,
		Msg:    "批量审核成功",
	}
}
