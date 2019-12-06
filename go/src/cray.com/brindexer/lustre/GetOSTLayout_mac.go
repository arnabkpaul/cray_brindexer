// +build darwin

package lustre

import ()

//Just to return some dummy on MAC, so we don't break build
func GetLayout(path string) *OSTLayout {

	osts := make([]uint32, 0, 64)
	layout := OSTLayout{0, "", osts}
	return &layout

}
