package gormx

// Pagination 分页参数处理
func Pagination(size, page int) (int, int) {
	if size > 100 {
		size = 100
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * size
	return offset, size
}
