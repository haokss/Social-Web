package serializer

import "todo_list/model"

func BuildColleaguesMapView(colleagues []model.Colleague) []map[string]interface{} {
	var result []map[string]interface{}
	for _, c := range colleagues {
		result = append(result, map[string]interface{}{
			"id":         c.ID,
			"name":       c.Name,
			"relation":   "", // 没有 relation 字段，填空
			"is_set_map": c.IsSetMap,
		})
	}
	return result
}
