-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
  id bigserial not null,
  username varchar(45) default null,
  password varchar(255) default null,
  email varchar(45) default null,
  user_type varchar(45) default null,
  uuid char(36),
  status char(1) default 'O',
  active boolean default null,
  last_login timestamptz null,
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz null,
  CONSTRAINT users_pkey PRIMARY KEY (id),
  CONSTRAINT username_alphanumeric_check CHECK (username ~ '^[a-zA-Z0-9]*$')
);

CREATE UNIQUE INDEX username_unique_idx ON users (username);
CREATE UNIQUE INDEX email_unique_idx ON users (email);

-- +migrate Down
DROP TABLE IF EXISTS users;