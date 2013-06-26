package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/bmizerany/pq"
	"strconv"
	"strings"
	"time"
)

type Parcel struct {
	ParcelId int `json:"parcelId"`
	// Pointers for sql.NullString handling in Scan, will marshall to null
	Address      *string `json:"address"`
	Owner1       *string `json:"owner1"`
	Owner2       *string `json:"owner2"`
	BuildingCode *string `json:"buildingCode"`
	BuildingDesc *string `json:"buildingDesc"`
	OpaId        *string `json:"opaId"`
	Geom         *string `json:"geom"`
	Pos          *string `json:"pos"`
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

func ScanCollectionRows(rs sql.Rows, expectList bool) ([]Collection, error) {
	var (
		cs  []Collection
		c   *Collection
		err error
	)

	for rs.Next() {
		if expectList {
			c, err = ScanCollectionListRow(&rs)
		} else {
			c, err = ScanCollectionRow(&rs)
		}
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

func ScanCollectionListRow(s Scanner) (*Collection, error) {
	var c Collection
	var pids string

	err := s.Scan(&c.Id, &c.Title, &c.Desc, &c.Owner, &c.Public, &c.Created, &c.Modified, &pids)
	sids := strings.Split(pids, ",")
	ids := make([]int, 0, len(sids))
	for _, sid := range sids {
		id, err := strconv.Atoi(sid)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	c.ParcelIds = ids
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
	sql := `SELECT id, title, description, owner, public, created, modified,
                array_to_string(array_agg(parcelid), ',')
            FROM collections c, collection_parcels cp 
            WHERE c.id = cp.collectionid AND owner = $1 
            GROUP BY id;`

	if s, err := DbConn.Prepare(sql); err != nil {
		return nil, err
	} else {
		if rs, err := s.Query(username); err != nil {
			return nil, err
		} else {
			return ScanCollectionRows(*rs, true)
		}
	}
}

func RemoveParcelFromCollection(username string, cid, pid int) (*Collection, error) {
	sql := `DELETE FROM collection_parcels WHERE collectionid = $1 and parcelid = $2;`
	c, err := CollectionById(cid)
	if err != nil {
		return nil, err
	}
	return ExecuteOnParcelCollection(sql, username, c, pid)
}

func AddParcelToCollectionById(username string, cid, pid int) (*Collection, error) {
	c, err := CollectionById(cid)
	if err != nil {
		return nil, err
	}
	return AddParcelToCollection(username, c, pid)
}

func AddParcelToCollection(username string, c *Collection, pid int) (*Collection, error) {
	sql := `INSERT INTO collection_parcels VALUES ($1, $2);`
	return ExecuteOnParcelCollection(sql, username, c, pid)
}

func ExecuteOnParcelCollection(sql, username string, c *Collection, pid int) (*Collection, error) {
	if c.Owner != username { // TODO: admin role
		return nil, errors.New("Not authorized to change collection")
	} else {
		if s, err := DbConn.Prepare(sql); err != nil {
			return nil, err
		} else {
			if _, err := s.Exec(c.Id, pid); err != nil {
				return nil, err
			} else {
				return c, nil
			}
		}
	}
}

func AddCollection(c *Collection) error {
	// The pq driver doesn't appear to support the LastInsertId command
	// so fake it with the native postgres `returning` statement
	sql := `INSERT INTO collections (title, description, owner, public)
                VALUES ($1, $2, $3, $4) returning id, created, modified;`
	if s, err := DbConn.Prepare(sql); err != nil {
		return err
	} else {
		err := s.QueryRow(c.Title, c.Desc, c.Owner, c.Public).Scan(&c.Id, &c.Created, &c.Modified)
		if err != nil {
			return err
		} else {
			for _, parcelId := range c.ParcelIds {
				AddParcelToCollection(c.Owner, c, parcelId)
			}
			return nil
		}
	}
}

type Collection struct {
	Id        int        `json:"id" schema:"-"`
	Title     string     `json:"title" schema:"title"`
	Desc      string     `json:"desc" schema:"desc"`
	Parcels   []Parcel   `json:"parcels,omitempty" schema:"-"`
	ParcelIds []int      `json:"parcelIds,omitempty" schema:"parcelIds"`
	Owner     string     `json:"owner" schema:"-"`
	Public    bool       `json:"public" schema:"public"`
	Created   *time.Time `json:"created"`
	Modified  *time.Time `json:"modified"`
}
