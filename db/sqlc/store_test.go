package sqlc

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T){
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println(">> before:", account1.Balance, account2.Balance)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	txResults := make(chan TransferTxResult)

	for i:=0; i<n; i++ {
		// txName := fmt.Sprintf("tx: %d", i+1)
		go func(){
			ctx := context.Background()
			result, err := store.transferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID: account2.ID,
				Amount: amount,
			})

			errs <- err
			txResults <- result
		}()
	}

	for i:=0; i<n; i++ {
		err := <- errs
		require.NoError(t, err)

		result := <- txResults
		require.NotEmpty(t, result)
		require.NotEmpty(t, result.Transfer)

		require.Equal(t, result.Transfer.FromAccountID, account1.ID)
		require.Equal(t, result.Transfer.ToAccountID, account2.ID)
		_, err = store.GetTransfer(context.Background(), result.Transfer.ID)
		require.Empty(t, err)

		require.Equal(t, result.FromEntry.Amount, int64(-10))
		require.Equal(t, result.ToEntry.Amount, int64(10))
		require.Equal(t, result.FromEntry.AccountID, account1.ID)
		require.Equal(t, result.ToEntry.AccountID, account2.ID)
		_, err = store.GetEntry(context.Background(), result.FromEntry.ID)
		require.Empty(t, err)

		_, err = store.GetEntry(context.Background(), result.ToEntry.ID)
		require.Empty(t, err)

		fromAcc := result.FromAccount;
		toAcc := result.ToAccount

		require.NotEmpty(t, fromAcc)
		require.NotEmpty(t, toAcc)

		require.Equal(t, fromAcc.ID, account1.ID)
		require.Equal(t, toAcc.ID, account2.ID)

		fmt.Println(">> tx:", fromAcc.Balance, toAcc.Balance)


		diff1 := account1.Balance - fromAcc.Balance
		diff2 := toAcc.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff2%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
	}
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.Empty(t, err)
	require.Equal(t, updatedAccount1.Balance, account1.Balance - int64(n) * amount)
	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)

	fmt.Println(">> after:", account1.Balance, account2.Balance)

	require.Empty(t, err)
	require.Equal(t, updatedAccount2.Balance, account2.Balance + int64(n) * amount)

}