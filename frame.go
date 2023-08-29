package cvm

import "fmt"

type Frame struct {
	StackOffset int
	HeapOffset  int
	FrameOffset int
	ReturnIP    uint32
}

func (f *Frame) String() string {
	return fmt.Sprintf("#stack: %d, #heap: %d,#return: %d", f.StackOffset, f.HeapOffset, f.ReturnIP)
}
