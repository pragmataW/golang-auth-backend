package models

type UserAlreadyExistsError struct {
	Code int
}

func (e UserAlreadyExistsError) Error() string {
	return "user is already exists"
}

type MailCouldNotSent struct {
	Code int
}

func (e MailCouldNotSent) Error() string {
	return "mail could not sent"
}

type CodeNotExpiredError struct {
	Code int
}

func (e CodeNotExpiredError) Error() string {
	return "code is not expired"
}

type UsernameLengthError struct {
	Code int
}

func (e UsernameLengthError) Error() string {
	return "username length must be greater than 3"
}

type PasswordLengthError struct {
	Code int
}

func (e PasswordLengthError) Error() string {
	return "password length must be greater than 8"
}

type UsernameAlreadyExistsError struct {
	Code int
}

func (e UsernameAlreadyExistsError) Error() string {
	return "username is already exists"
}

type MailAlreadyExistsError struct {
	Code int
}

func (e MailAlreadyExistsError) Error() string {
	return "mail is already exists"
}

type MailEmptyError struct {
	Code int
}

func (e MailEmptyError) Error() string {
	return "mail cannot be empty"
}

type UserNotFoundError struct {
	Code int
}

func (e UserNotFoundError) Error() string {
	return "user not found"
}

type CodeExpiredError struct {
	Code int
}

func (e CodeExpiredError) Error() string {
	return "code expired"
}

type CodeNotFoundError struct{
	Code int
}

func (e CodeNotFoundError) Error() string {
	return "code not found"
}

type CodeResent struct {
	Code int
}

func (e CodeResent) Error() string {
	return "code re-sent"
}

type NotVerifiedError struct{
	Code int
}

func (e NotVerifiedError) Error() string{
	return "not verified"
}