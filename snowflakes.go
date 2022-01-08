package snowflakes

import (
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

type ID int64

func ParseID(s string) (ID, error) {
	id, err := strconv.ParseInt(s, 10, 64)
	return ID(id), err
}

func (id ID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

func (id ID) Time() time.Time {
	t := int64(id) >> timeShift
	return getTimeFromMillis(t)
}

func (id ID) Node() int64 {
	return (int64(id) & nodeMask) >> int64(nodeShift)
}

func (id ID) Step() int64 {
	return int64(id) & stepMask
}

func (id ID) bits() string {
	return strconv.FormatInt(int64(id), 2)
}

func (id ID) IsNull() bool {
	return id == 0
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
	return ID(id)
}
