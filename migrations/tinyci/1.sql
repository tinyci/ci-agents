create index refs_repository_idx on refs (id, repository_id);
--
create index task_parent_idx on tasks (id, parent_id);
--
create index run_name_id on runs (id, name);
