package slicer

// RestartSlice trim comparable slice , return slice start with beginID
func RestartSliceInt[T comparable](inputs []T, beginID T) []T {

	if len(inputs) == 0 {
		return nil
	}

	var defaultValue T
	if beginID == defaultValue {
		return inputs
	}

	var indexID int
	for i, id := range inputs {
		if id == beginID {
			indexID = i
			break
		}
	}
	return inputs[indexID:]
}

type Identifier interface {
	Identity() uint64
}

// RestartSlice trim model.ItemAveragePrice slice , return slice start with beginID
func RestartSliceStruct[T Identifier](inputs []T, beginID uint64) []T {
	if len(inputs) == 0 {
		return nil
	}

	if beginID == 0 {
		return inputs
	}

	var indexID int
	for i, v := range inputs {
		if v.Identity() == beginID {
			indexID = i
			break
		}
	}
	return inputs[indexID:]
}
