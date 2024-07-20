package main

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	rows, err := s.db.Exec("insert into parcel (client, status, address, created_at) values (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))

	if err != nil {
		return 0, err
	}

	req, err := rows.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(req), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	p := Parcel{}
	rows := s.db.QueryRow("select number, client, status, address, created_at from parcel where number = :number",
		sql.Named("number", number))
	err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := s.db.Query("select number, client, status, address, created_at from parcel where client = :client",
		sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Parcel
	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("update parcel set status = :status where number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	_, err := s.db.Exec("update parcel set address = :address where number = :number",
		sql.Named("address", address),
		sql.Named("number", number))
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	_, err := s.db.Exec("delete from parcel where number = :number", sql.Named("number", number))
	if err != nil {
		return err
	}
	return nil
}
