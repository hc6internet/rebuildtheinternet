kubectl run db0 --image=database:v1 --port=3306 --image-pull-policy=Never \
    --env="MYSQL_USER=mysql" \
    --env="MYSQL_PASSWORD=passwd" \
    --env="MYSQL_DATABASE=analytics" 
kubectl expose deployment/db0 

kubectl run db1 --image=database:v1 --port=3306 --image-pull-policy=Never \
    --env="MYSQL_USER=mysql" \
    --env="MYSQL_PASSWORD=passwd" \
    --env="MYSQL_DATABASE=analytics" 
kubectl expose deployment/db1
