-- Write your down sql migration here
DROP TABLE users_lists;
DROP TABLE users_tasks;
DROP INDEX idx_users_lists_user_id;
DROP INDEX idx_users_lists_list_id;
DROP INDEX idx_users_tasks_user_id;
DROP INDEX idx_users_tasks_task_id;