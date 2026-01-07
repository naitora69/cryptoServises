CREATE TABLE proposal (
    id serial primary key,
    hex_id varchar(256) unique not null ,
    title text,
    author varchar(256),
    created_at timestamp,
    start_at timestamp,
    end_at timestamp,
    snapshot bigint,
    state text,
    choices json,
    space_id varchar(256),
    space_name varchar(256)
);

CREATE TABLE event_outbox (
    id serial primary key,
    hex_id varchar(256) references proposal(hex_id) on delete cascade,
    event_type text,
    created_at timestamp,
    processed_at timestamp NULL
);