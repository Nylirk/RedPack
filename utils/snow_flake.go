package utils

import (
	"errors"
	"sync"
	"time"
)

const (
	epoch            int64 = 1704067200000
	workerIDBits           = 5  // 工作节点ID位数
	datacenterIDBits       = 5  // 数据中心ID位数
	sequenceBits           = 12 // 序列号位数

	maxWorkerID     = -1 ^ (-1 << workerIDBits)     // 最大工作节点ID
	maxDatacenterID = -1 ^ (-1 << datacenterIDBits) // 最大数据中心ID
	maxSequence     = -1 ^ (-1 << sequenceBits)     // 最大序列号

	workerIDShift     = sequenceBits
	datacenterIDShift = sequenceBits + workerIDBits
	timestampShift    = sequenceBits + workerIDBits + datacenterIDBits
)

type Snowflake struct {
	mu            sync.Mutex
	epoch         int64
	workerID      int64
	datacenterID  int64
	sequence      int64
	lastTimestamp int64
}

func NewSnowflake(workerID, datacenterID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, errors.New("worker ID 超出范围")
	}
	if datacenterID < 0 || datacenterID > maxDatacenterID {
		return nil, errors.New("数据中心ID 超出范围")
	}

	return &Snowflake{
		epoch:         epoch,
		workerID:      workerID,
		datacenterID:  datacenterID,
		sequence:      0,
		lastTimestamp: -1,
	}, nil
}

func (sf *Snowflake) Generate() (int64, error) {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	timestamp := time.Now().UnixNano() / 1000000

	if timestamp < sf.lastTimestamp {
		return 0, errors.New("时钟回退，拒绝生成ID")
	}

	if timestamp == sf.lastTimestamp {
		sf.sequence = (sf.sequence + 1) & maxSequence
		if sf.sequence == 0 {
			timestamp = sf.tilNextMillis(sf.lastTimestamp)
		}
	} else {
		sf.sequence = 0
	}

	sf.lastTimestamp = timestamp

	return ((timestamp - sf.epoch) << timestampShift) |
		(sf.datacenterID << datacenterIDShift) |
		(sf.workerID << workerIDShift) |
		sf.sequence, nil
}

func (sf *Snowflake) tilNextMillis(lastTimestamp int64) int64 {
	timestamp := time.Now().UnixNano() / 1000000
	for timestamp <= lastTimestamp {
		timestamp = time.Now().UnixNano() / 1000000
	}
	return timestamp
}
