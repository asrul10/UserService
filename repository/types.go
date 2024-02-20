// This file contains types that are used in the repository layer.
package repository

type CreateUserInput struct {
	PhoneNumber string
	FullName    string
	Password    string
}

type CreateUserOutput struct {
	UserId int
}

type GetUserByPhoneNumberInput struct {
	PhoneNumber string
}

type GetUserByPhoneNumberOutput struct {
	UserId   int
	Password string
}

type GetUserByIdInput struct {
	UserId int
}

type GetUserByIdOutput struct {
	UserId      int
	FullName    string
	PhoneNumber string
}

type UpdateUserByIdInput struct {
	UserId      int
	FullName    string
	PhoneNumber string
}

type UpdateUserByIdOutput struct {
	UserId      int
	FullName    string
	PhoneNumber string
}

type SuccessLoginCountInput struct {
	UserId int
}

type SuccessLoginCountOutput struct {
	UserId int
}

type IsPhoneNumberChangedInput struct {
	UserId      int
	PhoneNumber string
}

type IsPhoneNumberChangedOutput struct {
	IsChanged bool
}
