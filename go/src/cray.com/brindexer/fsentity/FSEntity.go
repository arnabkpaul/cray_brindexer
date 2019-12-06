package fsentity

import ()

type FSEntity interface {
	CommitSqls() string
	PathMd5() string
	RPath() string
	Name() string
}
