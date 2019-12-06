package lustre

/*

#include <stdint.h>

typedef struct {
	uint32_t   Hsm_compat;
	uint32_t   Hsm_flags;
	uint64_t   Hsm_arch_id;
	uint64_t   Hsm_arch_ver;
}C_hsmattrs;
*/
import "C"
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/pkg/xattr"
)

const HSM_ATTR = "trusted.hsm"

type HSMAttrs struct {
	HsmCompat  uint32
	HsmFlags   uint32
	HsmArchId  uint64
	HsmArchVer uint64
}

func (ha *HSMAttrs) Dump() {
	fmt.Printf("0x%02x, 0x%02x, 0x%02x ,0x%02x", ha.HsmCompat, ha.HsmFlags, ha.HsmArchId, ha.HsmArchVer)
}

func GetHsmAttrs(path string) *HSMAttrs {
	attrs := HSMAttrs{0, 0, 0, 0}
	data, err := xattr.LGet(path, HSM_ATTR)

	if err != nil {
		//fmt.Println("Failed to get HSM attributes")
		//Return a default one
		return &attrs
	}
	buf := bytes.NewReader(data)
	rerr := binary.Read(buf, binary.BigEndian, &attrs)
	if rerr != nil {
		fmt.Println("Failed to read attrs:", rerr)
		return &attrs
	}

	return &attrs

}
