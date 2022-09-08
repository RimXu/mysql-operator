package constants

// 可以自定义my.cnf文件,目前只会替换server_id=MYSQL_SERVER_ID和innodb_buffer_pool_size=MYSQL_BUFFER_POOL_SIZE
// 目前是使用strings.Replace方法替换配置
// 可以修改mysql-operator/controllers/configmaps中ReadMycnf方法
const ProxyCfg = `datadir="/var/lib/proxysql"
errorlog="/var/lib/proxysql/proxysql.log"

admin_variables=
{
        admin_credentials="admin:admin"
        mysql_ifaces="0.0.0.0:6032"
}

mysql_variables=
{
        threads=4
        max_connections=2048
        default_query_delay=0
        default_query_timeout=36000000
        have_compress=true
        poll_timeout=2000
        interfaces="0.0.0.0:6033"
        default_schema="information_schema"
        stacksize=1048576
        server_version="5.5.30"
        connect_timeout_server=3000
        monitor_username="MONITOR-USER"
        monitor_password="MONITOR-PASSWORD"
        monitor_history=600000
        monitor_connect_interval=60000
        monitor_ping_interval=10000
        monitor_read_only_interval=1500
        monitor_read_only_timeout=500
        ping_interval_server_msec=120000
        ping_timeout_server=500
        commands_stats=true
        sessions_sort=true
        connect_retries_on_failure=10
}


# defines all the MySQL servers
mysql_servers =
(
    {
        address ="MYSQL-ADDR-master.NAMESPACE"
        port = 3306
        hostgroup = 1
        weight = 1
        compression = 0
        max_replication_lag = 0
        max_connections = 1000
    },
    {   address ="MYSQL-ADDR-slave.NAMESPACE"
        port = 3306
        hostgroup = 2
        weight = 1
        compression = 0
        max_replication_lag = 0
        max_connections = 1000
    }

)


# defines all the MySQL users
mysql_users:
(
MYSQL_USERS
)





#defines MySQL Query Rules
mysql_query_rules:
(
)

scheduler=
(
)


mysql_replication_hostgroups=
(
    {
        writer_hostgroup = 1
        reader_hostgroup = 2
        check_type = "read_only"
        comment="MySQL HA"
    }
)
`
