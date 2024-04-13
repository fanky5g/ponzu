package content

func (s *service) UpdateContent(entityType, entityId string, update map[string]interface{}) (interface{}, error) {
	return s.repository(entityType).UpdateById(entityId, update)
}
