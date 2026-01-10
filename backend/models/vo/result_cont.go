package vo

// ResultCont
// @Description: unified response format
type ResultCont struct {
	Code    int         `json:"code"`    //response code
	Message string      `json:"message"` //response message
	Data    interface{} `json:"data"`    //response data
}
