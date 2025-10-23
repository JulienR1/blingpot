alter table categories
add column order_index integer;

update categories set order_index = 8 where id = 1;
update categories set order_index = 2 where id = 2;
update categories set order_index = 5 where id = 3;
update categories set order_index = 6 where id = 4;
update categories set order_index = 3 where id = 5;
update categories set order_index = 1 where id = 6;
update categories set order_index = 4 where id = 7;
update categories set order_index = 7 where id = 8;
