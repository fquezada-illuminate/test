package pagination

import (
	"fmt"
	"net/http/httptest"
	"testing"
)

func TestPagination_Size(t *testing.T) {
	// Test minimum page size
	r := httptest.NewRequest("GET", fmt.Sprintf("/?size=%d", minimumPageSize-1), nil)
	page := NewPagination(r)
	size := page.Size()
	if size != minimumPageSize {
		t.Errorf("Expected page size %d, got %d", minimumPageSize, size)
	}

	// Test maximum page size
	r = httptest.NewRequest("GET", fmt.Sprintf("/?size=%d", maximumPageSize+1), nil)
	page = NewPagination(r)
	size = page.Size()
	if size != maximumPageSize {
		t.Errorf("Expected page size %d, got %d", maximumPageSize, size)
	}

	// Test page size in range
	var inputSize uint64 = 25
	r = httptest.NewRequest("GET", fmt.Sprintf("/?size=%d", inputSize), nil)
	page = NewPagination(r)
	size = page.Size()
	if size != inputSize {
		t.Errorf("Expected page size %d, got %d", inputSize, size)
	}

	// Test default page size
	r = httptest.NewRequest("GET", "/", nil)
	page = NewPagination(r)
	size = page.Size()
	if size != defaultPageSize {
		t.Errorf("Expected page size %d, got %d", defaultPageSize, size)
	}
}

func TestPagination_Offset(t *testing.T) {
	// Test default page
	r := httptest.NewRequest("GET", fmt.Sprintf("/?page=%d", defaultPage-1), nil)
	page := NewPagination(r)
	offset := page.Offset()
	if offset != 0 {
		t.Errorf("Expected page offset 0, got %d", offset)
	}

	// Test page in range
	var inputPage uint64 = 5
	r = httptest.NewRequest("GET", fmt.Sprintf("/?page=%d", inputPage), nil)
	page = NewPagination(r)
	offset = page.Offset()
	expectedOffset := (inputPage - 1) * page.Size()
	if offset != expectedOffset {
		t.Errorf("Expected page offset %d, got %d", expectedOffset, offset)
	}
}

func TestPagination_TotalPages(t *testing.T) {
	// Test default page size
	r := httptest.NewRequest("GET", "/", nil)
	page := NewPagination(r)
	count := 100
	totalPages := page.TotalPages(count)
	expectedOutput := count/defaultPageSize + 1
	if totalPages != expectedOutput {
		t.Errorf("Expected total pages %d, got %d", expectedOutput, totalPages)
	}

	// Test minimum page size
	inputSize := minimumPageSize - 1
	r = httptest.NewRequest("GET", fmt.Sprintf("/?size=%d", inputSize), nil)
	page = NewPagination(r)
	totalPages = page.TotalPages(count)
	expectedOutput = count/minimumPageSize + 1
	if totalPages != expectedOutput {
		t.Errorf("Expected total pages %d, got %d", expectedOutput, totalPages)
	}

	// Test minimum page size
	inputSize = maximumPageSize + 1
	r = httptest.NewRequest("GET", fmt.Sprintf("/?size=%d", inputSize), nil)
	page = NewPagination(r)
	totalPages = page.TotalPages(count)
	expectedOutput = count/maximumPageSize + 1
	if totalPages != expectedOutput {
		t.Errorf("Expected total pages %d, got %d", expectedOutput, totalPages)
	}

	// Test page size in range
	inputSize = (maximumPageSize + minimumPageSize) / 2
	r = httptest.NewRequest("GET", fmt.Sprintf("/?size=%d", inputSize), nil)
	page = NewPagination(r)
	totalPages = page.TotalPages(count)
	expectedOutput = count/inputSize + 1
	if totalPages != expectedOutput {
		t.Errorf("Expected total pages %d, got %d", expectedOutput, totalPages)
	}

	// Test count = 0
	inputSize = (maximumPageSize + minimumPageSize) / 2
	r = httptest.NewRequest("GET", fmt.Sprintf("/?size=%d", inputSize), nil)
	page = NewPagination(r)
	count = 0
	totalPages = page.TotalPages(count)
	expectedOutput = count/inputSize + 1
	if totalPages != expectedOutput {
		t.Errorf("Expected total pages %d, got %d", expectedOutput, totalPages)
	}
}

func TestPagination_CurrentPage(t *testing.T) {
	// Test default/minimum case
	r := httptest.NewRequest("GET", "/", nil)
	page := NewPagination(r)
	count := 100
	currentPage, _ := page.CurrentPage(count)
	expectedPage := 1
	if currentPage != expectedPage {
		t.Errorf("Expected current page %d, got %d", expectedPage, currentPage)
	}

	// Test error case with page input
	inputPage := 2
	r = httptest.NewRequest("GET", fmt.Sprintf("/?page=%d", inputPage), nil)
	page = NewPagination(r)
	count = 10
	currentPage, err := page.CurrentPage(count)
	if err == nil {
		t.Errorf("Expected error, got %d and error message %q", currentPage, err)
	}

	// Test error case with negative page input
	inputPage = -1
	r = httptest.NewRequest("GET", fmt.Sprintf("/?page=%d", inputPage), nil)
	page = NewPagination(r)
	count = 10
	currentPage, err = page.CurrentPage(count)
	if err == nil {
		t.Errorf("Expected current page %d, got %d and error message %q", defaultPage, currentPage, err)
	}

	// Test pass case with page input
	inputPage = 1
	r = httptest.NewRequest("GET", fmt.Sprintf("/?page=%d", inputPage), nil)
	page = NewPagination(r)
	count = 10
	currentPage, _ = page.CurrentPage(count)
	if currentPage != inputPage {
		t.Errorf("Expected current page %d, got %d", inputPage, currentPage)
	}
}

func TestPagination_PreviousPage(t *testing.T) {
	// Test default case
	r := httptest.NewRequest("GET", "/", nil)
	page := NewPagination(r)
	count := 100
	prevPage, _ := page.PreviousPage(count)
	expectedPage := 1
	if prevPage != expectedPage {
		t.Errorf("Expected previous page %d, got %d", expectedPage, prevPage)
	}

	// Test pass case with page input
	inputPage := 5
	r = httptest.NewRequest("GET", fmt.Sprintf("/?page=%d", inputPage), nil)
	page = NewPagination(r)
	expectedPage = inputPage - 1
	prevPage, _ = page.PreviousPage(count)
	if prevPage != expectedPage {
		t.Errorf("Expected previous page %d, got %d", expectedPage, prevPage)
	}

	// Test error case with page input
	inputPage = 0
	r = httptest.NewRequest("GET", fmt.Sprintf("/?page=%d", inputPage), nil)
	page = NewPagination(r)
	prevPage, err := page.PreviousPage(count)
	if err == nil {
		t.Errorf("Expected error, got %d and error message %q", prevPage, err)
	}
}

func TestPagination_NextPage(t *testing.T) {
	// Test default case
	r := httptest.NewRequest("GET", "/", nil)
	page := NewPagination(r)
	count := 100
	nextPage, _ := page.NextPage(count)
	expectedPage := 2
	if nextPage != expectedPage {
		t.Errorf("Expected next page %d, got %d", expectedPage, nextPage)
	}

	// Test pass case with page input
	inputPage := 5
	r = httptest.NewRequest("GET", fmt.Sprintf("/?page=%d", inputPage), nil)
	page = NewPagination(r)
	expectedPage = inputPage + 1
	nextPage, _ = page.NextPage(count)
	if nextPage != expectedPage {
		t.Errorf("Expected next page %d, got %d", expectedPage, nextPage)
	}

	// Test error case with page input
	inputPage = 0
	r = httptest.NewRequest("GET", fmt.Sprintf("/?page=%d", inputPage), nil)
	page = NewPagination(r)
	nextPage, err := page.NextPage(count)
	if err == nil {
		t.Errorf("Expected error, got %d and error message %q", nextPage, err)
	}
}
