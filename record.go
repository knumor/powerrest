package main

import (
	"log"
	"time"
)

type Record struct {
	Id        int        `json:"id"`
	DomainId  int        `json:"domain_id"`
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	Content   string     `json:"content"`
	Ttl       int        `json:"ttl"`
	Priority  *int       `json:"prio"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func FindRecord(id int) (*Record, error) {
	sql := "SELECT id, domain_id, name, type, content, ttl, prio FROM records WHERE id = $1"

	if config.DbType == "mysql" {
		sql = "SELECT id, domain_id, name, type, content, ttl, prio FROM records WHERE id = ?"
	}

	row := db.QueryRow(sql, id)

	r := &Record{}
	err := row.Scan(
		&r.Id,
		&r.DomainId,
		&r.Name,
		&r.Type,
		&r.Content,
		&r.Ttl,
		&r.Priority,
	)

	return r, err
}

func AllRecords() []*Record {
	rows, err := db.Query("SELECT id, domain_id, name, type, content, ttl, prio FROM records")
	if err != nil {
		log.Fatal(err)
	}

	records := []*Record{}

	for rows.Next() {
		r := &Record{}
		err := rows.Scan(
			&r.Id,
			&r.DomainId,
			&r.Name,
			&r.Type,
			&r.Content,
			&r.Ttl,
			&r.Priority,
		)
		if err != nil {
			log.Fatal(err)
		}

		records = append(records, r)
	}

	return records
}

func (r *Record) Create() error {
	sql := "INSERT INTO records (domain_id, name, type, content, ttl, prio, change_date) " +
		"VALUES ($1, $2, $3, $4, $5, $6, extract(epoch from now())::integer)"

	if config.DbType == "mysql" {
		sql = "INSERT INTO records (domain_id, name, type, content, ttl, prio) VALUES (?, ?, ?, ?, ?, ?)"
	}

	_, err := db.Exec(
		sql,
		r.DomainId,
		r.Name,
		r.Type,
		r.Content,
		r.Ttl,
		r.Priority,
	)
	return err
}

func (r *Record) Update() error {
	sql := "UPDATE records SET name=$1, type=$2, content=$3, ttl=$4, prio=$5, " +
		"change_date=extract(epoch from now())::integer WHERE id=$6"

	if config.DbType == "mysql" {
		sql = "UPDATE records SET name=?, type=?, content=?, ttl=?, prio=? WHERE id=?"
	}

	_, err := db.Exec(sql, r.Name, r.Type, r.Content, r.Ttl, r.Priority, r.Id)

	return err
}

func (r *Record) Delete() error {
	sql := "DELETE FROM records WHERE id = $1"

	if config.DbType == "mysql" {
		sql = "DELETE FROM records WHERE id = ?"
	}

	_, err := db.Exec(sql, r.Id)

	return err
}
