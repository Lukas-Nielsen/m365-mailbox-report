package main

import (
	"fmt"
	"strconv"
)

func formatBytes(s string) string {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return ""
	}
	const unit = 1024
	const units = "KMGTPE"

	if val < unit {
		return fmt.Sprintf("%d B", val)
	}
	exp, div := 0, int64(unit)
	for val/div >= unit && exp < len(units)-1 {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(val)/float64(div), units[exp])
}
