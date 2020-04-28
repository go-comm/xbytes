package xbytes

import (
	"bytes"
	"strconv"
	"testing"
)

func Test_Gen(t *testing.T) {
	var buf bytes.Buffer
	for i := 0; i <= 32; i++ {
		buf.WriteString(strconv.Itoa(1 << uint(i)))
		buf.WriteString(",")
	}

	t.Log(buf.String())
}
func Test_RoundUp(t *testing.T) {
	c := roundUp(33)

	t.Log(c)
}

func Test_RoundLog2(t *testing.T) {

	t.Log(roundLog2(262135))
	t.Log(roundLog2(65537))

	t.Log(roundLog2(1025))
	t.Log(roundLog2(1024))

	t.Log(roundLog2(65))
	t.Log(roundLog2(64))
	t.Log(roundLog2(63))
	t.Log(roundLog2(33))
	t.Log(roundLog2(32))
	t.Log(roundLog2(17))

	t.Log(roundLog2(16))
	t.Log(roundLog2(15))
	t.Log(roundLog2(9))
	t.Log(roundLog2(8))

}
