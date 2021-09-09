
CREATE TABLE IF NOT EXISTS entities
(
    entity_id   blob,
    replica_set bigint,
    created_on  bigint,
    edited_on   bigint,
    PRIMARY KEY ((entity_id))
);
