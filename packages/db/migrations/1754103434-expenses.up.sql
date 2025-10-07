CREATE TABLE IF NOT EXISTS categories (
    id integer primary key autoincrement,
    label text,
    color_fg text,
    color_bg text,
    icon_name text
);

INSERT INTO categories (id, label, color_fg, color_bg, icon_name)
VALUES (1, "Sans catégorie", "E0D071", "8B8874", "domino_mask");


INSERT INTO categories (label, color_fg, color_bg, icon_name) VALUES
("Bouffe", "F56027", "F5A527", "bakery_dining"),
("Divertissement", "F5BD51", "F5E451", "chess_knight"),
("Trucs pour être des humains", "F586E3", "F587A0", "shower"),
("Déplacement", "D74CF5", "F54CC8", "moved_location"),
("Paiements", "ADA0C3", "4B366E", "payment_card"),
("Maison", "2CCBE0", "43A8B6", "bungalow"),
("Minous", "2E8B35", "00E010", "pets")
;


CREATE TABLE IF NOT EXISTS expenses (
    id integer primary key autoincrement,
    spender_id text,
    label text,
    amount integer,
    datetime integer default (strftime('%s', 'now')),
    author_id text,
    category_id integer not null on conflict replace default(1),
    foreign key(spender_id) references profiles(id),
    foreign key(author_id) references profiles(id),
    foreign key(category_id) references categories(id)
);
