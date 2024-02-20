package repository

import (
	"context"
)

func (r *Repository) CreateUser(ctx context.Context, input CreateUserInput) (output CreateUserOutput, err error) {
	tx, err := r.Db.Begin()
	if err != nil {
		return
	}
	defer tx.Commit()

	err = r.Db.QueryRowContext(
		ctx,
		"INSERT INTO users (phone_number, full_name, password) VALUES ($1, $2, $3) RETURNING id",
		input.PhoneNumber,
		input.FullName,
		input.Password,
	).Scan(&output.UserId)
	if err != nil {
		tx.Rollback()
		return
	}

	return
}

func (r *Repository) GetUserByPhoneNumber(ctx context.Context, input GetUserByPhoneNumberInput) (output GetUserByPhoneNumberOutput, err error) {
	err = r.Db.QueryRowContext(
		ctx,
		"SELECT id, password FROM users WHERE phone_number = $1",
		input.PhoneNumber,
	).Scan(&output.UserId, &output.Password)
	if err != nil {
		return
	}

	return
}

func (r *Repository) GetUserById(ctx context.Context, input GetUserByIdInput) (output GetUserByIdOutput, err error) {
	err = r.Db.QueryRowContext(
		ctx,
		"SELECT id, full_name, phone_number FROM users WHERE id = $1",
		input.UserId,
	).Scan(&output.UserId, &output.FullName, &output.PhoneNumber)
	if err != nil {
		return
	}

	return
}

func (r *Repository) UpdateUserById(ctx context.Context, input UpdateUserByIdInput) (output UpdateUserByIdOutput, err error) {
	tx, err := r.Db.Begin()
	if err != nil {
		return
	}
	defer tx.Commit()

	_, err = r.Db.ExecContext(
		ctx,
		"UPDATE users SET full_name = $1, phone_number = $2 WHERE id = $3",
		input.FullName,
		input.PhoneNumber,
		input.UserId,
	)
	if err != nil {
		tx.Rollback()
		return
	}

	output.UserId = input.UserId
	output.FullName = input.FullName
	output.PhoneNumber = input.PhoneNumber

	return
}

func (r *Repository) SuccessLoginCount(ctx context.Context, input SuccessLoginCountInput) (output SuccessLoginCountOutput, err error) {
	tx, err := r.Db.Begin()
	if err != nil {
		return
	}
	defer tx.Commit()

	_, err = r.Db.ExecContext(
		ctx,
		"UPDATE users SET success_login_count = success_login_count + 1 WHERE id = $1",
		input.UserId,
	)
	if err != nil {
		tx.Rollback()
		return
	}

	output.UserId = input.UserId

	return
}

func (r *Repository) IsPhoneNumberChanged(ctx context.Context, input IsPhoneNumberChangedInput) (output IsPhoneNumberChangedOutput, err error) {
	var phoneNumber string
	err = r.Db.QueryRowContext(
		ctx,
		"SELECT phone_number FROM users WHERE id = $1",
		input.UserId,
	).Scan(&phoneNumber)
	if err != nil {
		return
	}

	output.IsChanged = phoneNumber != input.PhoneNumber

	return
}
