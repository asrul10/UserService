// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	CreateUser(
		ctx context.Context,
		input CreateUserInput,
	) (output CreateUserOutput, err error)
	GetUserByPhoneNumber(
		ctx context.Context,
		input GetUserByPhoneNumberInput,
	) (output GetUserByPhoneNumberOutput, err error)
	GetUserById(
		ctx context.Context,
		input GetUserByIdInput,
	) (output GetUserByIdOutput, err error)
	UpdateUserById(
		ctx context.Context,
		input UpdateUserByIdInput,
	) (output UpdateUserByIdOutput, err error)
	SuccessLoginCount(
		ctx context.Context,
		input SuccessLoginCountInput,
	) (output SuccessLoginCountOutput, err error)
	IsPhoneNumberChanged(
		ctx context.Context,
		input IsPhoneNumberChangedInput,
	) (output IsPhoneNumberChangedOutput, err error)
}
