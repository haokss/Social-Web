package serializer

import "todo_list/model"

func BuildClassmatesMapView(classmates []model.Classmate) []map[string]interface{} {
	var result []map[string]interface{}
	for _, c := range classmates {
		result = append(result, map[string]interface{}{
			"id":         c.ID,
			"name":       c.Name,
			"relation":   "", // 没有 relation 字段，填空
			"is_set_map": c.IsSetMap,
		})
	}
	return result
}
