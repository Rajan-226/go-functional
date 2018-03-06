// Package functional implements functional operations for slices of go primitives
package functional

// StringSliceFunctor implements functional operations for type []string
type StringSliceFunctor struct {
	slice []string
}

// LiftStringSlice creates an StringSliceFunctor from a []string
func LiftStringSlice(slice []string) StringSliceFunctor {
	return StringSliceFunctor{slice: slice}
}

// Collect returns the underlying []string
func (f StringSliceFunctor) Collect() []string {
	return f.slice
}

// Map returns a new StringSliceFunctor who's underlying slice is the result of
// applying the input operation to each of its members.
func (f StringSliceFunctor) Map(op func(string) string) StringSliceFunctor {
	mapped := make([]string, 0, len(f.slice))
	for _, i := range f.slice {
		mapped = append(mapped, op(i))
	}
	return LiftStringSlice(mapped)
}

// Filter returns a new StringSliceFunctor who's underlying slice has had
// members exluded that do not satisfy the input filter.
func (f StringSliceFunctor) Filter(op func(string) bool) StringSliceFunctor {
	var filtered []string
	for _, i := range f.slice {
		if op(i) {
			filtered = append(filtered, i)
		}
	}
	return LiftStringSlice(filtered)
}

// Fold applies its input operation to the initial input value and the first
// member of the underlying slice. It successively applies the input operation
// to the result of the previous and the next value in the underlying slice. It
// returns the final value successful operations. If the underlying slice is
// empty then Fold returns the initial input value.
func (f StringSliceFunctor) Fold(initial string, op func(string, string) string) string {
	for _, i := range f.slice {
		initial = op(initial, i)
	}
	return initial
}

// Take returns a new StringSliceFunctor who's underlying slice has had all
// members after the nth dropped. If n is larger than the length of the
// underlying slice, Take is a no-op.
func (f StringSliceFunctor) Take(n int) StringSliceFunctor {
	if n > len(f.slice) {
		return f
	}
	return LiftStringSlice(f.slice[0:n])
}

// Drop returns a new StringSliceFunctor who's underlying slice has had the
// first n members dropped. If n is larger than the length of the underlying
// slice, Drop returns an empty StringSliceFunctor.
func (f StringSliceFunctor) Drop(n int) StringSliceFunctor {
	if n > len(f.slice) {
		return LiftStringSlice([]string{})
	}
	return LiftStringSlice(f.slice[n:len(f.slice)])
}

// WithErrs creates an StringSliceErrFunctor from a StringSliceFunctor.
func (f StringSliceFunctor) WithErrs() StringSliceErrFunctor {
	return StringSliceErrFunctor{slice: f.slice}
}

// StringSliceErrFunctor behaves like StringSliceFunctor except that operations
// performed over the underlying slice are allowed to return errors. Should
// an error occur then the StringSliceErrFunctor's future operations do nothing
// except that Collect will return the error that occurred.
type StringSliceErrFunctor struct {
	slice []string
	err   error
}

// Collect returns the underlying []string or an error if one has occurred.
func (f StringSliceErrFunctor) Collect() ([]string, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.slice, nil
}

// Map returns a new StringSliceErrFunctor who's underlying slice is the result
// of applying the input operation to each of its members. Should an error
// occur, the underlying slice is lost and subsequent Collect calls will always
// return the error.
func (f StringSliceErrFunctor) Map(op func(string) (string, error)) StringSliceErrFunctor {
	if f.err != nil {
		return f
	}

	mapped := make([]string, len(f.slice))
	for i := range f.slice {
		new, err := op(f.slice[i])
		if err != nil {
			return StringSliceErrFunctor{err: err}
		}
		mapped[i] = new
	}
	return LiftStringSlice(mapped).WithErrs()
}