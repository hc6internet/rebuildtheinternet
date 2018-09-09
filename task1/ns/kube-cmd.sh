kubectl run pdns --image=pdns:v1 --image-pull-policy=Never --port=8000 \
    --env="PDNS_gmysql_host=mysql" \
    --env="PDNS_gmysql_port=3306" \
    --env="PDNS_gmysql_user=mysql" \
    --env="PDNS_gmysql_password=passwd" \
    --env="PDNS_gmysql_dbname=ns" \
    --env="PDNS_webserver_address=0.0.0.0" \
    --env="PDNS_webserver_port=8000" \
    --env="PDNS_webserver_password=passwd" \
    --env="PDNS_api_key=powerdns"
