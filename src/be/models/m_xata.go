package models

type Data struct {
	Id string `json:"id,omitempty"`
	Parent   string `json:"parent,omitempty"`
	Children []string `json:"children,omitempty"`
}

type DataRequest struct {
	Parent   string `json:"parent,omitempty"`
	Children []string `json:"children,omitempty"`
}

type DataQuery struct {
	Columns  []string `json:"columns,omitempty"`
}

type DataResponse struct {
	Id string `json:"id,omitempty"`
}
