package pagination

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

const defaultPageSize = 20
const maximumPageSize = 100
const minimumPageSize = 10
const defaultPage = 1

type Pagination struct {
	r *http.Request
}

func NewPagination(r *http.Request) *Pagination {
	return &Pagination{r}
}

func (p Pagination) Size() uint64 {
	var size uint64 = defaultPageSize
	if len(p.r.URL.Query()["size"]) > 0 {
		size, _ = strconv.ParseUint(p.r.URL.Query()["size"][0], 10, 0)
		if size > maximumPageSize {
			size = maximumPageSize
		}
		if size < minimumPageSize {
			size = minimumPageSize
		}
	}

	return size
}

func (p Pagination) Offset() uint64 {
	var offset uint64 = 0
	if len(p.r.URL.Query()["page"]) > 0 {
		page, _ := strconv.ParseUint(p.r.URL.Query()["page"][0], 10, 0)
		if page < defaultPage {
			page = defaultPage
		}
		offset = (page - 1) * p.Size()
	}

	return offset
}

func (p Pagination) TotalPages(count int) int {
	totalPages := count/int(p.Size()) + 1

	if totalPages > 1 {
		return totalPages
	}

	return 1
}

func (p Pagination) CurrentPage(count int) (int, error) {
	// get page from query parameter
	queryPage := int64(p.FirstPage())
	if len(p.r.URL.Query()["page"]) > 0 {
		queryPage, _ = strconv.ParseInt(p.r.URL.Query()["page"][0], 10, 0)
	}

	firstPage := p.FirstPage()
	lastPage := p.TotalPages(count)

	if int(queryPage) > lastPage || int(queryPage) < firstPage {
		return 0, errors.New(fmt.Sprintf("Requested page must be between 1 and %d", lastPage))
	}

	return int(queryPage), nil
}

func (p Pagination) PreviousPage(count int) (int, error) {
	prevPage := defaultPage
	currentPage, err := p.CurrentPage(count)
	if err != nil {
		return 0, err
	}

	if currentPage > 1 {
		prevPage = currentPage - 1
	}

	return prevPage, nil
}

func (p Pagination) NextPage(count int) (int, error) {
	currentPage, err := p.CurrentPage(count)
	if err != nil {
		return 0, err
	}

	return currentPage + 1, nil
}

func (p Pagination) FirstPage() int {
	return defaultPage
}
