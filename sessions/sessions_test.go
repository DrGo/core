package sessions

import (
	"fmt"
	"testing"
)

func TestSessions_GenSessionID(t *testing.T) {
	ss := Sessions{}
	sid := ss.GenSessionID()
	fmt.Println(string(sid[:]))
}

func BenchmarkGenSessionID(b *testing.B) {
	ss := Sessions{}
	for i := 0; i < b.N; i++ {
		_ = ss.GenSessionID()
	}
}
