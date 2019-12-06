package lustre

import (
	"fmt"
)

type OSTLayout struct {
	mirrorState uint32
	poolName    string
	ostIndice   []uint32
}

func (ol *OSTLayout) PoolName() string {
	return ol.poolName
}

func (ol *OSTLayout) MirrorState() uint32 {
	return ol.mirrorState
}

func (ol *OSTLayout) OstIndice() []uint32 {
	return ol.ostIndice
}

func (ol *OSTLayout) Dump() {
	fmt.Println(ol.mirrorState, ol.poolName)
	fmt.Println(ol.ostIndice)
}
