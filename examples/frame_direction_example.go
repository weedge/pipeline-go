package main

import (
	"fmt"

	"github.com/weedge/pipeline-go/pkg/processors"
)

func main() {
	// 示例展示 FrameDirection 枚举的 String() 方法
	upstream := processors.FrameDirectionUpstream
	downstream := processors.FrameDirectionDownstream

	fmt.Printf("FrameDirectionUpstream: %s\n", upstream.String())
	fmt.Printf("FrameDirectionDownstream: %s\n", downstream.String())

	// 也可以直接打印枚举值，会自动调用 String() 方法
	fmt.Printf("FrameDirectionUpstream: %v\n", upstream)
	fmt.Printf("FrameDirectionDownstream: %v\n", downstream)
}
