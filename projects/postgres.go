// Package projects will manage all projects requirements
package projects

import (
	"database/sql"

	"github.com/Lord-Y/cypress-parallel-api/commons"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/syyongx/php2go"
)

// create will insert projects in DB
func (p *projects) create() (z int64, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return z, err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO projects(project_name, team_id, repository, branch) VALUES($1, $2, $3, $4) RETURNING project_id")
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}
	err = stmt.QueryRow(
		php2go.Addslashes(p.Name),
		p.TeamID,
		php2go.Addslashes(p.Repository),
		php2go.Addslashes(p.Branch),
	).Scan(&z)
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}
	defer stmt.Close()
	return z, nil
}

// read will return all teams with range limit settings
func (p *getProjects) read() (z []map[string]interface{}, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT *, (SELECT count(project_id) FROM projects) total FROM projects ORDER BY date DESC OFFSET $1 LIMIT $2")
	if err != nil && err != sql.ErrNoRows {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(
		p.StartLimit,
		p.EndLimit,
	)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	columns, err := rows.Columns()
	if err != nil {
		return
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	m := make([]map[string]interface{}, 0)
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return
		}
		var value string
		sub := make(map[string]interface{})
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = php2go.Stripslashes(string(col))
			}
			sub[columns[i]] = value
		}
		m = append(m, sub)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return m, nil
}

// GetProjectIDForUnitTesting in only for unit testing purpose and will return project_id and team_id field
func GetProjectIDForUnitTesting() (z map[string]string, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return z, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT project_id,team_id FROM projects LIMIT 1")
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return z, err
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	m := make(map[string]string)
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return
		}
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = php2go.Stripslashes(string(col))
			}
			m[columns[i]] = value
		}
	}
	if err = rows.Err(); err != nil {
		return z, err
	}
	return m, nil
}

// update will insert environments in DB
func (p *updateProjects) update() (err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE projects SET project_name = $1, team_id = $2, repository = $3, branch = $4, specs = $5, scheduling = $6, scheduling_enabled = $7, max_pods = $8 WHERE project_id = $9")
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	err = stmt.QueryRow(
		php2go.Addslashes(p.Name),
		p.TeamID,
		php2go.Addslashes(p.Repository),
		php2go.Addslashes(p.Branch),
		php2go.Addslashes(p.Specs),
		php2go.Addslashes(p.Scheduling),
		p.SchedulingEnabled,
		p.MaxPods,
		p.ProjectID,
	).Scan()
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	defer stmt.Close()
	return nil
}

// delete will delete project in DB
func (p *deleteProject) delete() (err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM projects WHERE project_id = $1")
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	err = stmt.QueryRow(
		p.ProjectID,
	).Scan()
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	defer stmt.Close()
	return nil
}