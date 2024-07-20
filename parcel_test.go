package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	added, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, added)

	// get
	get, err := store.Get(added)
	parcel.Number = get.Number
	require.NoError(t, err)
	require.NotEmpty(t, parcel, get)
	require.Equal(t, parcel.Number, get.Number)

	// delete
	err = store.Delete(added)
	require.NoError(t, err)
	_, err = store.Get(added)
	require.Error(t, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add

	added, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, added)

	// set address
	newAddress := "new test address"
	err = store.SetAddress(added, newAddress)
	require.NoError(t, err)

	// check
	get, err := store.Get(added)
	require.NoError(t, err)
	require.Equal(t, newAddress, get.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	added, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, added)

	// set status
	newStatus := parcel.Status
	err = store.SetStatus(added, newStatus)
	require.NoError(t, err)

	// check
	get, err := store.Get(added)
	require.NoError(t, err)
	require.Equal(t, newStatus, get.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	for i := range parcels {
		parcels[i].Client = client
	}

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotEmpty(t, id)

		parcels[i].Number = id

		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Equal(t, len(parcels), len(storedParcels))
	require.NotEmpty(t, storedParcels)

	// check
	for _, parcel := range storedParcels {
		expectedParcel, ok := parcelMap[parcel.Number]
		require.True(t, ok, "unexpected parcel found")
		require.Equal(t, expectedParcel, parcel)
	}
}
