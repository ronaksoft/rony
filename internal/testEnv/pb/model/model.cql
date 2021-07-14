

CREATE TABLE tab_model_1 
(
	id  int,
	shard_key  int,
	p_1  blob,
	p_2  blob,
	p_5  bigint,
	enum  int,
	sdata  blob,
	PRIMARY KEY ((id), shard_key, enum)
) WITH CLUSTERING ORDER BY (shard_key ASC, enum ASC);

CREATE MATERIALIZED VIEW view_model_1_customer_sort AS
SELECT *
FROM tab_model_1
WHERE enum IS NOT null
AND shard_key IS NOT null
AND id IS NOT null
PRIMARY KEY ((enum), shard_key, id) WITH CLUSTERING ORDER BY (shard_key ASC, id ASC);




CREATE TABLE tab_model_2 
(
	id  bigint,
	shard_key  int,
	p_1  blob,
	p_2  blob,
	p_5  bigint,
	sdata  blob,
	PRIMARY KEY ((id, shard_key), p_1)
) WITH CLUSTERING ORDER BY (p_1 DESC);

CREATE MATERIALIZED VIEW view_model_2_p_1_shard_key_id AS
SELECT *
FROM tab_model_2
WHERE p_1 IS NOT null
AND shard_key IS NOT null
AND id IS NOT null
PRIMARY KEY ((p_1), shard_key, id) WITH CLUSTERING ORDER BY (shard_key ASC, id ASC);




CREATE TABLE tab_model_3 
(
	id  bigint,
	shard_key  int,
	p_1  blob,
	p_2  blob,
	p_5  blob,
	sdata  blob,
	PRIMARY KEY ((id, shard_key), p_1)
) WITH CLUSTERING ORDER BY (p_1 DESC);

CREATE MATERIALIZED VIEW view_model_3_p_1_shard_key_id AS
SELECT *
FROM tab_model_3
WHERE p_1 IS NOT null
AND shard_key IS NOT null
AND id IS NOT null
PRIMARY KEY ((p_1), shard_key, id) WITH CLUSTERING ORDER BY (shard_key ASC, id ASC);

CREATE MATERIALIZED VIEW view_model_3_p_1_id AS
SELECT *
FROM tab_model_3
WHERE p_1 IS NOT null
AND id IS NOT null
PRIMARY KEY ((p_1), id) WITH CLUSTERING ORDER BY (id ASC);


