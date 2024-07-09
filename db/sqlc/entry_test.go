package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func generateEntryRandom(t *testing.T, account Account) Entry {

	input := CreateEntryParams{
		AccountID: account.ID,
		Amount:    account.Balance,
	}

	entry, err := testQueries.CreateEntry(context.Background(), input)
	if err != nil {
		fmt.Printf("%# v\n", err)
	}

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, input.AccountID, entry.AccountID, "input and return account id should be same")
	require.NotEmpty(t, entry.EntriesUuid.String())
	require.NotNil(t, entry.CreatedAt)
	require.NotZero(t, entry.ID)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := generateAccount(t)
	generateEntryRandom(t, account)
}

func TestGetEntry(t *testing.T) {

	account := generateAccount(t)
	createRandEntry := generateEntryRandom(t, account)

	getFromEntry, err := testQueries.GetEntry(context.Background(), createRandEntry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, getFromEntry)
	require.Equal(t, createRandEntry.ID, getFromEntry.ID)
	require.Equal(t, createRandEntry.EntriesUuid.String(), getFromEntry.EntriesUuid.String())
	require.WithinDuration(t, createRandEntry.CreatedAt, getFromEntry.CreatedAt, time.Second)
}

func TestListEntry(t *testing.T) {
	account := generateAccount(t)

	for i := 0; i < 10; i++ {
		generateEntryRandom(t, account)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	// fmt.Printf("%# v\n", entries)
	require.NoError(t, err)
	require.NotEmpty(t, entries)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, account.ID, entry.AccountID)
	}
}
