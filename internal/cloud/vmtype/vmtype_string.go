// Code generated by "stringer -type=VMType"; DO NOT EDIT.

package vmtype

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Unknown-0]
	_ = x[AzureCVM-1]
	_ = x[AzureTrustedLaunch-2]
}

const _VMType_name = "UnknownAzureCVMAzureTrustedLaunch"

var _VMType_index = [...]uint8{0, 7, 15, 33}

func (i VMType) String() string {
	if i >= VMType(len(_VMType_index)-1) {
		return "VMType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _VMType_name[_VMType_index[i]:_VMType_index[i+1]]
}
