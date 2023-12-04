package main

type listReq struct {
	Dir string `json:"dir"`
}
type listFile struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir,omitempty"`
}
type listRes struct {
	List []listFile `json:"list"`
}

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type successRes struct {
	Message string `json:"message,omitempty"`
}
type errorRes struct {
	Error string `json:"error"`
}
