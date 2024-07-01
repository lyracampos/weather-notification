CREATE TABLE users (
      id bigserial PRIMARY KEY,
      first_name text NOT NULL,
      last_name text NOT NULL,
      email text NOT NULL,
      phone text NOT NULL,
      city_id int NOT NULL,
      opt_in boolean NOT NULL default true,
      created_at timestamp NOT NULL DEFAULT now(),
      updated_at timestamp NOT NULL DEFAULT now(),
      UNIQUE(email)
);