package conv

import "strconv"

// stringToInt gets an "string" value and returns it converted to "int"
func StringToInt(input string) int {
	output, err := strconv.Atoi(input)
	if err != nil {
		panic(err)
	}
	return output
}

// sliceIntToUint16 gets an "int" slice and returns it converted to "uint16"
func SliceIntToUint16(intSlice []int) []uint16 {
	uint16Slice := []uint16{}
	for _, value := range intSlice {
		uint16Slice = append(uint16Slice, uint16(value))
	}
	return uint16Slice
}

// sliceStringToUint16 gets an "string" slice and returns it converted to "uint16"
func SliceStringToUint16(stringSlice []string) []uint16 {
	uint16Slice := []uint16{}
	for _, strValue := range stringSlice {
		intValue := StringToInt(strValue)
		uint16Slice = append(uint16Slice, uint16(intValue))
	}
	return uint16Slice
}
