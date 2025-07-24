set role zulu_updater;

create table zulu.objects_geometry_log_new(
    rec_id 		    serial4,
    inserted_ts     TIMESTAMPTZ DEFAULT NOW(),
    elem_id		    INT4,
    zws_type        INT4,
    zws_mode        INT4,
    zws_geometry    GEOMETRY,
    zws_linecolor   INT8,
    parent_id       INT4,
);