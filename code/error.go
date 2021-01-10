package code

//EncoderError is error returned by any decoder function
type EncoderError struct {
	Msg string
}

func (e *EncoderError) Error() string { return e.Msg }
