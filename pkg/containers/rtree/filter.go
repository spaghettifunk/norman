package rtreego

import "github.com/twpayne/go-geom"

// Filter is an interface for filtering leaves during search. The parameters
// should be treated as read-only. If refuse is true, the current entry will
// not be added to the result set. If abort is true, the search is aborted and
// the current result set will be returned.
type Filter func(results []geom.T, object geom.T) (refuse, abort bool)

// ApplyFilters applies the given filters and returns whether the entry is
// refused and/or the search should be aborted. If a filter refuses an entry,
// the following filters are not applied for the entry. If a filter aborts, the
// search terminates without further applying any filter.
func applyFilters(results []geom.T, object geom.T, filters []Filter) (bool, bool) {
	for _, filter := range filters {
		refuse, abort := filter(results, object)
		if refuse || abort {
			return refuse, abort
		}
	}
	return false, false
}

// LimitFilter checks if the results have reached the limit size and aborts if so.
func LimitFilter(limit int) Filter {
	return func(results []geom.T, object geom.T) (refuse, abort bool) {
		if len(results) >= limit {
			return true, true
		}

		return false, false
	}
}
