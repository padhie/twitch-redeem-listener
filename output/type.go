package output

type Device interface {
	Toggle() error
}
