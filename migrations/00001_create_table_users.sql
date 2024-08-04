create table users (
    usr_id uuid primary key default gen_random_uuid(),
    usr_name text not null
);
