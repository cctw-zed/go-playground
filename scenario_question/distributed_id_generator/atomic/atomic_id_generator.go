package atomic

import (
	"errors"
	"sync/atomic"
	"time"
)

const (
	workerBits     = 10                        // 机器ID位数
	sequenceBits   = 12                        // 序列号位数
	workerMax      = -1 ^ (-1 << workerBits)   // 机器ID最大值
	sequenceMask   = -1 ^ (-1 << sequenceBits) // 序列号掩码
	timestampShift = workerBits + sequenceBits // 时间戳左移位数
	workerShift    = sequenceBits              // 机器ID左移位数
)

type Snowflake struct {
	sequence   atomic.Int64 // 使用原子操作的序列号
	timestamp  atomic.Int64 // 上次的时间戳
	workerId   int64        // 机器ID
	startEpoch int64        // 起始时间戳
}

func NewSnowflake(workerId int64, startEpoch int64) (*Snowflake, error) {
	// 检查机器ID是否合法
	if workerId < 0 || workerId > workerMax {
		return nil, errors.New("worker ID excess of quantity")
	}

	var seq, ts atomic.Int64
	seq.Store(0)
	ts.Store(0)
	return &Snowflake{
		sequence:   seq,
		timestamp:  atomic.Int64{},
		workerId:   0,
		startEpoch: 0,
	}, nil
}

func (s *Snowflake) NextId() (int64, error) {
	// 获取当前时间戳
	now := time.Now().UnixNano() / 1000000

	// CAS 操作更新时间戳和序列号
	for {
		lastTimestamp := s.timestamp.Load()
		sequence := s.sequence.Load()

		if now == lastTimestamp {
			// 同一毫秒内，尝试更新序列号
			nextSequence := (sequence + 1) & sequenceMask
			if nextSequence == 0 {
				// 序列号用完，等待下一毫秒
				for now <= lastTimestamp {
					now = time.Now().UnixNano() / 1000000
				}
				continue
			}
			if s.sequence.CompareAndSwap(sequence, nextSequence) {
				break
			}
			continue
		}

		// 不同毫秒，重置序列号
		if now > lastTimestamp {
			if s.timestamp.CompareAndSwap(lastTimestamp, now) {
				s.sequence.Store(0)
				break
			}
			continue
		}

		// 时钟回拨
		return 0, errors.New("clock moved backwards")
	}

	return ((now - s.startEpoch) << timestampShift) |
		(s.workerId << workerShift) |
		s.sequence.Load(), nil
}
