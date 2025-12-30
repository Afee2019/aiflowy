package dto

// PageRequest represents pagination request parameters
type PageRequest struct {
	PageNumber int    `query:"pageNumber" json:"pageNumber"`
	PageSize   int    `query:"pageSize" json:"pageSize"`
	SortKey    string `query:"sortKey" json:"sortKey"`
	SortType   string `query:"sortType" json:"sortType"` // asc or desc
}

// GetPageNumber returns page number with default value
func (p *PageRequest) GetPageNumber() int {
	if p.PageNumber <= 0 {
		return 1
	}
	return p.PageNumber
}

// GetPage is an alias for GetPageNumber
func (p *PageRequest) GetPage() int {
	return p.GetPageNumber()
}

// GetPageSize returns page size with default value
func (p *PageRequest) GetPageSize() int {
	if p.PageSize <= 0 {
		return 10
	}
	if p.PageSize > 100 {
		return 100
	}
	return p.PageSize
}

// GetOffset returns the offset for SQL query
func (p *PageRequest) GetOffset() int {
	return (p.GetPageNumber() - 1) * p.GetPageSize()
}

// PageResponse represents pagination response
type PageResponse struct {
	PageNumber int         `json:"pageNumber"`
	PageSize   int         `json:"pageSize"`
	TotalRow   int64       `json:"totalRow"`
	TotalPage  int         `json:"totalPage"`
	Rows       interface{} `json:"rows"`
}

// NewPageResponse creates a new PageResponse
func NewPageResponse(pageNumber, pageSize int, totalRow int64, rows interface{}) *PageResponse {
	totalPage := int(totalRow) / pageSize
	if int(totalRow)%pageSize > 0 {
		totalPage++
	}
	return &PageResponse{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalRow:   totalRow,
		TotalPage:  totalPage,
		Rows:       rows,
	}
}

// IDRequest represents a request with single ID
type IDRequest struct {
	ID int64 `json:"id" query:"id"`
}

// IDsRequest represents a request with multiple IDs
type IDsRequest struct {
	IDs []int64 `json:"ids"`
}
