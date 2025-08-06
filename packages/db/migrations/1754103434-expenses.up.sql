CREATE TABLE IF NOT EXISTS expenses (
    id integer primary key autoincrement,
    spender_id text,
    label text,
    amount integer,
    datetime integer default (strftime('%s', 'now')),
    author_id text,
    foreign key(spender_id) references profiles(id),
    foreign key(author_id) references profiles(id)
);
