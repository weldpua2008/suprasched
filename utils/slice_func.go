package utils
import (
    "sort"
)
// ContainsIntsSort .
func ContainsIntsSort(in []int, search int) bool {
    if !sort.IntsAreSorted(in){
        sort.Ints(in)
    }
    return sort.SearchInts(in, search) < len(in)
}
// ContainsIntsRange.
func ContainsIntsRange(in []int, search int) bool {
    for _, i := range in {
        if i == search {
            return true
        }
    }
    return false
}

func ContainsInts(in []int, search int) bool {
    return ContainsIntsRange(in, search)
}
