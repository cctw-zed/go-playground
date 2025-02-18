package basic

import (
	"errors"
	"sync"
	"time"
)

type Snowflake struct {
	mutex          sync.Mutex
	timestamp      int64 // 上次生成ID的时间戳
	workerId       int64 // 机器ID
	sequence       int64 // 序列号
	startEpoch     int64 // 起始时间戳
	backupSequence int64 // 备用序列号，用于时钟回拨
	maxBackTime    int64 // 最大允许回拨时间，比如 5ms
}

const (
	workerBits     = 10                        // 机器ID位数
	sequenceBits   = 12                        // 序列号位数
	workerMax      = -1 ^ (-1 << workerBits)   // 机器ID最大值
	sequenceMask   = -1 ^ (-1 << sequenceBits) // 序列号掩码
	timestampShift = workerBits + sequenceBits // 时间戳左移位数
	workerShift    = sequenceBits              // 机器ID左移位数
)

func NewSnowflake(workerId int64, startEpoch int64) (*Snowflake, error) {
	// 检查机器ID是否合法
	if workerId < 0 || workerId > workerMax {
		return nil, errors.New("worker ID excess of quantity")
	}

	return &Snowflake{
		timestamp:  0,
		workerId:   workerId,
		sequence:   0,
		startEpoch: startEpoch,
	}, nil
}

func (s *Snowflake) NextId() (int64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 获取当前时间戳
	now := time.Now().UnixNano() / 1000000 // 转换为毫秒

	if s.timestamp == now {
		// 同一毫秒内，序列号自增
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// 序列号用完，等待下一毫秒
			for now <= s.timestamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		// 不同毫秒，序列号重置
		s.sequence = 0
	}

	// 检查是否出现时钟回拨
	if now < s.timestamp {
		return 0, errors.New("clock moved backwards")
	}

	s.timestamp = now

	// 组合生成ID
	id := ((now - s.startEpoch) << timestampShift) |
		(s.workerId << workerShift) |
		s.sequence

	return id, nil
}
