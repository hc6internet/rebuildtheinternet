kubectl run partition-select --image=partition-select:v1 --port=9000 --image-pull-policy=Never \
    --env="PORT=9000" \
    --env="PARTITION_COUNT=2"

kubectl expose deployment/partition-select 
