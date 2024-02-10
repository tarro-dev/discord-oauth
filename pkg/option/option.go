package option

// Option is a type that represents an optional value.
type Option[T any] struct {
	value T
	ok    bool
}

// New returns an Option[T] with the given value.
func New[T any](v *T) Option[T] {
	if v == nil {
		return None[T]()
	}
	return Some[T](*v)
}

// Some returns an Option[T] with the given value.
func Some[T any](v T) Option[T] {
	return Option[T]{
		value: v,
		ok:    true,
	}
}

// None returns an Option[T] with no value.
func None[T any]() Option[T] {
	return Option[T]{
		ok: false,
	}
}

// IsSome returns true if the Option[T] has a value.
func (o Option[T]) IsSome() bool {
	return o.ok
}

// IsNone returns true if the Option[T] has no value.
func (o Option[T]) IsNone() bool {
	return !o.ok
}

// Unwrap returns the value and true of the Option[T] if it has a value.
// If it has no value, it returns undefined and false.
func (o Option[T]) Unwrap() (T, bool) {
	return o.value, o.ok
}

// UnwrapOr returns the value of the Option[T] if it has a value.
// If it has no value, it returns the given default value.
func (o Option[T]) UnwrapOr(def T) T {
	if o.ok {
		return o.value
	}
	return def
}

// UnwrapOrElse returns the value of the Option[T] if it has a value.
// If it has no value, it returns the result of the given function.
func (o Option[T]) UnwrapOrElse(f func() T) T {
	if o.ok {
		return o.value
	}
	return f()
}

// UnwrapDefault returns the value of the Option[T] if it has a value.
// If it has no value, it returns the zero value of T.
func (o Option[T]) UnwrapDefault() T {
	var def T
	return o.UnwrapOr(def)
}

// UnwrapPtr returns a pointer to the value of the Option[T] if it has a value.
func (o Option[T]) UnwrapPtr() *T {
	if o.ok {
		return &o.value
	}
	return nil
}

// Expect returns the value of the Option[T] if it has a value.
// If it has no value, it panics with the given message.
func (o Option[T]) Expect(msg string) T {
	if o.ok {
		return o.value
	}
	panic(msg)
}
