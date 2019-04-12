alter table users add column login_token varchar;
--
create index users_id_login_token_idx on users (id, login_token);
--
create index users_name_login_token_idx on users (username, login_token);
--
