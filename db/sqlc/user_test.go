package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lyb88999/Go-SimpleBank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandString(6))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.RandOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandOwner(),
		Email:          util.RandomEmail(),
	}
	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}
func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testStore.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.CreatedAt, user2.CreatedAt)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.Equal(t, user1.PasswordChangedAt, user2.PasswordChangedAt)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)
	newFullName := util.RandOwner()
	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: pgtype.Text{
			String: newFullName,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.Equal(t, updatedUser.Username, oldUser.Username)
	require.Equal(t, updatedUser.FullName, newFullName)
	require.Equal(t, updatedUser.Email, oldUser.Email)
	require.Equal(t, updatedUser.HashedPassword, oldUser.HashedPassword)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)
	newEmail := util.RandomEmail()
	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: pgtype.Text{
			String: newEmail,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.Equal(t, updatedUser.Email, newEmail)
	require.Equal(t, updatedUser.Username, oldUser.Username)
	require.Equal(t, updatedUser.FullName, oldUser.FullName)
	require.Equal(t, updatedUser.HashedPassword, oldUser.HashedPassword)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := createRandomUser(t)
	newHashPassword, err := util.HashPassword(util.RandString(6))
	require.NoError(t, err)
	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: pgtype.Text{
			String: newHashPassword,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.Equal(t, updatedUser.Email, oldUser.Email)
	require.Equal(t, updatedUser.Username, oldUser.Username)
	require.Equal(t, updatedUser.FullName, oldUser.FullName)
	require.Equal(t, updatedUser.HashedPassword, newHashPassword)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)
	newFullName := util.RandOwner()
	newEmail := util.RandomEmail()
	newHashPassword, err := util.HashPassword(util.RandString(6))
	require.NoError(t, err)

	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: pgtype.Text{
			String: newFullName,
			Valid:  true,
		},
		Email: pgtype.Text{
			String: newEmail,
			Valid:  true,
		},
		HashedPassword: pgtype.Text{
			String: newHashPassword,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.Equal(t, updatedUser.Email, newEmail)
	require.Equal(t, updatedUser.FullName, newFullName)
	require.Equal(t, updatedUser.HashedPassword, newHashPassword)
	require.Equal(t, updatedUser.Username, oldUser.Username)

}
