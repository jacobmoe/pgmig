package errors

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
