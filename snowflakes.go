package snowflakes

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var (
	stepBits uint8 = 12
	nodeBits uint8 = 10

	nodeShift = stepBits
	timeShift = nodeBits + stepBits

	stepMask int64 = -1 ^ (-1 << stepBits)
	nodeMask int64 = (-1 ^ (-1 << nodeBits)) << nodeShift

	ErrInvalidID error = fmt.Errorf("invalid ID format")

	epoch time.Time = time.Unix(0, 0)
)

func SetStepBits(StepBits uint8) {
	stepBits = StepBits
	updateValues()
}

func SetNodeBits(NodeBits uint8) {
	nodeBits = NodeBits
	updateValues()
}

func SetEpoch(Epoch time.Time) {
	epoch = Epoch
}

func GetEpoch() time.Time {
	return epoch
}

func updateValues() {
	nodeShift = stepBits
	timeShift = nodeBits + stepBits

	stepMask = -1 ^ (-1 << stepBits)
	nodeMask = (-1 ^ (-1 << nodeBits)) << nodeShift
}

func getMillisFromEpoch() int64 {
	return time.Since(epoch).Milliseconds()
}

func getTimeFromMillis(millis int64) time.Time {
	return epoch.Add(time.Millisecond * time.Duration(millis))
}

func toInt(id ID) (int64, bool) {
	i, err := strconv.ParseInt(id.String(), 10, 64)
	if err != nil {
		return 0, false
	}
	return i, true
}

func toStr(i int64) ID {
	return ID(fmt.Sprint(i))
}

type ID string

func ParseID(s string) (ID, error) {
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		err = ErrInvalidID
	}
	return toStr(id), err
}

func (id ID) String() string {
	return string(id)
}

func (id ID) Step() int64 {
	i, ok := toInt(id)
	if !ok {
		return 0
	}
	return i & stepMask
}

func (id ID) Node() int64 {
	i, ok := toInt(id)
	if !ok {
		return 0
	}
	return (i & nodeMask) >> nodeShift
}

func (id ID) Stamp() int64 {
	i, ok := toInt(id)
	if !ok {
		return 0
	}
	return i >> timeShift
}

func (id ID) Time() time.Time {
	i, ok := toInt(id)
	if !ok {
		return time.Now()
	}
	t := i >> timeShift
	return getTimeFromMillis(t)
}

/*
func (id ID) bits() string {
	i, ok := toInt(id)
	if !ok {
		return ""
	}
	return strconv.FormatInt(i, 2)
}
*/

func (id ID) IsValid() bool {
	_, ok := toInt(id)
	return ok && id != ""
}

type Node struct {
	last int64
	node int64
	step int64

	sync.Mutex
}

func NewNode(node int64) *Node {
	n := &Node{
		node: node,
	}
	return n
}

func (n *Node) NewID() ID {
	now := getMillisFromEpoch()
	n.Lock()
	if now == n.last {
		n.step = (n.step + 1) & stepMask
		if n.step == 0 {
			for now <= n.last {
				now = getMillisFromEpoch()
			}
		}
	} else {
		n.step = 0
	}
	n.last = now
	id := now<<timeShift | n.node<<nodeShift | n.step
	n.Unlock()
	return toStr(id)
}
