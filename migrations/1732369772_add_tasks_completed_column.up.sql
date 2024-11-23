-- Write your up sql migration here
ALTER TABLE tasks ADD COLUMN completed BOOLEAN NOT NULL DEFAULT FALSE;
