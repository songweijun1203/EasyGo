package xutils

import (
	"sync/atomic"
)

// TimeStampType 定义时间戳的底层类型
type TimeStampType uint64

// 包级别的原子变量，用于生成全局单调递增的时间戳
var globalTimeStamp uint64

type XTimestamp struct {
	modifiedTime TimeStampType
}

// NewXTimestamp 创建一个新的时间戳对象，并自动标记为已修改
func NewXTimestamp() *XTimestamp {
	ts := &XTimestamp{}
	ts.Modified()
	return ts
}

// Modified 更新时间戳为全局最新的递增值
func (ts *XTimestamp) Modified() {
	// 原子自增，返回自增后的新值
	newVal := atomic.AddUint64(&globalTimeStamp, 1)
	ts.modifiedTime = TimeStampType(newVal)
}

// GetMTime 获取当前的修改时间戳
func (ts *XTimestamp) GetMTime() TimeStampType {
	return ts.modifiedTime
}

// IsNewerThan 判断当前时间戳是否比 rhs 更新 (对应 C++ 的 operator>)
func (ts *XTimestamp) IsNewerThan(rhs *XTimestamp) bool {
	return ts.modifiedTime > rhs.modifiedTime
}

// IsOlderThan 判断当前时间戳是否比 rhs 更旧 (对应 C++ 的 operator<)
func (ts *XTimestamp) IsOlderThan(rhs *XTimestamp) bool {
	return ts.modifiedTime < rhs.modifiedTime
}
