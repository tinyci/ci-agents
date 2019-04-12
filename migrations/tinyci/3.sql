alter table users drop column login_token;
--
alter table users add column login_token bytea;
