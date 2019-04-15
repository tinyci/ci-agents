drop table owners; -- squash this migration before release
--
alter table repositories add column owner_id integer;
--
create table user_capabilities (
  user_id integer not null,
  name varchar not null,

  PRIMARY KEY (user_id, name)
);
