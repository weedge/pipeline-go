package frames

import "reflect"

// AppFrame is a user-defined frame.
type AppFrame struct {
	*BaseFrame
}

func NewAppFrame() *AppFrame {
	return &AppFrame{
		BaseFrame: NewBaseFrameWithName(reflect.TypeOf(AppFrame{}).Name()),
	}
}
