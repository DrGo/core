package str

// InsertToStringSlice inserts the value into the slice at the specified index,
// which must be in range.
// The slice must have room for the new element.
//from https://blog.golang.org/slices
func InsertToStringSlice(slice []string, index int, value string) []string {
	if len(slice) == 0 { //inserting into an empty slice; just append it ignoring index
		slice = append(slice, value)
		return slice
	}
	// Grow the slice by one element.
	slice = slice[0 : len(slice)+1]
	// Use copy to move the upper part of the slice out of the way and open a hole.
	copy(slice[index+1:], slice[index:])
	// Store the new value.
	slice[index] = value
	// Return the result.
	return slice
}
