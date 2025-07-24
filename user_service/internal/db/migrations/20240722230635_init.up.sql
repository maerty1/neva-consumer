CREATE TABLE users (
    id SERIAL PRIMARY KEY, 
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL, 
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS user_settings (
    bounding_box_top_left_lat DECIMAL(9, 6) NOT NULL,    -- широта верхнего левого угла
    bounding_box_top_left_lon DECIMAL(9, 6) NOT NULL,    -- долгота верхнего левого угла
    bounding_box_bottom_right_lat DECIMAL(9, 6) NOT NULL, -- широта правого нижнего угла
    bounding_box_bottom_right_lon DECIMAL(9, 6) NOT NULL, -- долгота правого нижнего угла
    default_zoom SMALLINT CHECK (default_zoom >= 0 AND default_zoom <= 16), -- зум от 0 до 16
    default_center_lat DECIMAL(9, 6) NOT NULL, -- широта центра
    default_center_lon DECIMAL(9, 6) NOT NULL, -- долгота центра

    user_id INTEGER NOT NULL UNIQUE,
    FOREIGN KEY (user_id) REFERENCES users(id)
)

