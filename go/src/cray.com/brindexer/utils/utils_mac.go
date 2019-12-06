// +build darwin

package utils

import (
	"syscall"
)

func TimesFromStat_s(stat *syscall.Stat_t) (int64, int64, int64) {
	return stat.Mtimespec.Nano(), stat.Atimespec.Nano(), stat.Ctimespec.Nano()
}
