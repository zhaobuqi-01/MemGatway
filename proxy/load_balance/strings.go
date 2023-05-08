package load_balance

import "strings"

func InStringSlice(slice []string, str string) bool {
	for _, item := range slice {
		if str == item {
			return true
		}
	}
	return false
}

// SplitStringByComma takes a string and splits it into a slice of strings using a comma as the separator.
// The input is a single string, and the output is a slice of strings.
//
// Usage example:
//
//	input := "apple,banana,orange"
//	output := SplitStringByComma(input)
//	fmt.Println(output) // Output: ["apple" "banana" "orange"]
func SplitStringByComma(data string) []string {
	return strings.Split(data, ",")
}
