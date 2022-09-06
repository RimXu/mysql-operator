package constants

const InitCfg = `
if [ ! -d "/data/3306/data" ]; then
    mkdir -p /data/3306/mysql
    mkdir -p /data/3306/tmp
    mkdir -p /data/3306/share
    mkdir -p /data/3306/data
	chown mysql:mysql -R /data/
fi
`
