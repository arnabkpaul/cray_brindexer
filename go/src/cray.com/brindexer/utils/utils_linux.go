// +build linux

package utils

import (
	"syscall"
)

func TimesFromStat_s(stat *syscall.Stat_t) (int64, int64, int64) {
	return stat.Mtim.Nano(), stat.Atim.Nano(), stat.Ctim.Nano()
}
