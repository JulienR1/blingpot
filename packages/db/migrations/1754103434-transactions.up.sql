CREATE TABLE IF NOT EXISTS transactions (
    id integer primary key autoincrement,
    profile_id text,
    label text,
    amount integer,
    is_expense integer default true,
    datetime integer default (strftime('%s', 'now')),
    author_id text,
    foreign key(profile_id) references profiles(id),
    foreign key(author_id) references profiles(id)
);
