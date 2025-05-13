package parallel

// Map manipulates a slice and transforms it to a slice of another type.
// `iteratee` is call in parallel. Result keep the same order.
func Map[T any, R any](collection []T, iteratee func(item T, index int) R, opts ...Option) []R {
	result := make([]R, len(collection))

	eg := getErrGroup(opts)

	for i, item := range collection {
		_i, _item := i, item
		eg.Go(func() error {
			iteratee(_item, _i)
			return nil
		})
	}

	eg.Wait()

	return result
}

// ForEach iterates over elements of collection and invokes iteratee for each element.
// `iteratee` is call in parallel.
func ForEach[T any](collection []T, iteratee func(item T, index int), opts ...Option) {
	eg := getErrGroup(opts)

	for i, item := range collection {
		_i, _item := i, item
		eg.Go(func() error {
			iteratee(_item, _i)
			return nil
		})
	}
	eg.Wait()
}

// Times invokes the iteratee n times, returning an array of the results of each invocation.
// The iteratee is invoked with index as argument.
// `iteratee` is call in parallel.
func Times[T any](count int, iteratee func(index int) T, opts ...Option) []T {
	result := make([]T, count)

	eg := getErrGroup(opts)

	for i := 0; i < count; i++ {
		_i := i
		eg.Go(func() error {
			result[_i] = iteratee(_i)
			return nil
		})
	}

	eg.Wait()

	return result
}

// GroupBy returns an object composed of keys generated from the results of running each element of collection through iteratee.
// The order of grouped values is determined by the order they occur in the collection.
// `iteratee` is call in parallel.
func GroupBy[T any, U comparable, Slice ~[]T](collection Slice, iteratee func(item T) U, opts ...Option) map[U]Slice {
	result := map[U]Slice{}

	keys := Map(collection, func(item T, _ int) U {
		return iteratee(item)
	}, opts...)

	for i, item := range collection {
		result[keys[i]] = append(result[keys[i]], item)
	}

	return result
}

// PartitionBy returns an array of elements split into groups. The order of grouped values is
// determined by the order they occur in collection. The grouping is generated from the results
// of running each element of collection through iteratee.
// The order of groups is determined by their first appearance in the collection.
// `iteratee` is call in parallel.
func PartitionBy[T any, K comparable, Slice ~[]T](collection Slice, iteratee func(item T) K, opts ...Option) []Slice {
	result := []Slice{}
	seen := map[K]int{}

	keys := Map(collection, func(item T, _ int) K {
		return iteratee(item)
	}, opts...)

	for i, item := range collection {
		if resultIndex, ok := seen[keys[i]]; ok {
			result[resultIndex] = append(result[resultIndex], item)
		} else {
			resultIndex = len(result)
			seen[keys[i]] = resultIndex
			result = append(result, Slice{item})
		}
	}

	return result
}

// ForEachE iterates over elements of collection and invokes iteratee for each element.
// `iteratee` is call in parallel.
func ForEachE[T any](collection []T, iteratee func(item T, index int) error, opts ...Option) error {
	eg := getErrGroup(opts)

	for i, item := range collection {
		_i, _item := i, item
		eg.Go(func() error {
			return iteratee(_item, _i)
		})
	}
	return eg.Wait()
}
