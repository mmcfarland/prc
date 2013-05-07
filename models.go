package main

import (
	"database/sql"
	_ "github.com/bmizerany/pq"
	"time"
)

type Parcel struct {
	ParcelId int
	// Pointers for sql.NullString handling in Scan, will marshall to null
	Address      *string
	Owner1       *string
	Owner2       *string
	BuildingCode *string
	BuildingDesc *string
	OpaId        *string
	GeomWkt      *string
	Pop          *string
}

type Scanner interface {
	Scan(dest ...interface{}) error
}

func ScanParcelRow(s Scanner) (*Parcel, error) {
	var p Parcel
	err := s.Scan(&p.ParcelId, &p.Address, &p.Owner1, &p.Owner2, &p.BuildingCode, &p.BuildingDesc, &p.OpaId, &p.GeomWkt)
	return &p, err
}

func ScanParcelRows(rs sql.Rows) (*[]Parcel, error) {
	var parcels []Parcel
	for rs.Next() {
		p, err := ScanParcelRow(&rs)
		if err != nil {
			return nil, err
		} else {
			parcels = append(parcels, *p)
		}
	}
	return &parcels, nil
}

func ParcelById(id int) (*Parcel, error) {
	sql := "SELECT parcelid, address, owner1, owner2, bldg_code, bldg_desc, brt_id, ST_AsGeoJSON(geom) FROM pwd_parcels where parcelid = $1;"
	if s, err := DbConn.Prepare(sql); err != nil {
		return nil, err
	} else {
		return ScanParcelRow(s.QueryRow(id))
	}
}

func ParcelsByCid(cid int) (*[]Parcel, error) {
	sql := "SELECT p.parcelid, p.address, p.owner1, p.owner2, p.bldg_code, p.bldg_desc, p.brt_id, ST_AsGeoJSON(p.geom) FROM pwd_parcels p, collection_parcels c where p.parcelid = c.parcelid and c.collectionid = $1;"
	if s, err := DbConn.Prepare(sql); err != nil {
		return nil, err
	} else {
		if rs, err := s.Query(cid); err != nil {
			return nil, err
		} else {
			return ScanParcelRows(*rs)
		}
	}
}

type Collection struct {
	Id       int
	Title    string
	Desc     string
	Parcels  []Parcel
	Owner    int
	IsPublic bool
	Created  time.Time
}

type User struct {
	Id    int
	Name  string
	email string
}
