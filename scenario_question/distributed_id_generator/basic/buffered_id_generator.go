package basic

import (
	"errors"
	"time"
)

type BufferedSnowflake struct {
	generator *Snowflake
	idChan    chan int64
	bufSize   int
}

func NewBufferedSnowflake(workerId int64, startEpoch int64, bufSize int) (*BufferedSnowflake, error) {
	generator, err := NewSnowflake(workerId, startEpoch)
	if err != nil {
		return nil, err
	}

	bf := &BufferedSnowflake{
		generator: generator,
		idChan:    make(chan int64, bufSize),
		bufSize:   bufSize,
	}

	go bf.generateIds()
	return bf, nil
}

func (bf *BufferedSnowflake) generateIds() {
	for {
		// 批量生成ID
		ids := make([]int64, 0, bf.bufSize/2)
		for i := 0; i < bf.bufSize/2; i++ {
			id, err := bf.generator.NextId()
			if err != nil {
				time.Sleep(time.Millisecond)
				break
			}
			ids = append(ids, id)
		}

		// 将生成的ID放入缓冲通道
		for _, id := range ids {
			bf.idChan <- id
		}
	}
}

func (bf *BufferedSnowflake) NextId() (int64, error) {
	select {
	case id := <-bf.idChan:
		return id, nil
	case <-time.After(time.Millisecond * 100):
		return 0, errors.New("get id timeout")
	}
}
