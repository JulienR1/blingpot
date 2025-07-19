CREATE TABLE IF NOT EXISTS auth_providers (
    id integer primary key,
    label text
);

INSERT INTO auth_providers (id, label)
VALUES (1, "google");

CREATE TABLE IF NOT EXISTS profiles (
    sub text primary key,
    first_name text,
    last_name text,
    email text,
    picture text nullable,
    provider_id integer,
    foreign key(provider_id) references auth_providers(id)
);
