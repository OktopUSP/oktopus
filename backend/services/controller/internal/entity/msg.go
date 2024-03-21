package entity

type DataType interface {
	[]map[string]interface{} | *string | Device
}

type MsgAnswer[T DataType] struct {
	Code int
	Msg  T
}
