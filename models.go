package main

import (
	"database/sql"
	"fmt"
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
	Geom         *string
	Pos          *string
}

type Scanner interface {
	Scan(dest ...interface{}) error
}

func ScanParcelRow(s Scanner) (*Parcel, error) {
	var p Parcel
	err := s.Scan(&p.ParcelId, &p.Address, &p.Owner1, &p.Owner2, &p.BuildingCode, &p.BuildingDesc, &p.OpaId, &p.Geom, &p.Pos)
	return &p, err
}

func ScanParcelRows(rs sql.Rows) ([]Parcel, error) {
	var parcels []Parcel
	for rs.Next() {
		p, err := ScanParcelRow(&rs)
		if err != nil {
			return nil, err
		} else {
			parcels = append(parcels, *p)
		}
	}
	return parcels, nil
}

func ParcelById(id int) (*Parcel, error) {
	sql := `SELECT parcelid, address, owner1, owner2, bldg_code, bldg_desc, brt_id,
                ST_AsGeoJSON(geom), ST_AsGeoJSON(pos)
            FROM pwd_parcels
            WHERE parcelid = $1;`
	if s, err := DbConn.Prepare(sql); err != nil {
		return nil, err
	} else {
		return ScanParcelRow(s.QueryRow(id))
	}
}

func ParcelsByCid(cid int) ([]Parcel, error) {
	sql := `SELECT p.parcelid, p.address, p.owner1, p.owner2, p.bldg_code, 
                p.bldg_desc, p.brt_id, ST_AsGeoJSON(p.geom), ST_AsGeoJSON(p.pos)
            FROM pwd_parcels p, collection_parcels c
            WHERE p.parcelid = c.parcelid and c.collectionid = $1;`
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

func ParcelByLocation(lat, lon float64) (*Parcel, error) {
	pointWkt := fmt.Sprintf("POINT (%f %f)", lon, lat)
	sql := `SELECT parcelid, address, owner1, owner2, bldg_code, bldg_desc,
                brt_id, ST_AsGeoJSON(geom), ST_AsGeoJSON(pos)
            FROM pwd_parcels
            WHERE ST_Intersects(ST_GeomFromText($1, 4326), geom) = true;`
	if s, err := DbConn.Prepare(sql); err != nil {
		return nil, err
	} else {
		return ScanParcelRow(s.QueryRow(pointWkt))
	}
}

func ScanCollectionRows(rs sql.Rows) ([]Collection, error) {
	var cs []Collection
	for rs.Next() {
		c, err := ScanCollectionRow(&rs)
		if err != nil {
			return nil, err
		} else {
			cs = append(cs, *c)
		}
	}
	return cs, nil
}

func ScanCollectionRow(s Scanner) (*Collection, error) {
	var c Collection
	err := s.Scan(&c.Id, &c.Title, &c.Desc, &c.Owner, &c.Public, &c.Created, &c.Modified)
	return &c, err
}

func CollectionById(id int) (*Collection, error) {
	sql := `SELECT id, title, description, owner, public, created, modified 
            FROM Collections 
            WHERE id = $1;`
	if s, err := DbConn.Prepare(sql); err != nil {
		return nil, err
	} else {
		c, err := ScanCollectionRow(s.QueryRow(id))
		if err != nil {
			return nil, err
		}
		if c.Parcels, err = ParcelsByCid(id); err != nil {
			return c, err
		} else {
			return c, nil
		}
	}
}

func CollectionListByUser(username string) ([]Collection, error) {
	sql := `SELECT id, title, desc, owner
                FROM collections
                WHERE owner = $1;`
	if s, err := DbConn.Prepare(sql); err != nil {
		return nil, err
	} else {
		if rs, err := s.Query(username); err != nil {
			return nil, err
		} else {
			return ScanCollectionRows(*rs)
		}
	}
}

type Collection struct {
	Id       int        `json:"id"`
	Title    string     `json:"title"`
	Desc     string     `json:"desc"`
	Parcels  []Parcel   `json:"parcels"`
	Owner    string     `json:"owner"`
	Public   bool       `json:"public"`
	Created  *time.Time `json:"created"`
	Modified *time.Time `json:"modified"`
}
