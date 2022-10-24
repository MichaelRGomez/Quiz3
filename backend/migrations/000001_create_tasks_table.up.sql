--File: todoApi/backend/migrations/000001_create_tasks_table.up.sql
create table if not exists task_list(
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone not null default now(),
    title text not null,
    description text not null,
    completed boolean not null,
    version int not null default 1
);