package basic

import (
	"testing"
	"time"
)

func BenchmarkSnowflake(b *testing.B) {
	sf, _ := NewSnowflake(1, time.Now().UnixNano()/1000000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sf.NextId()
		}
	})
}

func BenchmarkBufferedSnowflake(b *testing.B) {
	sf, _ := NewBufferedSnowflake(1, time.Now().UnixNano()/1000000, 1024)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sf.NextId()
		}
	})
}
