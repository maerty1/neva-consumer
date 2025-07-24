set role zulu_updater;

create table zulu.object_records_new(
    rec_id 		serial4,
    inserted_ts TIMESTAMPTZ DEFAULT NOW(),
    elem_id		INT4,
    val_name	varchar(64),
    val			varchar(64)
);