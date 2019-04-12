create table sessions (
  key varchar(4096) primary key,
  values varchar(4096) not null,
  expires_on timestamptz not null
);
--
create table o_auths (
  state varchar(4096) primary key,
  expires_on timestamptz not null
);
--
create table users (
  id serial primary key,
  username varchar not null unique,
  token varchar not null,
  last_scanned_repos timestamptz
);
--
create table user_errors (
  id serial primary key,
  user_id int not null,
  error varchar not null
);
--
create table subscriptions (
  user_id int not null,
  repository_id int not null,
  PRIMARY KEY(user_id, repository_id)
);
--
create table repositories (
  id serial primary key,
  name varchar not null unique,
  private boolean not null,
  github text not null,
  disabled bool default false,
  auto_created bool not null,
  hook_secret varchar not null
);
--
create index repo_name_idx on repositories (id, name);
--
create table owners (
  user_id int not null,
  repository_id int not null,

  PRIMARY KEY(user_id, repository_id)
);
--
create table refs (
  id serial primary key,
  repository_id int not null,
  ref varchar not null,
  sha varchar not null,
  UNIQUE(repository_id, ref, sha)
);
--
create table tasks (
  id serial primary key,
  ref_id int not null,
  parent_id int not null,
  base_sha varchar not null,
  pull_request_id int not null,
  status boolean,
  task_settings text not null, 
  created_at timestamptz not null default now(),
  started_at timestamptz,
  finished_at timestamptz,
  canceled boolean not null default false
);
--
create index task_ref_idx on tasks (id, ref_id);
--
create table runs (
  id serial primary key,
  task_id int not null,
  name varchar not null,
  run_settings text not null, 
  status boolean,
  created_at timestamptz not null default now(),
  started_at timestamptz,
  finished_at timestamptz
);
--
create index run_task_idx on runs (id, task_id);
--
create table queue_items (
  id serial primary key,
  run_id int unique,
  running bool not null default false,
  running_on varchar,
  queue_name varchar not null,
  started_at timestamptz
);
--
create index queue_running_idx on queue_items (run_id, running);
--
create index queue_name_idx on queue_items (run_id, running, queue_name);
