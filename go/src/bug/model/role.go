package model

type Table_roles struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	rolelist string `json:"rolelist"`
}

type Data_roles struct {
	Id       int64    `json:"id"`
	Name     string   `json:"name"`
	RoleList []string `json:"rolelist"`
	Code     int      `json:"statuscode"`
}

type List_roles struct {
	DataList []*Data_roles `json:"datalist"`
	Code     int           `json:"statuscode"`
}

type Get_roles struct {
	Roles []string `json:"roles"`
	Code  int      `json:"statuscode"`
}

type Updata_role struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Oldname string `json:"oldname"`
}
