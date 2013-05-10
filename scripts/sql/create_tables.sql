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
