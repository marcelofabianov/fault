package fault

type Option func(*Error)

func New(message string, opts ...Option) *Error {
	err := &Error{Message: message}
	for _, opt := range opts {
		opt(err)
	}
	return err
}

func Wrap(err error, message string, opts ...Option) *Error {
	opts = append(opts, WithWrappedErr(err))
	return New(message, opts...)
}

func WithWrappedErr(err error) Option {
	return func(e *Error) {
		e.Err = err
	}
}

func WithCode(code Code) Option {
	return func(e *Error) {
		e.Code = code
	}
}

func WithContext(key string, value any) Option {
	return func(e *Error) {
		if e.Context == nil {
			e.Context = make(map[string]any)
		}
		e.Context[key] = value
	}
}

func WithDetails(details ...*Error) Option {
	return func(e *Error) {
		e.Details = append(e.Details, details...)
	}
}

func NewValidationError(err error, message string, context map[string]any) *Error {
	opts := []Option{WithCode(Invalid)}
	if err != nil {
		opts = append(opts, WithWrappedErr(err))
	}
	for k, v := range context {
		opts = append(opts, WithContext(k, v))
	}
	return New(message, opts...)
}

func NewInternalError(err error, context map[string]any) *Error {
	opts := []Option{WithCode(Internal)}
	if err != nil {
		opts = append(opts, WithWrappedErr(err))
	}
	for k, v := range context {
		opts = append(opts, WithContext(k, v))
	}
	return New("An unexpected internal error occurred.", opts...)
}
