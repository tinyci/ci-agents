create table submissions (
  id serial primary key,
  user_id int,
  head_ref_id int,
  base_ref_id int not null,
  created_at timestamptz not null default now()
);
--
alter table tasks add column submission_id int not null default 0;
--
alter table tasks alter column submission_id drop default;
--
create index task_submission_idx on tasks (id, submission_id);
