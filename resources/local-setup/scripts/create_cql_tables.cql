CREATE KEYSPACE video_views_keyspace WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'}
                                      AND durable_writes = true;

CREATE TABLE video_views_keyspace.video_views
(
    video_id         text,
    timestamp_in_min int,
    aggregate_views  int,
    PRIMARY KEY (video_id, timestamp_in_min)
) WITH CLUSTERING ORDER BY (timestamp_in_min DESC);
