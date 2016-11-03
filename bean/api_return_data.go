package bean

type Page struct {
	Total     int64 `json:"total,omitempty"`
	PageTotal int   `json:"page_total,omitempty"`
	Page      int   `json:"page,omitempty"`
	PageSize  int   `json:"page_size,omitempty"`
}

type ApiData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
	Ext  string      `json:"ext,omitempty"`
}

type ApiDataWithPage struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
	Page Page        `json:"page,omitempty"`
	Ext  string      `json:"ext,omitempty"`
}
