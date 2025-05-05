package serializer

import "todo_list/model"

type Relative struct {
	ID              uint        `json:"id"`
	Name            string      `json:"name"`
	Relation        string      `json:"relation"`
	Gender          string      `json:"gender"`
	Address         string      `json:"address"`
	Contact         string      `json:"contact"`
	Wechat          string      `json:"wechat"`
	HasDebtRelation bool        `json:"hasDebtRelation"`
	DebtType        string      `json:"debtType"`
	DebtProof       string      `json:"debtProof"`
	Note            string      `json:"note"`
	Avatar          string      `json:"avatar"`
	ParentID        *uint       `json:"parentId"`
	Children        []*Relative `json:"children,omitempty"`
}

func BuildRelative(node model.RelativeInfo) Relative {
	return Relative{
		ID:              node.ID,
		Name:            node.Name,
		Relation:        node.Relation,
		Gender:          node.Gender,
		Address:         node.Address,
		Contact:         node.Contact,
		Wechat:          node.WeChat,
		HasDebtRelation: node.HasDebtRelation,
		DebtType:        node.DebtType,
		DebtProof:       node.DebtProof,
		Note:            node.Note,
		Avatar:          node.Avatar,
		ParentID:        node.ParentID,
	}
}

func BuildRelativeTree(relatives []model.RelativeInfo) []Relative {
	idMap := make(map[uint]*Relative)

	rootNode := &Relative{
		ID:       0,
		Name:     "本人",
		Children: []*Relative{},
	}
	idMap[0] = rootNode

	// 第一步：全部转换为序列化对象
	for _, r := range relatives {
		node := Relative{
			ID:              r.ID,
			Name:            r.Name,
			Relation:        r.Relation,
			Gender:          r.Gender,
			Address:         r.Address,
			Contact:         r.Contact,
			Wechat:          r.WeChat,
			HasDebtRelation: r.HasDebtRelation,
			DebtType:        r.DebtType,
			DebtProof:       r.DebtProof,
			Note:            r.Note,
			Avatar:          r.Avatar,
			ParentID:        r.ParentID,
			Children:        []*Relative{},
		}
		idMap[r.ID] = &node
	}

	// 构建父子树结构
	for _, node := range idMap {
		if node.ParentID != nil {
			if parent, ok := idMap[*node.ParentID]; ok {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	return []Relative{*rootNode}
}
