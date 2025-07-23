package field

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

// ID 唯一标识
type ID int64

func NewID(id int64) ID {
	return ID(id)
}

func (id ID) Value() int64 {
	return int64(id)
}

// ExtractTimestamp 从ID中提取时间戳
func (id ID) ExtractTimestamp() time.Time {
	ms := (int64(id) >> timeShift) + epoch
	return time.Unix(ms/1000, (ms%1000)*1e6)
}

// ExtractMachineID 从ID中提取机器ID
func (id ID) ExtractMachineID() int64 {
	return (int64(id) >> machineShift) & 0x3FF // 10位机器ID
}

// ExtractSequence 从ID中提取序列号
func (id ID) ExtractSequence() int64 {
	return int64(id) & 0xFFF // 12位序列号
}

// Snowflake 分布式ID生成器
// | 1位符号 | 41位时间戳 | 10位机器ID | 12位序列号 |
// |   0    |  时间差值  |   机器编号  |   序号    |
type Snowflake struct {
	machineID int64
	lastStamp int64
	sequence  int64
	mu        sync.Mutex
}

const (
	timeShift    = 22            // 时间戳左移22位(10位机器ID + 12位序列号)，可使用约69年，2^41秒
	machineShift = 12            // 机器ID左移12位(12位序列号)，支持1024台机器同时工作
	maxSequence  = 4095          // 序列号最大值(2^12 - 1)，每台机器每秒生成的最大id数
	epoch        = 1767196800000 // 自定义纪元起始 2026-01-01 00:00:00(毫秒)
)

func NewSnowflake(machines int) *Snowflake {
	s := &Snowflake{}
	s.machineID = s.getMachineID(machines) // 获取机器唯一ID
	return s
}

// getMachineID 获取机器唯一ID的示例函数
func (s *Snowflake) getMachineID(machines int) int64 {
	// 生产环境中应使用稳定的ID分配方案
	//if id := os.Getenv("MACHINE_ID"); id != "" {
	//	i, _ := strconv.Atoi(id)
	//	return int64(i % 1024)
	//}

	// 最大支持1024台机器
	maxMachines := 1 << (timeShift - machineShift)

	// 计算机器ID范围
	machineRange := machines
	if machineRange <= 0 {
		machineRange = 1
	} else if machineRange > maxMachines {
		machineRange = maxMachines
	}

	// 默认使用随机分配 (开发环境)
	u := uuid.New()
	return int64(u.ID() % uint32(machineRange))
}

func (s *Snowflake) Generate() ID {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 获取当前毫秒时间戳
	now := time.Now().UnixNano() / 1e6

	// 如果当前时间小于上次记录时间，说明时钟回拨
	if now < s.lastStamp {
		panic(fmt.Sprintf("Snowflake >>> 时钟回拨异常: %d < %d", now, s.lastStamp))
	}

	// 同一毫秒内的序列号递增
	if now == s.lastStamp {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 { // 当前毫秒序列号用完
			for now <= s.lastStamp {
				now = time.Now().Unix()
			}
		}
	} else {
		s.sequence = 0
	}
	s.lastStamp = now

	// 生成ID: (时间差 << 时间位移) | (机器ID << 机器位移) | 序列号
	return ID(((now - epoch) << timeShift) | (s.machineID << machineShift) | s.sequence)
}
