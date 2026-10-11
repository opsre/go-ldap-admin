package tools

import (
	"fmt"
	"strconv"
	"strings"
)

// Compare string slices
func ArrStrCmp(src []string, dest []string) ([]string, []string) {
	msrc := make(map[string]byte) //Build an index from the source slice
	mall := make(map[string]byte) //Build an index from all source and destination elements
	var set []string              //Intersection
	//1. Build a map from the source slice
	for _, v := range src {
		msrc[v] = 0
		mall[v] = 0
	}
	//2. Elements that cannot be inserted from the destination slice are duplicates; the remaining elements form the union
	for _, v := range dest {
		l := len(mall)
		mall[v] = 1
		if l != len(mall) {
			continue
		} else {
			set = append(set, v)
		}
	}
	//3. For each intersection element, find and remove it from the union; the remaining elements are the difference (union minus intersection)
	for _, v := range set {
		delete(mall, v)
	}
	//4. mall is now the difference: elements found in the source were removed, and the rest are newly added
	var added, deleted []string
	for v := range mall {
		_, exist := msrc[v]
		if exist {
			deleted = append(deleted, v)
		} else {
			added = append(added, v)
		}
	}
	return added, deleted
}

// Compare uint slices
func ArrUintCmp(src []uint, dest []uint) ([]uint, []uint) {
	msrc := make(map[uint]byte) //Build an index from the source slice
	mall := make(map[uint]byte) //Build an index from all source and destination elements
	var set []uint              //Intersection
	//1. Build a map from the source slice
	for _, v := range src {
		msrc[v] = 0
		mall[v] = 0
	}
	//2. Elements that cannot be inserted from the destination slice are duplicates; the remaining elements form the union
	for _, v := range dest {
		l := len(mall)
		mall[v] = 1
		if l != len(mall) {
			continue
		} else {
			set = append(set, v)
		}
	}
	//3. For each intersection element, find and remove it from the union; the remaining elements are the difference (union minus intersection)
	for _, v := range set {
		delete(mall, v)
	}
	//4. mall is the difference; elements found in the source were removed
	var added, deleted []uint
	for v := range mall {
		_, exist := msrc[v]
		if exist {
			deleted = append(deleted, v)
		} else {
			added = append(added, v)
		}
	}
	return added, deleted
}

// Convert a string slice to a uint slice
func SliceToString(src []uint, delim string) string {
	return strings.Trim(strings.ReplaceAll(fmt.Sprint(src), " ", delim), "[]")
}

// Convert a string slice to a uint slice
func StringToSlice(src string, delim string) []uint {
	var dest []uint
	if src == "" {
		return dest
	}
	strs := strings.SplitSeq(src, delim)
	for v := range strs {
		t, _ := strconv.Atoi(v)
		dest = append(dest, uint(t))
	}
	return dest
}
