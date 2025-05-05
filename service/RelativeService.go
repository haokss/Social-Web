package service

import (
	"fmt"
	"strings"
	"todo_list/model"
	"todo_list/serializer"
)

type CreateRelativeService struct {
	Name            string `json:"name"`
	Label           string `json:"label"`
	Relation        string `json:"relation"`
	Gender          string `json:"gender"`
	Address         string `json:"address"`
	Contact         string `json:"contact"`
	Wechat          string `json:"wechat"`
	HasDebtRelation bool   `json:"hasDebtRelation"`
	DebtType        string `json:"debtType"`
	DebtProof       string `json:"debtProof"`
	Note            string `json:"note"`
	Avatar          string `json:"avatar"`
	ParentID        *uint  `json:"parentId"` // 可选
}

type ListRelativeService struct{}

type UpdateRelativeService struct {
	Name            string `json:"name"`
	Relation        string `json:"relation"`
	Gender          string `json:"gender"`
	Address         string `json:"address"`
	Contact         string `json:"contact"`
	Wechat          string `json:"wechat"`
	HasDebtRelation bool   `json:"hasDebtRelation"`
	DebtType        string `json:"debtType"`
	DebtProof       string `json:"debtProof"`
	Note            string `json:"note"`
	Avatar          string `json:"avatar"`
	ParentID        *uint  `json:"parentId"`
}

type DeleteRelativeService struct {
	IDs []uint `json:"ids" binding:"required"`
}

// 创建亲属关系
func (service *CreateRelativeService) CreateRelative(id uint) serializer.Response {
	node := model.RelativeInfo{
		Name:            service.Name,
		Relation:        service.Relation,
		Gender:          service.Gender,
		Address:         service.Address,
		Contact:         service.Contact,
		WeChat:          service.Wechat,
		HasDebtRelation: service.HasDebtRelation,
		DebtType:        service.DebtType,
		DebtProof:       service.DebtProof,
		Note:            service.Note,
		Avatar:          service.Avatar,
		ParentID:        service.ParentID,
	}

	if err := model.DB.Create(&node).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "创建失败: " + err.Error(),
		}
	}

	return serializer.Response{
		Status: 200,
		Data:   serializer.BuildRelative(node),
		Msg:    "创建成功",
	}
}

// 获取所有亲属信息
func (service *ListRelativeService) GetRelativeList(id uint) serializer.Response {
	var relatives []model.RelativeInfo
	if err := model.DB.Find(&relatives).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "获取亲属信息失败: " + err.Error(),
		}
	}

	tree := serializer.BuildRelativeTree(relatives)

	return serializer.Response{
		Status: 200,
		Data:   tree,
		Msg:    "获取成功",
	}
}

// 更新亲属信息
func (service *UpdateRelativeService) Update(id uint, tid string) serializer.Response {
	var relative model.RelativeInfo
	if err := model.DB.First(&relative, tid).Error; err != nil {
		return serializer.Response{
			Status: 404,
			Msg:    "亲属信息不存在",
		}
	}

	// 更新字段
	relative.Name = service.Name
	relative.Relation = service.Relation
	relative.Gender = service.Gender
	relative.Address = service.Address
	relative.Contact = service.Contact
	relative.WeChat = service.Wechat
	relative.HasDebtRelation = service.HasDebtRelation
	relative.DebtType = service.DebtType
	relative.DebtProof = service.DebtProof
	relative.Note = service.Note
	relative.Avatar = service.Avatar
	relative.ParentID = service.ParentID

	if err := model.DB.Save(&relative).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "更新失败: " + err.Error(),
		}
	}

	return serializer.Response{
		Status: 200,
		Msg:    "更新成功",
	}
}

func (service *DeleteRelativeService) DeleteRelative() serializer.Response {
	var allIDs []uint

	// 递归收集所有子节点ID
	var collectChildIDs func(ids []uint)
	collectChildIDs = func(ids []uint) {
		allIDs = append(allIDs, ids...)
		var children []model.RelativeInfo
		if err := model.DB.Where("parent_id IN ?", ids).Find(&children).Error; err != nil {
			return
		}
		if len(children) > 0 {
			var childIDs []uint
			for _, child := range children {
				childIDs = append(childIDs, child.ID)
			}
			collectChildIDs(childIDs)
		}
	}

	// 开始收集
	collectChildIDs(service.IDs)

	// 执行删除
	if len(allIDs) == 0 {
		return serializer.Response{
			Status: 400,
			Msg:    "无可删除的节点",
		}
	}

	placeholders := strings.TrimRight(strings.Repeat("?,", len(allIDs)), ",")
	// 将 []uint 转换成 []interface{} 以用于 Exec 参数
	args := make([]interface{}, len(allIDs))
	for i, id := range allIDs {
		args[i] = id
	}

	sql := fmt.Sprintf("DELETE FROM relative_info WHERE id IN (%s)", placeholders)
	if err := model.DB.Exec(sql, args...).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "删除失败: " + err.Error(),
		}
	}

	return serializer.Response{
		Status: 200,
		Msg:    "删除成功",
		Data: map[string]interface{}{
			"deleted_ids": allIDs,
		},
	}
}
