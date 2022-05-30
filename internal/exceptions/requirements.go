package exceptions

type ErrMissingRequirements struct {
	MissingRequirements []string
}

func (ErrMissingRequirements) Error() string {
	return "ErrMissingRequirements"
}
