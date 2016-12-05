CREATE TABLE logs
(
    build_id uuid,
    created_at timestamp,
    entry_id uuid,
    container_id text,
    container_name text,
    payload text,
    PRIMARY KEY (build_id, created_at, entry_id, container_id)
) WITH CLUSTERING ORDER BY (created_at ASC);
