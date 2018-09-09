kubectl run mysql --image=mysql:v1 --port=3306 --image-pull-policy=Never \
    --env="MYSQL_USER=mysql" \
    --env="MYSQL_PASSWORD=passwd" \
    --env="MYSQL_DATABASE=ns"
