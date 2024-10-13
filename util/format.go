package util

import (
	"fmt"
	"time"
)

// FmtBytes converts the numeric byte size value to the appropriate magnitude
// size in KB, MB, GB, TB, PB, or EB.
func FmtBytes(value interface{}) string {
	var size float64

	switch value.(type) {
	case int64:
		size = float64(value.(int64))
	case float64:
		size = value.(float64)
	}

	unit := float64(1024)
	BYTE := unit
	KBYTE := BYTE * unit
	MBYTE := KBYTE * unit
	GBYTE := MBYTE * unit
	TBYTE := GBYTE * unit
	PBYTE := TBYTE * unit

	switch {
	case size < BYTE:
		return fmt.Sprintf("%0.f B", size)
	case size < KBYTE:
		return fmt.Sprintf("%.1f KB", size/BYTE)
	case size < MBYTE:
		return fmt.Sprintf("%.1f MB", size/KBYTE)
	case size < GBYTE:
		return fmt.Sprintf("%.1f GB", size/MBYTE)
	case size < TBYTE:
		return fmt.Sprintf("%.1f TB", size/GBYTE)
	case size < PBYTE:
		return fmt.Sprintf("%.1f PB", size/TBYTE)
	default:
		return fmt.Sprintf("%0.f B", size)
	}

}

// FmtTime shows a human-readable time based on the timestamp
func FmtTime(t int64) string {
	return time.Unix(t/1000, 0).Format("03:04 PM Jan 2, 2006") + " (UTC)"
}
