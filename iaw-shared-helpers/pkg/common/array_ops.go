// ------------------------------------------------------ {COPYRIGHT-TOP} ---
// IBM Confidential
// OCO Source Materials
// 5900-AEO
//
// Copyright IBM Corp. 2021
//
// The source code for this program is not published or otherwise
// divested of its trade secrets, irrespective of what has been
// deposited with the U.S. Copyright Office.
// ------------------------------------------------------ {COPYRIGHT-END} ---
package common

import "sort"

// CombineStringSlices combines a set of []string into a single slice
// this slice is then ordered and deduplicated
func CombineStringSlices(sliceSet ...[]string) []string {
	uniqueValues := map[string]struct{}{}
	for _, slice := range sliceSet {
		for _, val := range slice {
			uniqueValues[val] = struct{}{}
		}
	}
	newSlice := []string{}
	for key := range uniqueValues {
		newSlice = append(newSlice, key)
	}
	sort.Strings(newSlice)
	return newSlice
}
