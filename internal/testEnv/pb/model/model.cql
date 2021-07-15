

CREATE TABLE IF NOT EXISTS tab_model_1 
(
	id  int,
	shard_key  int,
	enum  int,
	sdata  blob,
	PRIMARY KEY ((id), shard_key, enum)
) WITH CLUSTERING ORDER BY (shard_key ASC, enum ASC);

CREATE MATERIALIZED VIEW IF NOT EXISTS view_model_1_customer_sort AS
SELECT *
FROM tab_model_1
WHERE id IS NOT null
AND shard_key IS NOT null
AND enum IS NOT null
PRIMARY KEY ((enum), shard_key, id) WITH CLUSTERING ORDER BY (shard_key ASC, id ASC);




CREATE TABLE IF NOT EXISTS tab_model_2 
(
	id  bigint,
	shard_key  int,
	p_1  blob,
	sdata  blob,
	PRIMARY KEY ((id, shard_key), p_1)
) WITH CLUSTERING ORDER BY (p_1 DESC);

CREATE MATERIALIZED VIEW IF NOT EXISTS view_model_2_p_1_shard_key_id AS
SELECT *
FROM tab_model_2
WHERE id IS NOT null
AND shard_key IS NOT null
AND p_1 IS NOT null
PRIMARY KEY ((p_1), shard_key, id) WITH CLUSTERING ORDER BY (shard_key ASC, id ASC);




CREATE TABLE IF NOT EXISTS tab_model_3 
(
	id  bigint,
	shard_key  int,
	p_1  blob,
	sdata  blob,
	PRIMARY KEY ((id, shard_key), p_1)
) WITH CLUSTERING ORDER BY (p_1 DESC);

CREATE MATERIALIZED VIEW IF NOT EXISTS view_model_3_p_1_shard_key_id AS
SELECT *
FROM tab_model_3
WHERE id IS NOT null
AND shard_key IS NOT null
AND p_1 IS NOT null
PRIMARY KEY ((p_1), shard_key, id) WITH CLUSTERING ORDER BY (shard_key ASC, id ASC);

CREATE MATERIALIZED VIEW IF NOT EXISTS view_model_3_p_1_id_shard_key AS
SELECT *
FROM tab_model_3
WHERE id IS NOT null
AND shard_key IS NOT null
AND p_1 IS NOT null
PRIMARY KEY ((p_1), id, shard_key) WITH CLUSTERING ORDER BY (id ASC, shard_key ASC);


