package tools

import (
	"math"
)

type Paging struct {
	Page      int64 //当前页面
	PageSize  int64 //每页条数
	Total     int64 //文章总条数
	PageCount int64 //分页后总页数
}

func CreatePaging(page, pagesize, total int64) *Paging {
	if page < 1 {
		page = 1
	}
	if pagesize < 1 {
		pagesize = 10
	}

	page_count := math.Ceil(float64(total) / float64(pagesize))

	paging := new(Paging)
	paging.Page = page
	paging.PageSize = pagesize
	paging.Total = total
	paging.PageCount = int64(page_count)
	return paging
}
