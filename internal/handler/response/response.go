package response

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type ResponseError struct {
	Meta  Meta        `json:"meta"`
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}

type ResponsePagination struct {
	Meta MetaPagination `json:"meta"`
	Data interface{}    `json:"data"`
}

type Meta struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type MetaPagination struct {
	Code          int    `json:"code"`
	Status        string `json:"status"`
	Message       string `json:"message"`
	TotalFiltered int    `json:"totalFiltered"`
	TotalRecords  int    `json:"totalRecords"`
	Page          int    `json:"page"`
	PerPage       int    `json:"perPage"`
}
