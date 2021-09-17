package blogServer

type InputError struct{ ErrMsg string }

func (err InputError) Error() string { return err.ErrMsg }
func NewInputError(msg string) error { return InputError{msg} }
