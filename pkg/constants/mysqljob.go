package constants

const MySQLJob = `
#!/bin/sh
db_name=$1
db_user=$2
db_pass=$3


func_repl() {

    mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $m_mysql -e  "CREATE USER IF NOT EXISTS '$REPL_USER'@'%' IDENTIFIED BY '$REPL_PASS';"
    mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $s_mysql -e  "CREATE USER IF NOT EXISTS '$REPL_USER'@'%' IDENTIFIED BY '$REPL_PASS';"

    mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $m_mysql -e  "GRANT REPLICATION SLAVE ON *.* to '$REPL_USER'@'%' IDENTIFIED  by '$REPL_PASS';"
    mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $s_mysql -e  "GRANT REPLICATION SLAVE ON *.* to '$REPL_USER'@'%' IDENTIFIED  by '$REPL_PASS';"

    mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $m_mysql -e  "CREATE USER IF NOT EXISTS '$MONITOR_USER'@'%' IDENTIFIED BY '$MONITOR_PASS';"
    mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $s_mysql -e  "CREATE USER IF NOT EXISTS '$MONITOR_USER'@'%' IDENTIFIED BY '$MONITOR_PASS';"

    mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $m_mysql -e  "GRANT SELECT ON *.* TO '$MONITOR_USER'@'%' IDENTIFIED  by '$MONITOR_PASS';"
    mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $s_mysql -e  "GRANT SELECT ON *.* TO '$MONITOR_USER'@'%' IDENTIFIED  by '$MONITOR_PASS';"

    mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $m_mysql -e  "CREATE USER IF NOT EXISTS '$EXPORTER_USER'@'127.0.0.1' IDENTIFIED BY '$EXPORTER_PASS' ;"
    mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $s_mysql -e  "CREATE USER IF NOT EXISTS '$EXPORTER_USER'@'127.0.0.1' IDENTIFIED BY '$EXPORTER_PASS' ;"

    mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $m_mysql -e  "GRANT PROCESS, REPLICATION CLIENT, SELECT ON *.* TO '$EXPORTER_USER'@'%' IDENTIFIED  by '$EXPORTER_PASS';"
    mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $s_mysql -e  "GRANT PROCESS, REPLICATION CLIENT, SELECT ON *.* TO '$EXPORTER_USER'@'%' IDENTIFIED  by '$EXPORTER_PASS';"


    gtid=$(mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $m_mysql -e "show master status \G;" | grep Executed_Gtid_Set|awk '{print $2}')
    mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $s_mysql -e "reset master;stop slave; reset slave all;"
    mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $s_mysql -e "set global gtid_purged='$gtid';"
    mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $s_mysql -e "change master to master_host='$m_mysql',master_user='$REPL_USER',master_password='$REPL_PASS',master_auto_position=1,master_port=3306 ;"
    mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $s_mysql -e "start slave;"

    mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $m_mysql -e  "CREATE USER IF NOT EXISTS '$db_user'@'%' IDENTIFIED BY '$db_pass';"
    #mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $s_mysql -e  "CREATE USER IF NOT EXISTS '$db_user'@'%' IDENTIFIED BY '$db_pass';"

    mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $m_mysql -e  "GRANT ALL PRIVILEGES ON $db_name.* TO '$db_user'@'%' IDENTIFIED BY '$db_pass' WITH GRANT OPTION;"
    #mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $s_mysql -e  "GRANT ALL PRIVILEGES ON $db_name.* TO '$db_user'@'%' IDENTIFIED BY '$db_pass' WITH GRANT OPTION;"

    mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $m_mysql -e "CREATE DATABASE IF NOT EXISTS $db_name;"



    echo $db_name,$db_user,$db_pass,$MONITOR_USER,$MONITOR_PASS,$EXPORTER_USER,$EXPORTER_PASS
}


func_app() {
    mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $m_mysql -e  "CREATE USER IF NOT EXISTS '$db_user'@'%' IDENTIFIED BY '$db_pass';"
    mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $m_mysql -e  "GRANT ALL PRIVILEGES ON $db_name.* TO '$db_user'@'%' IDENTIFIED BY '$db_pass' WITH GRANT OPTION;"
    mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $m_mysql -e "CREATE DATABASE IF NOT EXISTS $db_name;"
}


while true
do
    sleep 10
    m_conn=$(mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $m_mysql --connect-timeout=1 -e "show databases")
    m_status=$?
    s_conn=$(mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $s_mysql --connect-timeout=1 -e "show databases")
    s_status=$?
    echo INFO: MySQL-Master status:$m_status,MySQL-Slave status:$s_status
    if [ $m_status -eq 0 -a $s_status -eq 0 ];then
        echo "INFO: MySQL connects successful"
        repl_conn=$(mysql -uroot -p$MYSQL_ROOT_PASSWORD -h $m_mysql --connect-timeout=1 -e "select user,host from mysql.user;")
        repl_exists=$(echo $repl_conn|grep $REPL_USER|wc -l)
        if [ $repl_exists -eq 0 ];then
            echo "INFO: MySQL replication config"
            func_repl
            break
        else
            m_exists=$(echo $m_conn|grep $db_name|wc -l)
            s_exists=$(echo $s_conn|grep $db_name|wc -l)
            if [ $m_exists -gt 0 -o $s_exists -gt 0 ];then
                echo "WARN: MySQL databases already exists,exit 0"
                exit 0
                break
            else
                echo "INFO: MySQL database add"
                func_app
                break
            fi
        fi
    fi
done
echo "INFO: MySQL connects check break"

`
