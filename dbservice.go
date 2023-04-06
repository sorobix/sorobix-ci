package main

import (
	"github.com/genjidb/genji"
	"github.com/genjidb/genji/document"
	"github.com/genjidb/genji/types"
)

type Deployment struct {
	ID     string `genji:"id" json:"ID,omitempty"`
	StdOut []byte `genji:"stdout" json:"stdOut,omitempty"`
	StdErr []byte `genji:"stderr" json:"stdErr,omitempty"`
	Time   string `genji:"creation" json:"time,omitempty"`
}

type Repository interface {
	InsertDeployment(deployment Deployment) error
	FetchDeployments() ([]Deployment, error)
	CreateTable() error
}

type repository struct {
	db *genji.DB
}

func (r repository) InsertDeployment(deployment Deployment) error {
	err := r.db.Exec(`INSERT INTO deployments VALUES ?`, &deployment)
	if err != nil {
		return err
	}
	return nil
}

func (r repository) FetchDeployments() ([]Deployment, error) {
	var deps []Deployment
	res, err := r.db.Query(`SELECT * FROM deployments`)
	if err != nil {
		return nil, err
	}
	err = res.Iterate(func(d types.Document) error {
		var dep Deployment
		err = document.StructScan(d, &dep)
		if err != nil {
			return err
		}
		deps = append(deps, dep)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return deps, nil
}

func (r repository) CreateTable() error {
	err := r.db.Exec(`
CREATE TABLE IF NOT EXISTS deployments (
    	id TEXT NOT NULL,
        stdout blob,
    	stderr blob,
    	creation TEXT 
    );
`)
	if err != nil {
		return err
	}
	return nil
}

// NewRepo creates a new instance of this repository
func NewRepo(db *genji.DB) Repository {
	return &repository{
		db: db,
	}
}
