package main

import (
	"log"
	"time"
)

type Domain struct {
	Id        int        `json:"id"`
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func NewDomain(name string) *Domain {
	return &Domain{Name: name, Type: "NATIVE"}
}

func FindDomain(id int) (*Domain, error) {
	sql := "SELECT id, name, type FROM domains WHERE id = $1"

	if conf.DbType == "mysql" {
		sql = "SELECT id, name, type FROM domains WHERE id = ?"
	}

	row := db.QueryRow(sql, id)

	d := &Domain{}
	err := row.Scan(
		&d.Id,
		&d.Name,
		&d.Type,
	)

	return d, err
}

func AllDomains() []*Domain {
	rows, err := db.Query("select id, name, type from domains")
	if err != nil {
		log.Fatal(err)
	}

	domains := []*Domain{}

	for rows.Next() {
		d := &Domain{}
		err := rows.Scan(
			&d.Id,
			&d.Name,
			&d.Type,
		)
		if err != nil {
			log.Fatal(err)
		}

		domains = append(domains, d)
	}

	return domains
}

// Create new domain, returning the id
func (d *Domain) Create() (int64, error) {

	if conf.DbType == "mysql" {
		// MySQL supports LastInsertId()
		sql := "INSERT INTO domains (name, type) VALUES (?, ?)"

		res, err := db.Exec(
			sql,
			d.Name,
			d.Type,
		)

		if err != nil {
			return -1, err
		}

		domain_id, id_err := res.LastInsertId()
		if id_err != nil {
			return -1, id_err
		}

		return domain_id, err
	} else {
		// PostgreSQL driver does not support it, use RETURNING instead
		sql := "INSERT INTO domains (name, type) VALUES ($1, $2) RETURNING id"
		var domain_id int64
		err := db.QueryRow(
			sql,
			d.Name,
			d.Type,
		).Scan(&domain_id)

		if err != nil {
			return -1, err
		}

		return domain_id, err
	}
}

func (d *Domain) Update() error {
	sql := "UPDATE domains SET name=$1 WHERE id=$2"

	if conf.DbType == "mysql" {
		sql = "UPDATE domains SET name=? WHERE id=?"
	}

	_, err := db.Exec(sql, d.Name, d.Id)

	return err
}

func (d *Domain) Delete() error {
	sql := "DELETE FROM domains WHERE id = $1"
	recSql := "DELETE FROM records WHERE domain_id = $1"

	if conf.DbType == "mysql" {
		sql = "DELETE FROM domains WHERE id = ?"
		recSql = "DELETE FROM records WHERE domain_id = ?"
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Delete records that belong to this domain and rollback on err
	_, err = tx.Exec(recSql, d.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete domain and rollback on err
	_, err = tx.Exec(sql, d.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
