package entity

type DataType interface {
	[]map[string]interface{} | *string | Device | int64 | []Device
}

type MsgAnswer[T DataType] struct {
	Code int
	Msg  T
}
