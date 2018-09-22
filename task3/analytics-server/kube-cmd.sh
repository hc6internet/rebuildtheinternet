kubectl run server0 --image=analytics-server:v1 --port=80 --image-pull-policy=Never \
    --env="SERVER_ID=0" \
    --env="PORT=80" \
    --env="PARTITION_SERVICE=partition-select:9000" \
    --env="DB_USER=mysql" \
    --env="DB_PASSWD=passwd" \
    --env="DB_DATABASE=analytics"
kubectl expose deployment/server0 --type=LoadBalancer

kubectl run server1 --image=analytics-server:v1 --port=80 --image-pull-policy=Never \
    --env="SERVER_ID=1" \
    --env="PORT=80" \
    --env="PARTITION_SERVICE=partition-select:9000" \
    --env="DB_USER=mysql" \
    --env="DB_PASSWD=passwd" \
    --env="DB_DATABASE=analytics"
kubectl expose deployment/server1 --type=LoadBalancer
