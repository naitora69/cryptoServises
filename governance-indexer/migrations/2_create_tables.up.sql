CREATE TABLE IF NOT EXISTS spaces
(
    id serial primary key,
    space_id varchar(256) unique not null ,
    name text,
    about text,
    network text,
    symbol text,
    created integer,
    strategies_name text,
    admins json,
    members json,
    filters_min_score integer,
    filters_only_members bool
);

CREATE TABLE IF NOT EXISTS spaces_outbox
(
    id serial primary key,
    space_id varchar(256) unique not null references spaces(space_id) on delete cascade,
    event_type text,
    created_at timestamp,
    processed_at timestamp NULL
);