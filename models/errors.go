package models

type PickupAlreadyExistsError struct{}

func (p PickupAlreadyExistsError) Error() string {
	return "pickup already exists"
}
