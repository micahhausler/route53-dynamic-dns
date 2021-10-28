set -e

# Load environment variables
. /mnt/data/route53-dynamic-dns/route53-dynamic-dns.env

# Setup nightly cron job
CRON_FILE='/etc/cron.d/rout53-dynamic-dns'
if [ ! -f "${CRON_FILE}" ]; then
	echo "2 3 * * * sh ${DYNAMIC_DNS_PATH}/rout53-dynamic-dns.sh renew" > ${CRON_FILE}
	chmod 644 ${CRON_FILE}
	/etc/init.d/crond reload ${CRON_FILE}
fi

CONTAINER_NAME="route53-dynamic-dns"
DOCKER_VOLUMES="-v ${CREDENTIAL_DIRECTORY_PATH}:/root/.aws/ -v ${DYNAMIC_DNS_PATH}:${DYNAMIC_DNS_PATH}"
PODMAN_CMD="podman run --env-file=${DYNAMIC_DNS_PATH}/route53-dynamic-dns.env -it --name=${CONTAINER_NAME} --network=host --rm ${DOCKER_VOLUMES} ${CONTAINER_IMAGE}:${CONTAINER_IMAGE_TAG}"
CMD_ARGS="-z ${ZONE_ID} -i ${IFACE} --config ${DYNAMIC_DNS_PATH}/records.json"

function sync_dns() {
	$PODMAN_CMD $CMD_ARGS
}

case $1 in
sync)
	echo "Synching DNS"
	sync_dns
	;;
*)
	echo '$(basename "$0") sync'
	;;
esac