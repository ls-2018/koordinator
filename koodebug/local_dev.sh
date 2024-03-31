export MANAGER_IMG=acejilam/manager:v1
export SCHEDULER_IMG=acejilam/sc:v1
export DESCHEDULER_IMG=acejilam/desc:v1
export KOORDLET_IMG=acejilam/koordlet:v1

docker build -t ${MANAGER_IMG} -f       koodebug/docker/koord-manager.dockerfile .
docker build -t ${SCHEDULER_IMG} -f     koodebug/docker/koord-scheduler.dockerfile .
docker build -t ${DESCHEDULER_IMG} -f   koodebug/docker/koord-descheduler.dockerfile .
docker build -t ${KOORDLET_IMG} -f      koodebug/docker/koordlet.dockerfile .


kind load docker-image -n koord ${MANAGER_IMG}
kind load docker-image -n koord ${SCHEDULER_IMG}
kind load docker-image -n koord ${DESCHEDULER_IMG}
kind load docker-image -n koord ${KOORDLET_IMG}

./hack/deploy_kind.sh

