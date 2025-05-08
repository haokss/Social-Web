package serializer

import "todo_list/model"


func BuildFriendsMapView(friends []model.Friend) []map[string]interface{} {
	var result []map[string]interface{}
	for _, f := range friends {
		result = append(result, map[string]interface{}{
			"id":         f.ID,
			"name":       f.Name,
			"relation":   "", // Friend 没有 relation，直接填 ""
			"is_set_map": f.IsSetMap,
		})
	}
	return result
}