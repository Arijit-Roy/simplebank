package sqlc

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"simplebank/util"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account{
		args := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomInt(1000, 1000000),
		Currency: util.RandomCurrency(),
	}

	acc, err := testQueries.CreateAccount(context.Background(), args)
		require.NoError(t, err)
	require.NotEmpty(t, acc)
	require.Equal(t, acc.Balance, args.Balance)
	require.Equal(t, acc.Owner, args.Owner)
	require.Equal(t, acc.Currency, args.Currency)

	require.NotZero(t, acc.ID)
	require.NotZero(t, acc.CreatedAt)
	return acc
}

func TestCreateAccount(t *testing.T){
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T){
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Balance, account2.Balance)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

}

func TestUpdateAccount(t *testing.T){
	account1 := createRandomAccount(t)

	args := UpdateAccountParams{
		ID: account1.ID,
		Balance: util.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, args.Balance, account2.Balance)

}

func TestDeleteAccount(t *testing.T){
	account1 := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err1 := testQueries.GetAccount(context.Background(), account1.ID)
	require.Empty(t, account2)
	require.Error(t, err1)
	require.EqualError(t, err1, sql.ErrNoRows.Error())
}

func TestListAccount(t *testing.T){
	for i:=0; i<11; i++ {
		createRandomAccount(t)
	}

	args := ListAccountParams{
		Limit: 5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccount(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}