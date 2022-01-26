package snowflakes_test

import (
	"fmt"
	"testing"

	"github.com/Amqp-prtcl/snowflakes"
)

func TestIfNull(t *testing.T) {
	node := snowflakes.NewNode(1)

	id := node.NewID()
	fmt.Printf("id: %v\n", id)
	fmt.Printf("valid: %v\n", id.IsValid())
	fmt.Printf("node: %v\n", id.Node())
	fmt.Printf("step: %v\n", id.Step())
	fmt.Printf("time: %v\n", id.Time())
	fmt.Printf("str: %v\n", id.String())
}
