package db

import (
	"context"
	"testing"

	"github.com/amirrmonfared/SimpleMicroService/authentication-service/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Email:     util.RandomEmail(),
		FirstName: util.RandomName(),
		LastName:  util.RandomName(),
		Password:  util.RandomString(10),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, user.LastName)
	require.Equal(t, arg.Password, user.Password)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateRow(t *testing.T) {
	createRandomUser(t)
}
