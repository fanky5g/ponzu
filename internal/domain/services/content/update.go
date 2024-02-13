package content

func (s *service) UpdateContent(entityType, entityId string, update map[string]interface{}) (interface{}, error) {
	return s.repository.UpdateEntity(entityType, entityId, update)
}
