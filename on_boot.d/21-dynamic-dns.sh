#!/bin/sh

# Load Environmnt Variables
. /mnt/data/route53-dynamic-dns/route53-dynamic-dns.env

if [ ! -f /etc/cron.d/route53-dynamic-dns ]; then
	# Sleep for 5 minutes to avoid restarting
	# services during system startup.
	sleep 300
	RESTART_SERVICES=true sh ${DYNAMIC_DNS_PATH}/route53-dynamic-dns.sh sync
fi