
package dataset

type Dataset[T any] struct {
	data []T
}

func New[T any](data []T) *Dataset[T] {
	return &Dataset[T]{data: data}
}

func (d *Dataset[T]) Map[U any](fn func(T) U) *Dataset[U] {
	out := make([]U, len(d.data))
	for i, v := range d.data {
		out[i] = fn(v)
	}
	return New(out)
}

func (d *Dataset[T]) Filter(fn func(T) bool) *Dataset[T] {
	var out []T
	for _, v := range d.data {
		if fn(v) {
			out = append(out, v)
		}
	}
	return New(out)
}

func (d *Dataset[T]) Batch(batchSize int) *Dataset[[]T] {
	var batches [][]T
	for i := 0; i < len(d.data); i += batchSize {
		end := i + batchSize
		if end > len(d.data) {
			end = len(d.data)
		}
		batches = append(batches, d.data[i:end])
	}
	return New(batches)
}

func (d *Dataset[T]) Repeat(n int) *Dataset[T] {
	var out []T
	for i := 0; i < n; i++ {
		out = append(out, d.data...)
	}
	return New(out)
}

func (d *Dataset[T]) ToSlice() []T {
	return d.data
}