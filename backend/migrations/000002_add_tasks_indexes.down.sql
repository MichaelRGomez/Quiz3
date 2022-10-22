--File: todoApi/backend/migrations/000001_add_tasks_indexes.down.sql
drop index if exists tasks_title_idx;
drop index if exists tasks_description_idx;