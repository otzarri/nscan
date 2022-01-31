package conv

import (
	"reflect"
	"testing"
)

func TestStringToInt(t *testing.T) {
	tables := []struct {
		strNum string
		intNum int
	}{
		{"-6789", -6789},
		{"-345", -345},
		{"-12", -12},
		{"0", 0},
		{"12", 12},
		{"345", 345},
		{"6789", 6789},
	}

	for _, table := range tables {
		intNum := StringToInt(table.strNum)

		if intNum != table.intNum {
			t.Errorf("Error converting value from string (%v) to int (%v).\n",
				table.strNum, table.intNum)
		}
	}
}

func TestSliceIntToUint16(t *testing.T) {
	tables := []struct {
		intSlice    []int
		uint16Slice []uint16
	}{
		{[]int{20, 21, 22}, []uint16{uint16(20), uint16(21), uint16(22)}},
		{[]int{80, 443}, []uint16{uint16(80), uint16(443)}},
	}

	for _, table := range tables {
		uint16Slice := SliceIntToUint16(table.intSlice)

		if !reflect.DeepEqual(uint16Slice, table.uint16Slice) {
			t.Errorf("Error: Converted slice (%v) is different from the desired one(%v)\n",
				uint16Slice, table.uint16Slice)
		}
	}
}

func TestSliceStringToUint16(t *testing.T) {
	tables := []struct {
		stringSlice []string
		uint16Slice []uint16
	}{
		{[]string{"20", "21", "22"}, []uint16{uint16(20), uint16(21), uint16(22)}},
		{[]string{"80", "443"}, []uint16{uint16(80), uint16(443)}},
	}

	for _, table := range tables {
		uint16Slice := SliceStringToUint16(table.stringSlice)

		if !reflect.DeepEqual(uint16Slice, table.uint16Slice) {
			t.Errorf("Error: Converted slice (%v) is different from the desired one(%v)\n",
				uint16Slice, table.uint16Slice)
		}
	}
}
