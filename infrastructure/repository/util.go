package repository

const (
	defaultSize = 30
)

func CalcPage(page, pageSize int32) (offset, limit int) {
	currPage := page - 1
	if currPage < 0 {
		currPage = 1
	}

	if pageSize < 0 || pageSize > defaultSize {
		pageSize = defaultSize
	}
	offset = int(currPage * defaultSize)
	limit = int(pageSize)

	return
}
