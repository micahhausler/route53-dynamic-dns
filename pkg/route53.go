package pkg

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/pkg/errors"
	"inet.af/netaddr"
)

type Updater interface {
	Update(config Config, ips []netaddr.IP, dryRun bool) error
}

type Route53Updater struct {
	client        *route53.Route53
	defaultZoneId string
	defaultTtl    int64
}

type Route53UpdaterOpt func(*Route53Updater)

// SetDefaultTTL sets the TTL on a Route53Updater
func SetDefaultTTL(ttl int64) Route53UpdaterOpt {
	return func(u *Route53Updater) {
		if ttl > 0 {
			u.defaultTtl = ttl
		}
	}
}

// NewRoute53Updater initializes a Route53Updater
func NewRoute53Updater(defaultZoneId string, opts ...Route53UpdaterOpt) (Updater, error) {
	sess := session.Must(session.NewSession())
	svc := route53.New(sess)

	resp := &Route53Updater{
		defaultZoneId: defaultZoneId,
		client:        svc,
		defaultTtl:    3600,
	}
	for _, opt := range opts {
		opt(resp)
	}

	return resp, nil
}

// Update upserts a config's record with the requested IPs
func (u *Route53Updater) Update(config Config, ips []netaddr.IP, dryRun bool) error {
	zoneId := u.defaultZoneId
	if config.ZoneID != "" {
		zoneId = config.ZoneID
	}
	ttl := u.defaultTtl
	if config.TTL > 0 {
		ttl = config.TTL
	}
	changes := []*route53.Change{}
	for _, record := range config.Records {
		ipv4s := []*route53.ResourceRecord{}
		ipv6s := []*route53.ResourceRecord{}
		for _, ip := range ips {
			if ip.Is4() {
				ipv4s = append(ipv4s, &route53.ResourceRecord{Value: aws.String(ip.String())})
			}
			if ip.Is6() {
				ipv6s = append(ipv6s, &route53.ResourceRecord{Value: aws.String(ip.String())})
			}
		}

		if len(ipv4s) > 0 && config.ARecordAllowed() {
			changes = append(changes, &route53.Change{
				Action: aws.String(route53.ChangeActionUpsert),
				ResourceRecordSet: &route53.ResourceRecordSet{
					Name:            aws.String(record),
					ResourceRecords: ipv4s,
					TTL:             aws.Int64(ttl),
					Type:            aws.String(route53.RRTypeA),
				},
			})
		}
		if len(ipv6s) > 0 && config.AAAARecordAllowed() {
			changes = append(changes, &route53.Change{
				Action: aws.String(route53.ChangeActionUpsert),
				ResourceRecordSet: &route53.ResourceRecordSet{
					Name:            aws.String(record),
					ResourceRecords: ipv6s,
					TTL:             aws.Int64(ttl),
					Type:            aws.String(route53.RRTypeAaaa),
				},
			})
		}
	}

	if len(changes) == 0 {
		return nil
	}
	output, err := json.MarshalIndent(changes, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("Submitting change:\n%s \n", string(output))
	if dryRun {
		return nil
	}

	_, err = u.client.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: changes,
		},
		HostedZoneId: aws.String(zoneId),
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
