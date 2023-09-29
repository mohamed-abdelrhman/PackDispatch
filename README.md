# PackDispatch
This README provides instructions on how to use the calculatePacks API endpoint provided by PackDispatch.


## Overview
The calculatePacks endpoint calculates the number and sizes of packs needed to fulfill a given order quantity. It accepts a JSON payload with the quantity field, which represents the quantity of the order.

## API Call
`curl --location --request POST 'http://161.35.244.35/api/v1/orders/calculate-packs' 
--header 'Content-Type: application/json'
--data '{
"quantity":251
}'`


## Run the application Locally
 * clone the repository
` git clone git@github.com:mohamed-abdelrhman/pack-dispatch.git`
 * build and run the application `docker-compose up`

## Deploy on K8 cluster

* `kubectl apply -f psql.deployment.yaml`
* `kubectl apply -f psql.service.yaml`
* `kubectl apply -f app.deployment.yaml`
* `kubectl apply -f app.service.yaml`
* `kubectl - get svc -w `