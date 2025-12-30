package snowflake

import (
	"fmt"
	"sync"
	"time"
)

const (
	// Epoch is the custom epoch (2024-01-01 00:00:00 UTC)
	// This should match the Java version's epoch for compatibility
	epoch int64 = 1704067200000

	// Bit allocation
	workerIDBits     = 5
	datacenterIDBits = 5
	sequenceBits     = 12

	// Max values
	maxWorkerID     = -1 ^ (-1 << workerIDBits)
	maxDatacenterID = -1 ^ (-1 << datacenterIDBits)
	maxSequence     = -1 ^ (-1 << sequenceBits)

	// Bit shifts
	workerIDShift     = sequenceBits
	datacenterIDShift = sequenceBits + workerIDBits
	timestampShift    = sequenceBits + workerIDBits + datacenterIDBits
)

// Snowflake is a distributed unique ID generator
type Snowflake struct {
	mu           sync.Mutex
	workerID     int64
	datacenterID int64
	sequence     int64
	lastTime     int64
}

var (
	instance *Snowflake
	once     sync.Once
)

// Init initializes the snowflake generator with worker and datacenter IDs
func Init(workerID, datacenterID int64) error {
	if workerID < 0 || workerID > maxWorkerID {
		return fmt.Errorf("worker ID must be between 0 and %d", maxWorkerID)
	}
	if datacenterID < 0 || datacenterID > maxDatacenterID {
		return fmt.Errorf("datacenter ID must be between 0 and %d", maxDatacenterID)
	}

	once.Do(func() {
		instance = &Snowflake{
			workerID:     workerID,
			datacenterID: datacenterID,
			sequence:     0,
			lastTime:     0,
		}
	})

	return nil
}

// GetInstance returns the singleton instance
func GetInstance() *Snowflake {
	if instance == nil {
		// Default initialization with worker ID 1, datacenter ID 1
		Init(1, 1)
	}
	return instance
}

// NextID generates the next unique ID
func (s *Snowflake) NextID() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()

	if now < s.lastTime {
		return 0, fmt.Errorf("clock moved backwards, refusing to generate ID for %d milliseconds", s.lastTime-now)
	}

	if now == s.lastTime {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			// Sequence overflow, wait for next millisecond
			for now <= s.lastTime {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTime = now

	id := ((now - epoch) << timestampShift) |
		(s.datacenterID << datacenterIDShift) |
		(s.workerID << workerIDShift) |
		s.sequence

	return id, nil
}

// NextIDString generates the next unique ID as string
func (s *Snowflake) NextIDString() (string, error) {
	id, err := s.NextID()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", id), nil
}

// GenerateID is a convenience function to generate a new ID
func GenerateID() (int64, error) {
	return GetInstance().NextID()
}

// GenerateIDString is a convenience function to generate a new ID as string
func GenerateIDString() (string, error) {
	return GetInstance().NextIDString()
}

// MustGenerateID generates a new ID, panics on error
func MustGenerateID() int64 {
	id, err := GenerateID()
	if err != nil {
		panic(err)
	}
	return id
}

// MustGenerateIDString generates a new ID as string, panics on error
func MustGenerateIDString() string {
	id, err := GenerateIDString()
	if err != nil {
		panic(err)
	}
	return id
}

// ParseID extracts the components from a snowflake ID
func ParseID(id int64) (timestamp time.Time, datacenterID, workerID, sequence int64) {
	timestamp = time.UnixMilli(((id >> timestampShift) & 0x1FFFFFFFFFF) + epoch)
	datacenterID = (id >> datacenterIDShift) & maxDatacenterID
	workerID = (id >> workerIDShift) & maxWorkerID
	sequence = id & maxSequence
	return
}
