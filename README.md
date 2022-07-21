# Dynamic Route53 DNS

Dynamically update DNS records with your interface's public IPs

```
Usage of /usr/bin/route53-dynamic-dns:
      --config string         Path to a config JSON file
      --default-ttl int       Default time to live on record
  -z, --default-zone string   Default Zone ID to update
  -i, --interfaces strings    List of interfaces to evaluate (default [eth0])
```

## Config file format

```json
[
  {
    "zoneId": "<hosted zone ID>", // Optional
    "records": [
      "vpn.home.example.com",
      "ui.home.example.com"
    ],
    "recordTypes": ["A", "AAAA"],
    "ttl": 300
  },
  {
    "zoneId": "<hosted zone ID>", // Optional
    "records": [
      "vpn6.home.example.com"
    ],
    "recordTypes": ["AAAA"]
  }
]
```

## Installation on Unifi Dream Machine

First, set up `udm-utilities` from https://github.com/boostchicken/udm-utilities on your Unifi Dream Machine. This utilitiy is very useful with [kchristensen/udm-le](https://github.com/kchristensen/udm-le/) for Lets Encrypt TLS and [wireguard-go](https://github.com/boostchicken/udm-utilities/tree/master/wireguard-go) in udm-utilities.

1. Copy the contents of this repo to your device at `/mnt/data/route53-dynamic-dns`.
   ```sh
   docker run -it --rm -v /mnt/data/:/mnt/data/ --net host alpine /bin/sh
   $ apk -U add git
   $ cd /mnt/data
   $ git clone https://github.com/micahhausler/route53-dynamic-dns.git
   ```
2. Edit `route53-dynamic-dns.env` and tweak variables to meet your needs.
3. Edit `records.json` with the records you want created
4. Run `/mnt/data/route53-dynamic-dns/route53-dynamic-dns.sh sync`.
   This will handle your initial DNS record and setup a cron task at `/etc/cron.d/route53-dynamic-dns` to attempt a DNS update each morning at 03:02.
5. Copy `on_boot.d/21-dynamic-dns.sh` to `/mnt/data/on_boot.d/`.
   This will ensure that the DNS gets updated and cron is re-created after a system update.

## Required IAM permissions

This application needs the same IAM permissions as [kchristensen/udm-le](https://github.com/kchristensen/udm-le/).
I ended up just creating an IAM user with the following permissions and shoving credentials in `/mnt/data/udm-le/.secrets/credentials` along with a corresponding profile in `/mnt/data/udm-le/.secrets/config`

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "",
            "Effect": "Allow",
            "Action": [
                "route53:GetChange",
                "route53:ListResourceRecordSets"
            ],
            "Resource": [
                "arn:aws:route53:::hostedzone/*",
                "arn:aws:route53:::change/*"
            ]
        },
        {
            "Sid": "",
            "Effect": "Allow",
            "Action": [
                "route53:ChangeResourceRecordSets"
            ],
            "Resource": [
                "arn:aws:route53:::hostedzone/<YOURZONEID>"
            ]
        },
        {
            "Sid": "",
            "Effect": "Allow",
            "Action": "route53:ListHostedZonesByName",
            "Resource": "*"
        }
    ]
}
```

## LICENSE

Apache 2.0. See [LICENSE](LICENSE)
