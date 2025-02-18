package basic

import (
	"errors"
	"sync/atomic"
	"time"
)

type RingBuffer struct {
	buffer []int64
	size   uint32
	mask   uint32
	cursor atomic.Uint32
	next   atomic.Uint32
}

func NewRingBuffer(size uint32) *RingBuffer {
	// 确保 size 是 2 的幂
	size = NextPowerOfTwo(size)
	return &RingBuffer{
		buffer: make([]int64, size),
		size:   size,
		mask:   size - 1,
	}
}

type RingSnowflake struct {
	generator *Snowflake
	ring      *RingBuffer
	padding   [56]byte // 避免false sharing
}

func (rs *RingSnowflake) produceIds() {
	for {
		cursor := rs.ring.cursor.Load()
		next := rs.ring.next.Load()

		// 检查是否还有空间
		if next-cursor >= rs.ring.size {
			time.Sleep(time.Microsecond)
			continue
		}

		// 生成新ID
		id, err := rs.generator.NextId()
		if err != nil {
			time.Sleep(time.Millisecond)
			continue
		}

		// 写入 RingBuffer
		idx := next & rs.ring.mask
		rs.ring.buffer[idx] = id
		rs.ring.next.Add(1)
	}
}

func (rs *RingSnowflake) NextId() (int64, error) {
	var cursor uint32
	var next uint32

	for {
		cursor = rs.ring.cursor.Load()
		next = rs.ring.next.Load()

		if cursor == next {
			return 0, errors.New("ring buffer is empty")
		}

		if rs.ring.cursor.CompareAndSwap(cursor, cursor+1) {
			return rs.ring.buffer[cursor&rs.ring.mask], nil
		}
	}
}

func NextPowerOfTwo(size uint32) uint32 {
	if size <= 1 {
		return 1
	}
	size--
	size |= size >> 1
	size |= size >> 2
	size |= size >> 4
	size |= size >> 8
	size |= size >> 16
	size++
	return size
}
