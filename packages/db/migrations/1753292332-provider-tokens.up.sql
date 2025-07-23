CREATE TABLE IF NOT EXISTS provider_tokens (
    sub text primary key,
    access_token text,
    refresh_token text,
    foreign key(sub) references profiles(sub)
);
