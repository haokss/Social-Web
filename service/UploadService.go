package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"todo_list/config"
	"todo_list/model"
	sse "todo_list/package/SSE"
	"todo_list/package/errorcode"
	"todo_list/package/utils"
	"todo_list/serializer"

	"github.com/gin-gonic/gin"
)

type UploadService struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

// 上传头像处理
func (service *UploadService) UploadAva(uid uint) serializer.Response {
	// 创建存储目录
	uploadDir := "./www/uploads/Ava"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, 0755)
	}
	// 生成唯一文件名, 并保存
	fileName, _ := utils.GenerateFileName(uid, service.File)
	filePath := filepath.Join(uploadDir, fileName)
	if err := saveUploadedFile(service.File, filePath); err != nil {
		return serializer.Response{
			Status: errorcode.AvaSaveError,
			Data:   nil,
			Msg:    "upload ava failed!",
		}
	}

	avaUrl := fmt.Sprintf("http://%s:%s/www/uploads/Ava/%s", config.HttpIp, config.NgnixImagePort, fileName)

	// 更新数据库
	result := model.DB.Model(&model.User{}).Where("id = ?", uid).Update("ava_url", avaUrl)
	if result.Error != nil {
		return serializer.Response{
			Status: 40003,
			Data:   nil,
			Msg:    "更新头像URL失败",
		}
	}

	return serializer.Response{
		Status: 200,
		Data: serializer.UploadAvaSerializer{
			ID:  uid,
			URL: avaUrl,
		},
		Msg: "upload ava success!",
	}
}

// 上传图片处理
func (service *UploadService) UploadImage(c *gin.Context, uid uint) serializer.Response {
	// 创建存储目录
	uploadDir := "./www/uploads/Image/" + strconv.Itoa(int(uid))
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, 0755)
	}

	// 生成唯一文件名
	fileName, _ := utils.GenerateFileName(uid, service.File)
	filePath := filepath.Join(uploadDir, fileName)

	// 保存文件
	if err := saveUploadedFile(service.File, filePath); err != nil {
		return serializer.Response{
			Status: 40001,
			Data:   nil,
			Msg:    "upload image failed!",
		}
	}

	fileURL := fmt.Sprintf("http://%s:%s/www/uploads/Image/%d/%s", config.HttpIp, config.NgnixImagePort, uid, fileName)

	// 保存到数据库
	image := model.Image{
		Uid:         uid,
		ImageName:   fileName,
		Url:         fileURL,
		Description: "",
		CreateTime:  time.Now(),
	}
	err := model.DB.Create(&image).Error
	if err != nil {
		fmt.Println("DB create error:", err)
		return serializer.Response{
			Status: 500,
			Msg:    "create failed!",
		}
	}
	// 通知管理员审核
	broker := c.MustGet("sseBroker").(*sse.Broker)
	broker.Notify(sse.Message{
		Event:     "instant_notification",
		Data:      map[string]interface{}{"title": "您有新的图片待审核：" + fileName, "content": fileName},
		TargetIDs: []uint{utils.AdminUid},
	})

	return serializer.Response{
		Status: 200,
		Data: serializer.UploadAvaSerializer{
			ID:  uid,
			URL: fileURL,
		},
		Msg: "upload ava success!",
	}
}

// 获取存储的所有图片
func (service *UploadService) GetImage(uid uint) serializer.Response {
	var images []model.Image
	if err := model.DB.Where("uid = ?", uid).Find(&images).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "查询图片失败",
		}
	}

	var imageItems []serializer.ImageItem
	for _, img := range images {
		imageItems = append(imageItems, serializer.ImageItem{
			URL:         img.Url,
			Name:        img.ImageName,
			Description: img.Description,
		})
	}

	return serializer.Response{
		Status: 200,
		Data: serializer.ImagesSerializer{
			ID:   uid,
			URLS: imageItems,
		},
		Msg: "获取图片成功！",
	}
}

// 功能函数保存文件
func saveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
