create index users_username_idx on users (id, username);
--
create index user_errors_user_id_idx on user_errors (id, user_id);
--
create index session_key_expires_idx on sessions (key, expires_on);
--
create index runs_task_idx on runs (id, task_id);
