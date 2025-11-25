package domain

import "errors"

var (
	ErrNoAvailability = errors.New("One of the selected dates is no longer available.")
	
	ErrInvalidHotelID = errors.New("The hotel ID is invalid.")

	ErrEmptyBooking = errors.New("The reservation must contain at least one room.")

	ErrInvalidDateFormat = errors.New("The date format is invalid. Please use YYYY-MM-DD.")
	ErrInvalidDateRange  = errors.New("The check-out date must be after the check-in date.")
	ErrPastDate          = errors.New("Cannot be reserved in the past")

	ErrEmailAlreadyExists = errors.New("The email address is already registered.")
	ErrInvalidCredentials = errors.New("Invalid credentials")
)
