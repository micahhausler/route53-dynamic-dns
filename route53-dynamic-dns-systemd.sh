#!/bin/bash

set -e

# Load environment variables
set -a
source /data/route53-dynamic-dns/route53-dynamic-dns.env
set +a


ROUTE53_DYNAMIC_DNS_PATH='/data/route53-dynamic-dns'
ROUTE53_DYNAMIC_DNS_BINARY="route53-dynamic-dns"
ROUTE53_DYNAMIC_DNS_BINARY_PATH="${ROUTE53_DYNAMIC_DNS_PATH}/${ROUTE53_DYNAMIC_DNS_BINARY}"


create_services() {
	# Create systemd service and timers (for renewal)
	echo "create_services(): Creating route53-dynamic-dns systemd service and timer"
	cp -f "${ROUTE53_DYNAMIC_DNS_PATH}/resources/systemd/route53-dynamic-dns.service" /etc/systemd/system/route53-dynamic-dns.service
	cp -f "${ROUTE53_DYNAMIC_DNS_PATH}/resources/systemd/route53-dynamic-dns.timer" /etc/systemd/system/route53-dynamic-dns.timer

	systemctl daemon-reload
	systemctl enable route53-dynamic-dns.timer
}

install_binary() {
	if [ ! -f "${ROUTE53_DYNAMIC_DNS_BINARY_PATH}" ]; then
		echo "install_binary(): Downloading route53-dynamic-dns binary"
		wget -qO "/tmp/r53_d_d_release-${ROUTE53_DYNAMIC_DNS_VERSION}.tar.gz" "${ROUTE53_DYNAMIC_DNS_DOWNLOAD_URL}"
		chmod +x "${ROUTE53_DYNAMIC_DNS_BINARY_PATH}"

		echo "install_binary(): Extracting route53-dynamic-dns binary from release and placing at ${ROUTE53_DYNAMIC_DNS_BINARY_PATH}"
		tar -xozvf "/tmp/r53_d_d_release-${ROUTE53_DYNAMIC_DNS_VERSION}.tar.gz" --directory="${ROUTE53_DYNAMIC_DNS_PATH}" route53-dynamic-dns

		# Verify lego binary integrity
		echo "install_binary(): Verifying integrity of route53-dynamic-dns binary"
		R53DD_HASH=$(sha1sum "${ROUTE53_DYNAMIC_DNS_BINARY}" | awk '{print $1}')
		if [ "${R53DD_HASH}" = "${ROUTE53_DYNAMIC_DNS_SHA1}" ]; then
			echo "install_binary(): Verified route53-dynamic-dns v${ROUTE53_DYNAMIC_DNS_VERSION}:${R53DD_HASH}"
			chmod +x "${ROUTE53_DYNAMIC_DNS_BINARY}"
		else
			echo "install_binary(): Verification failure, route53-dynamic-dns binary sha1 was ${R53DD_HASH}, expected ${ROUTE53_DYNAMIC_DNS_SHA1}. Cleaning up and aborting"
			rm -f "${ROUTE53_DYNAMIC_DNS_BINARY_PATH}" "/tmp/r53_d_d_release-${ROUTE53_DYNAMIC_DNS_VERSION}.tar.gz"
			exit 1
		fi
	else
		echo "install_binary(): route53-dynamic-dns binary is already installed at ${ROUTE53_DYNAMIC_DNS_BINARY_PATH}, no operation necessary"
	fi
}

case $1 in
setup)
	echo "Setting up"
	install_binary
	create_services
	;;
*)
	echo '$(basename "$0") setup'
	;;
esac
