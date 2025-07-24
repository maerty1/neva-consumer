SET timezone = 'UTC';

CREATE TABLE if not exists scada_consumer
(
    data_type VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'), 
    other_data JSONB
);

CREATE TABLE if not exists accounts
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    token TEXT NOT NULL,
    server_host VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS measure_point_types (
    id INTEGER PRIMARY KEY,
    title VARCHAR(255) NOT NULL
);


CREATE TABLE IF NOT EXISTS measure_point_type_parameters (
    title VARCHAR(255) NOT NULL,
    type_id INTEGER NOT NULL,
    type_key VARCHAR(255) DEFAULT 'default',
    min_zoom SMALLINT DEFAULT 14,
    max_zoom SMALLINT DEFAULT NULL,
    FOREIGN KEY (type_id) REFERENCES measure_point_types(id)
);

CREATE TABLE if not exists measure_points 
(
    id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    address VARCHAR(255) NULL,
    full_title VARCHAR(255) NULL,
    system_type VARCHAR(255) NULL,
    device_id INTEGER NULL,
    description TEXT,
    type_id INTEGER NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),
    elem_id INTEGER NULL,

    FOREIGN KEY (account_id) REFERENCES accounts(id),
    FOREIGN KEY (type_id) REFERENCES measure_point_types(id)
);

CREATE TABLE if not exists measure_points_metadata
(   
    lat DOUBLE PRECISION,
    lon DOUBLE PRECISION,
    measure_point_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),

    FOREIGN KEY (measure_point_id) REFERENCES measure_points(id)
);

CREATE TABLE if not exists measure_points_data
(   
    measure_point_id INTEGER NOT NULL,
    datetime TIMESTAMP NOT NULL,
    values JSONB,

    PRIMARY KEY (measure_point_id, datetime)
);

CREATE TABLE if not exists measure_points_data_day
(   
    measure_point_id INTEGER NOT NULL,
    datetime TIMESTAMP NOT NULL,
    values JSONB,

    PRIMARY KEY (measure_point_id, datetime)
);

CREATE TABLE if not exists accounts_sync_log
(
    account_id INTEGER NOT NULL,
    measure_point_id INTEGER,
    message TEXT,
    level VARCHAR(255) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC')
);

CREATE TABLE if not exists measure_points_poll_log 
(
    poll_id INTEGER NOT NULL,
    -- Ответ из Лэрс во время попытки опроса --
    message TEXT,
    -- Результат опроса --
    status TEXT,
    measure_point_id INTEGER,
    account_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC')
);

CREATE TABLE if not exists measure_points_poll_retry 
(
    original_poll_id INTEGER NOT NULL,
    retrying_poll_id INTEGER NOT NULL,
    status TEXT DEFAULT 'PENDING', 
    created_at TIMESTAMP NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC')
);

CREATE INDEX ix_measure_points_data ON public.measure_points_data  USING btree ("datetime", measure_point_id);


CREATE TABLE if not exists static_measure_points 
(
    id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL
);

CREATE TABLE if not exists static_measure_points_data
(
    static_measure_point_id INTEGER NOT NULL,
    datetime TIMESTAMP NOT NULL,
    value TEXT,

    PRIMARY KEY (static_measure_point_id, datetime)
);

CREATE TABLE measure_types (
    id SERIAL PRIMARY KEY,
    measure_type VARCHAR(50) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    unit VARCHAR(50)
);

CREATE TABLE if not exists scada_measure_points 
(
    id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),
    elem_id INTEGER NULL,

    FOREIGN KEY (account_id) REFERENCES accounts(id),
    UNIQUE (account_id, title)
);

CREATE TABLE if not exists scada_consumer
(
    data_type VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'), 
    other_data JSONB
);

CREATE TABLE if not exists scada_rawdata 
(
    scada_measure_point_id INTEGER NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    varname TEXT NULL,
    measurement_type_id INTEGER NULL,
    value TEXT,
    raw_packet JSONB
);

CREATE TABLE IF NOT EXISTS public.weather_data (
	"hour" integer,
	date text,
	"year" integer,
	outdoor_temperature real,
	dew_point_temperature real,
	humidity integer,
	effective_outdoor_temperature integer,
	wind_direction integer,
	wind_speed character varying(50),
	max_wind_speed integer,
	min_wind_speed integer,
	wind_gust_speed integer,
	"day" integer,
	"month" integer
);


INSERT INTO measure_point_types (id, title)
VALUES
    (1, 'Труба'),
    (2, 'Компенсатор'),
    (3, 'Тепловая камера'),
    (4, 'Источник'),
    (5, 'ЦТП'),
    (6, 'Задвижка'),
    (7, 'Потребитель'),
    (8, 'Насосная станция');


INSERT INTO measure_types (measure_type, title, unit)
VALUES 
('T_in', 'Температура прямой', '°C'),
('T_out', 'Температура обратной', '°C'),
('T_delta', 'Разность температур', '°C'),
('T_cw', 'Температура холодной', '°C'),
('P_in', 'Давление прямой', 'МПа'),
('P_out', 'Давление обратной', 'МПа'),
('P_delta', 'Перепад давления', 'МПа'),
('P_cw', 'Давление холодной', 'МПа'),
('M_in', 'Расход прямой', 'т/ч'),
('M_out', 'Расход обратной', 'т/ч'),
('M_delta', 'Массовый расход', ''),
('V_in', 'Объем по подающей магистрали', ''),
('V_out', 'Объем по обратной магистрали', ''),
('V_delta', 'Объемный расход', ''),
('Q_in', 'Тепло прямой', 'Гк/мес'),
('Q_out', 'Тепло обратной', 'Гк/мес'),
('Q_delta', 'Тепло', 'Гк/мес');

CREATE TABLE weather_data_raw (
    dateutc TIMESTAMP NOT NULL, 
    tempinf DOUBLE PRECISION,
    humidityin DOUBLE PRECISION,
    baromrelin DOUBLE PRECISION,
    baromabsin DOUBLE PRECISION,
    tempf DOUBLE PRECISION, 
    humidity DOUBLE PRECISION,
    winddir DOUBLE PRECISION, 
    windspeedmph DOUBLE PRECISION,
    windgustmph DOUBLE PRECISION,
    maxdailygust DOUBLE PRECISION,
    solarradiation DOUBLE PRECISION, 
    uv DOUBLE PRECISION,                  
    rainratein DOUBLE PRECISION,
    eventrainin DOUBLE PRECISION, 
    hourlyrainin DOUBLE PRECISION, 
    dailyrainin DOUBLE PRECISION,
    weeklyrainin DOUBLE PRECISION,
    monthlyrainin DOUBLE PRECISION,
    yearlyrainin DOUBLE PRECISION,
    totalrainin DOUBLE PRECISION,
    wh65batt DOUBLE PRECISION,    
    freq VARCHAR(20)
);