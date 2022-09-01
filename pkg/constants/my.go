package constants

// 可以自定义my.cnf文件,目前只会替换server_id=MYSQL_SERVER_ID和innodb_buffer_pool_size=MYSQL_BUFFER_POOL_SIZE
// 目前是使用strings.Replace方法替换配置
// 可以修改mysql-operator/controllers/configmaps中ReadMycnf方法
const MySQLCfg = `[client]
port=3306
socket=/data/3306/tmp/mysql.sock
default-character-set=utf8mb4

[mysql]
port=3306
#prompt=\\u@\\d \\r:\\m:\\s>
default-character-set=utf8mb4
no-auto-rehash

[mysqld]
#server_setting
skip-ssl
port=3306
server_id=MYSQL_SERVER_ID
back_log=1024
slow_query_log=1
long_query_time=3
skip-name-resolve
skip-external-locking
performance_schema=on
character-set-server=utf8mb4
lower_case_table_names=1
log_slow_admin_statements=1
group_concat_max_len=102400
explicit_defaults_for_timestamp=true
transaction-isolation=READ-COMMITTED
sql_mode='ONLY_FULL_GROUP_BY,NO_AUTO_VALUE_ON_ZERO,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION,PIPES_AS_CONCAT,ANSI_QUOTES'
#secure_file_priv=
#binlog_checksum=NONE

#limit_setting
connect_timeout=8
net_read_timeout=30
net_write_timeout=60
open_files_limit=65536
max_connections=4500
max_connect_errors=65536
max_allowed_packet=128M
max_user_connections=4000

#cache_setting
thread_stack=512K
thread_cache_size=256
sort_buffer_size=256K
join_buffer_size=128K
read_buffer_size=128K
read_rnd_buffer_size=128K
table_open_cache=10000
table_definition_cache=3000

#dir_setting
tmpdir=/data/3306/tmp
basedir=/data/3306/mysql
datadir=/data/3306/data
socket=/data/3306/tmp/mysql.sock
pid-file=/data/3306/tmp/mysql.pid
log_bin=/data/3306/binlog
#lc-messages-dir=/data/3306/share
#log-error=/data/3306/error.log
#general_log_file=/data/3306/general.log
#slow_query_log_file=/data/3306/slow.log
innodb_data_home_dir=/data/3306/data
innodb_log_group_home_dir=/data/3306/data
slave_load_tmpdir=/data/3306/tmp
relay_log=/data/3306/relaylog
relay_log_info_file=/data/3306/relay-log.info
relay_log_index=/data/3306/mysqld-relay-bin.index



#myisam_setting
#key_buffer=64M
concurrent_insert=2
delayed_insert_timeout=300
myisam_sort_buffer_size=64M


#innodb_setting
innodb_purge_threads=1
innodb_io_capacity=1000
innodb_open_files=60000
innodb_read_io_threads=4
innodb_write_io_threads=4
innodb_thread_concurrency=8
innodb_flush_sync=1
innodb_doublewrite=1
innodb_flush_log_at_trx_commit=1
innodb_old_blocks_time=1000
innodb_max_dirty_pages_pct=60
innodb_flush_method=O_DIRECT
innodb_change_buffering=all
innodb_buffer_pool_size=MYSQL_BUFFER_POOL_SIZE
innodb_data_file_path=ibdata1:16M:autoextend
innodb_file_per_table=1

#binlog_setting
sync_binlog=1
log-slave-updates=0
expire_logs_days=3
binlog_format=ROW
binlog_cache_size=32K
max_binlog_size=500M
#max_binlog_cache_size=2G

#replication_setting
gtid_mode=on
enforce_gtid_consistency=on
slave_net_timeout=4
sync_master_info=1000
sync_relay_log_info=1000
master-info-repository=table
#rpl_semi_sync_slave_enabled=1
#rpl_semi_sync_master_enabled=1
#rpl_semi_sync_master_timeout=1000
#rpl_semi_sync_master_wait_no_slave=1
slave_type_conversions="ALL_NON_LOSSY"
slave-parallel-type=LOGICAL_CLOCK
slave-parallel-workers=8`
