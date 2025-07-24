-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists icons_on_zws_type_id (
icon_id SERIAL PRIMARY KEY,
elem_id INTEGER UNIQUE,
FOREIGN KEY (elem_id) REFERENCES zulu.elems_metadata(elem_id),
icon TEXT NOT NULL,
path TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS icons_on_zws_type_id;
-- +goose StatementEnd
