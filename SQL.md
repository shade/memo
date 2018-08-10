# SQL

This file contains any useful SQL queries,
including any SQL that should be ran manually when upgrading.
No manual SQL should be needed for new instances.


## Manual DB Updates

#### 2018-07-30 - Fix transaction input index

```sql
DROP INDEX `transaction_in_index` ON `transaction_ins`;

SELECT previous_out_point_hash, previous_out_point_index, COUNT(1) AS count
FROM transaction_ins
GROUP BY previous_out_point_hash, previous_out_point_index
ORDER BY count DESC;

ALTER TABLE `transaction_ins` ADD UNIQUE `previous_out`(`previous_out_point_hash`, `previous_out_point_index`);

CREATE TABLE `transaction_ins_early` SELECT * FROM `transaction_ins` WHERE `previous_out_point_index` = 4294967295 AND `previous_out_point_hash` = "";
SELECT COUNT(*) FROM `transaction_ins` WHERE `previous_out_point_index` = 4294967295 AND `previous_out_point_hash` = "";
DELETE FROM `transaction_ins` WHERE `previous_out_point_index` = 4294967295 AND `previous_out_point_hash` = "";
```

#### 2018-08-04 Initial Schema Updates

```sql
ALTER TABLE `transaction_ins` MODIFY `unlock_string` VARCHAR(1000);
ALTER TABLE `transaction_outs` MODIFY `lock_string` VARCHAR(1000);
ALTER TABLE `memo_posts` MODIFY `message` VARCHAR(500);
```

## Useful Queries

#### Memos by hour

```sql
SELECT
    COUNT(*) AS count,
    COUNT(DISTINCT pk_hash) AS users,
    DATE_FORMAT(`timestamp`, '%Y-%m-%d %H:00') AS hour
FROM memo_tests
JOIN blocks ON (memo_tests.block_id = blocks.id)
GROUP BY hour
ORDER BY hour DESC
LIMIT 1000;
```

#### Memos by day

```sql
SELECT
    COUNT(*) AS count,
    COUNT(DISTINCT pk_hash) AS users,
    DATE_FORMAT(`timestamp`, '%Y-%m-%d') AS day
FROM memo_tests
JOIN blocks ON (memo_tests.block_id = blocks.id)
GROUP BY day
ORDER BY day DESC
LIMIT 1000;
```
