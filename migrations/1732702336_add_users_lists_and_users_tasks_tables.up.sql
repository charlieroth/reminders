-- Write your up sql migration here

-- Create users_lists junction table
CREATE TABLE users_lists (
    user_id uuid NOT NULL,
    list_id uuid NOT NULL,
    CONSTRAINT users_lists_pkey PRIMARY KEY (user_id, list_id),
    CONSTRAINT users_lists_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT users_lists_list_id_fkey FOREIGN KEY (list_id) REFERENCES lists(id)
);

-- Create users_tasks junction table
CREATE TABLE users_tasks (
    user_id uuid NOT NULL,
    task_id uuid NOT NULL,
    CONSTRAINT users_tasks_pkey PRIMARY KEY (user_id, task_id),
    CONSTRAINT users_tasks_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT users_tasks_task_id_fkey FOREIGN KEY (task_id) REFERENCES tasks(id)
);

-- Create indexes for faster lookups
CREATE INDEX idx_users_lists_user_id ON users_lists(user_id);
CREATE INDEX idx_users_lists_list_id ON users_lists(list_id);
CREATE INDEX idx_users_tasks_user_id ON users_tasks(user_id);
CREATE INDEX idx_users_tasks_task_id ON users_tasks(task_id);