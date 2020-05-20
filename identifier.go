package core

import (
	"github.com/sony/sonyflake"
	"time"
)

// NewSonyflakeID generates a distributed unique identifier based on
// time & machine's private IP
func NewSonyflakeID() uint64 {
	id, err := sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime:      time.Time{},
		MachineID:      nil,
		CheckMachineID: nil,
	}).NextID()

	if err != nil {
		return 0
	}

	return id
}
