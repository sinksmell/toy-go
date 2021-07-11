package protocol

// JsonConvertReq request of json convert to go
type JsonConvertReq struct {
	JsonStr string `json:"json_str"`
	OmitEmpty bool `json:"omit_empty"`
}

// JsonConvertResp resp of json convert to go
type JsonConvertResp struct {
	GoStructStr string `json:"go_struct_str"`
}

// JsonConvertError error of json convert
type JsonConvertError struct {
	Message string `json:"message"`
}
