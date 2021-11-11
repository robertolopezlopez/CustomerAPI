create database customer;

create schema roberto;

create table roberto.customer
(
    mailing_id  integer primary key not null,
    email       varchar(50),
    insert_time timestamptz default current_timestamp,
    title       varchar(50),
    content     varchar(150)
);

select * from roberto.customer;

create function delete_old() returns trigger
language plpgsql
as $$
begin
    delete from roberto.customer where insert_time < current_timestamp - interval '5 minutes';
end;
$$;

CREATE TRIGGER trigger_select AFTER INSERT ON roberto.customer
    EXECUTE PROCEDURE delete_old();
