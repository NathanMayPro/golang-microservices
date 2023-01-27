package datamodel


type Document struct {
	Filename  string `json:"filename"`
	Content   string `json:"content"`
	Extension string `json:"extension"`
	Uid       string `json:"uid"`
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	//Data    interface{} `json:"data"`
}