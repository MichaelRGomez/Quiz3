--File: todoApi/backend/migrations/000001_add_tasks_indexes.up.sql
create index if not exists tasks_title_idx on task_list using gin(to_tsvector('simple', title));
create index if not exists tasks_description_idx on task_list using gin(to_tsvector('simple', description));