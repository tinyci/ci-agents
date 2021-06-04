-- +migrate Up

-- +migrate StatementBegin
CREATE TABLE o_auths (
    state character varying(4096) NOT NULL primary key,
    expires_on timestamp with time zone NOT NULL,
    scopes character varying DEFAULT ''::character varying NOT NULL
);
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TABLE sessions (
    key character varying(4096) NOT NULL primary key,
    "values" character varying(4096) NOT NULL,
    expires_on timestamp with time zone NOT NULL
);
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TABLE users (
    id bigserial NOT NULL primary key,
    username character varying NOT NULL,
    token bytea NOT NULL,
    last_scanned_repos timestamp with time zone,
    login_token bytea,

    UNIQUE(username)
);
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TABLE repositories (
    id bigserial NOT NULL primary key,
    name character varying NOT NULL,
    private boolean NOT NULL,
    github jsonb NOT NULL,
    disabled boolean DEFAULT false,
    auto_created boolean NOT NULL,
    hook_secret character varying NOT NULL,
    owner_id bigint NOT NULL,

    FOREIGN KEY (owner_id) REFERENCES users(id),

    UNIQUE(name)
);
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TABLE refs (
    id bigserial NOT NULL primary key,
    repository_id bigint NOT NULL,
    ref character varying NOT NULL,
    sha character varying NOT NULL,

    UNIQUE (repository_id, ref, sha),

    FOREIGN KEY (repository_id) REFERENCES repositories(id)
);
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TABLE submissions (
    id bigserial NOT NULL primary key,
    user_id bigint,
    head_ref_id bigint,
    base_ref_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    ticket_id bigint,

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (head_ref_id) REFERENCES refs(id),
    FOREIGN KEY (base_ref_id) REFERENCES refs(id)
);
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TABLE tasks (
    id bigserial NOT NULL primary key,
    status boolean,
    task_settings jsonb NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    started_at timestamp with time zone,
    finished_at timestamp with time zone,
    canceled boolean DEFAULT false NOT NULL,
    path character varying DEFAULT ''::character varying NOT NULL,
    submission_id bigint NOT NULL,

    FOREIGN KEY (submission_id) REFERENCES submissions(id)
);
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TABLE runs (
    id bigserial NOT NULL primary key,
    task_id bigint NOT NULL,
    name character varying NOT NULL,
    run_settings jsonb NOT NULL,
    status boolean,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    started_at timestamp with time zone,
    finished_at timestamp with time zone,
    ran_on character varying,

    FOREIGN KEY (task_id) REFERENCES tasks(id)
);
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TABLE queue_items (
    id bigserial NOT NULL primary key,
    run_id bigint NOT NULL,
    running boolean DEFAULT false NOT NULL,
    running_on character varying,
    queue_name character varying NOT NULL,
    started_at timestamp with time zone,

    FOREIGN KEY (run_id) REFERENCES runs(id),
    UNIQUE(run_id)
);
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TABLE subscriptions (
    user_id bigint NOT NULL,
    repository_id bigint NOT NULL,

    PRIMARY KEY (user_id, repository_id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TABLE user_capabilities (
    user_id bigint NOT NULL,
    name character varying NOT NULL,

    PRIMARY KEY (user_id, name),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TABLE user_errors (
    id bigserial NOT NULL primary key,
    user_id bigint NOT NULL,
    error character varying NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
-- +migrate StatementEnd

CREATE INDEX queue_name_idx ON queue_items USING btree (run_id, running, queue_name);
CREATE INDEX queue_running_idx ON queue_items USING btree (run_id, running);
CREATE INDEX refs_repository_idx ON refs USING btree (id, repository_id);
CREATE INDEX repo_name_idx ON repositories USING btree (id, name);
CREATE INDEX run_name_id ON runs USING btree (id, name);
CREATE INDEX run_task_idx ON runs USING btree (id, task_id);
CREATE INDEX session_key_expires_idx ON sessions USING btree (key, expires_on);
CREATE INDEX submission_base_ref_id ON submissions USING btree (id, base_ref_id);
CREATE INDEX submission_head_ref_id ON submissions USING btree (id, head_ref_id);
CREATE INDEX task_submission_idx ON tasks USING btree (id, submission_id);
CREATE INDEX user_errors_user_id_idx ON user_errors USING btree (id, user_id);
CREATE INDEX users_id_login_token_idx ON users USING btree (id, login_token);
CREATE INDEX users_name_login_token_idx ON users USING btree (username, login_token);
CREATE INDEX users_username_idx ON users USING btree (id, username);
