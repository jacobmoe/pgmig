package errors

// ErrUninitialized is an initialization error
type ErrUninitialized struct{}

// NewErrUninitialized returns a new ErrUninitialized
func NewErrUninitialized() *ErrUninitialized {
	return &ErrUninitialized{}
}

// Error returns the error message
func (e *ErrUninitialized) Error() string {
	return "must call Init use"
}

// ErrSQLLoad is an error loading a SQL file
type ErrSQLLoad struct {
	message string
}

// NewSQLLoadErr returns a new ErrSQLLoad
func NewSQLLoadErr(message string) *ErrSQLLoad {
	return &ErrSQLLoad{message: message}
}

// Error returns the error message
func (e *ErrSQLLoad) Error() string {
	return e.message
}
