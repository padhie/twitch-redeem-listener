package output

type Device interface {
	Toggle(ToggleInput) error
}

type ToggleInput struct {
	User       string
	RedeemName string
}
