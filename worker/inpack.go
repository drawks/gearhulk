package worker

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"

	rt "github.com/drawks/gearhulk/pkg/runtime"
)

// Worker side job
type inPack struct {
	dataType             rt.PT
	data                 []byte
	handle, uniqueId, fn string
	a                    *agent
}

// Create a new job
func getInPack() *inPack {
	return &inPack{}
}

func (inpack *inPack) Data() []byte {
	return inpack.data
}

func (inpack *inPack) Fn() string {
	return inpack.fn
}

func (inpack *inPack) Handle() string {
	return inpack.handle
}

func (inpack *inPack) UniqueId() string {
	return inpack.uniqueId
}

func (inpack *inPack) Err() error {
	if inpack.dataType == rt.PT_Error {
		return getError(inpack.data)
	}
	return nil
}

// Send some datas to client.
// Using this in a job's executing.
func (inpack *inPack) SendData(data []byte) {
	outpack := getOutPack()
	outpack.dataType = rt.PT_WorkData
	hl := len(inpack.handle)
	l := hl + len(data) + 1
	outpack.data = rt.NewBuffer(l)
	copy(outpack.data, []byte(inpack.handle))
	copy(outpack.data[hl+1:], data)
	inpack.a.write(outpack)
}

func (inpack *inPack) SendWarning(data []byte) {
	outpack := getOutPack()
	outpack.dataType = rt.PT_WorkWarning
	hl := len(inpack.handle)
	l := hl + len(data) + 1
	outpack.data = rt.NewBuffer(l)
	copy(outpack.data, []byte(inpack.handle))
	copy(outpack.data[hl+1:], data)
	inpack.a.write(outpack)
}

// Update status.
// Tall client how many percent job has been executed.
func (inpack *inPack) UpdateStatus(numerator, denominator int) {
	n := []byte(strconv.Itoa(numerator))
	d := []byte(strconv.Itoa(denominator))
	outpack := getOutPack()
	outpack.dataType = rt.PT_WorkStatus
	hl := len(inpack.handle)
	nl := len(n)
	dl := len(d)
	outpack.data = rt.NewBuffer(hl + nl + dl + 2)
	copy(outpack.data, []byte(inpack.handle))
	copy(outpack.data[hl+1:], n)
	copy(outpack.data[hl+nl+2:], d)
	inpack.a.write(outpack)
}

// Decode job from byte slice
func decodeInPack(data []byte) (inpack *inPack, l int, err error) {
	if len(data) < rt.MinPacketLength { // valid package should not less 12 bytes
		err = fmt.Errorf("Invalid data: %v", data)
		return
	}
	dl := int(binary.BigEndian.Uint32(data[8:12]))
	if len(data) < (dl + rt.MinPacketLength) {
		err = fmt.Errorf("Not enough data: %v", data)
		return
	}
	dt := data[rt.MinPacketLength : dl+rt.MinPacketLength]
	if len(dt) != int(dl) { // length not equal
		err = fmt.Errorf("Invalid data: %v", data)
		return
	}
	inpack = getInPack()
	inpack.dataType, err = rt.NewPT(binary.BigEndian.Uint32(data[4:8]))
	if err != nil {
		return
	}
	switch inpack.dataType {
	case rt.PT_JobAssign:
		s := bytes.SplitN(dt, []byte{'\x00'}, 3)
		if len(s) == 3 {
			inpack.handle = string(s[0])
			inpack.fn = string(s[1])
			inpack.data = s[2]
		}
	case rt.PT_JobAssignUniq:
		s := bytes.SplitN(dt, []byte{'\x00'}, 4)
		if len(s) == 4 {
			inpack.handle = string(s[0])
			inpack.fn = string(s[1])
			inpack.uniqueId = string(s[2])
			inpack.data = s[3]
		}
	default:
		inpack.data = dt
	}
	l = dl + rt.MinPacketLength
	return
}
