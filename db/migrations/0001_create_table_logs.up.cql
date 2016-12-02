CREATE TABLE logs
(
    build_id uuid,
    created_at timestamp,
    container_id text,
    container_name text,
    payload text,
    PRIMARY KEY (build_id, created_at, container_id)
) WITH CLUSTERING ORDER BY (created_at ASC);
