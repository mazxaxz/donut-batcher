# donut-batcher

`docker-compose -f ./docker-compose.resources.yml up`  
wait until rabbit and mongo are alive (without implementing healthchecks "depends_on" does not work)  

`docker-compose -f ./docker-compose.app.yml up --build`

`rest.http` for app testing
