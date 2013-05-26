create table if not exists collections (
	id 		serial PRIMARY KEY,
	title 		varchar(255) NOT NULL,
	description	text,
	owner 		varchar(55),
	public 		boolean,
	created 	timestamp DEFAULT current_timestamp,
	modified 	timestamp DEFAULT current_timestamp
);

create index collectionid_idx on collections (id);

create table if not exists collection_parcels (
	collectionid	integer NOT NULL,
	parcelid	integer NOT NULL
);
create index collection_parcels_idx on collection_parcels (collectionid);

create table if not exists users (
	username	varchar(255),
	email		varchar(255),
	password	varchar(255),
	joined		timestamp DEFAULT current_timestamp
);
create index on users (username);

create table if not exists opa (
	acct_num	varchar(255),
	address		varchar(255),
	unit		varchar(255),
	homestd_ex	varchar(255),
	prop_cat	varchar(255),
	prop_type	varchar(255),
	num_stor	real,
	mktval_13	integer,
	landval_13	integer,
	impval_13	integer,
	abat_ex_13	integer,
	mktval_14	integer,
	landval_14	integer,
	impval_14	integer,
	abat_ex_14	integer,
	plotarea	real,
	zone		varchar(255),
	extcond		varchar(255),
	totdwellarea	integer,
	ownr_nam	varchar(255),
	secd_nam	varchar(255),
	dist_from	varchar(255),
	ploc_zip5	varchar(255),
	titl_date	date,
	sale_price	integer,
	geom		geometry(Point, 4326)
);

create index on opa (acct_num);
