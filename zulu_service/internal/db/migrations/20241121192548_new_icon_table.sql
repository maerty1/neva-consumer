-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS icons (
 icon_id SERIAL PRIMARY KEY,
 icon TEXT NOT NULL,
 path TEXT NOT NULL
);

CREATE TABLE if not exists point_to_icon (
 icon_id INTEGER NOT NULL,
 FOREIGN KEY(icon_id) REFERENCES public.icons(icon_id),
 elem_id INTEGER NOT NULL, 
 FOREIGN KEY (elem_id) REFERENCES zulu.elems_metadata(elem_id),
 CONSTRAINT metadata_icon_pkey PRIMARY KEY (elem_id, icon_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS icons;
DROP TABLE IF EXISTS point_to_icon;
-- +goose StatementEnd
