package snowflakes_test

import (
	"testing"

	"github.com/Amqp-prtcl/snowflakes"
)

func TestIfNull(t *testing.T) {
	var id1 snowflakes.ID

	t.Log(id1)

	if id1 != 0 {
		t.Fail()
	}

	if !id1.IsNull() {
		t.Fail()
	}

	node := snowflakes.NewNode(0)

	id2 := node.NewID()
	t.Log(id2)

	if id2 == 0 {
		t.Fail()
	}

	if id2.IsNull() {
		t.Fail()
	}
}
