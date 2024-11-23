-- Write your up sql migration here
CREATE TABLE lists_tasks (
    list_id UUID NOT NULL REFERENCES lists(id),
    task_id UUID NOT NULL REFERENCES tasks(id),
    PRIMARY KEY (list_id, task_id)
);
