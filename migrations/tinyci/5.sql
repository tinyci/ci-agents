alter table tasks drop column ref_id;
--
alter table tasks drop column base_sha;
--
alter table tasks drop column parent_id;
--
alter table submissions alter column head_ref_id set not null;
