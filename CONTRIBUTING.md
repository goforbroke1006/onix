# onix

### Before writing code

```shell
bash ./setup.sh
```

### Setup local environment

```shell
docker-compose down --volumes --remove-orphans
docker-compose build
docker-compose up -d

echo 'Open dev build http://localhost:3002/'
echo 'Open prod build http://localhost:3001/'

make test/functional
make test/integration

```

### Make and apply local DB dump

```shell
export PGPASSWORD=onix
pg_dump --host=localhost --port=5432 --username=onix --data-only -T public.schema_migrations onix > dump-data.tmp.sql
rm -f dump-data.sql
mv dump-data.tmp.sql dump-data.sql
```

```shell
export PGPASSWORD=onix
psql --host=localhost --port=5432 --username=onix onix < dump-data.sql
```


### Debug dashboard-main API

```shell
curl -X GET 'http://127.0.0.1:8082/api/dashboard-main/service'
curl -X GET 'http://127.0.0.1:8082/api/dashboard-main/release?service=<service>'
curl -X GET 'http://127.0.0.1:8082/api/dashboard-main/compare?service=<service>&release_one_title=1.17.1&release_one_start=1643894400&release_two_title=1.19.0&release_two_start=1643894940&period=1h'
```

### Create migration

```shell
dbmate -d "./migrations" new <migration_name>
```
