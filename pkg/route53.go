package pkg

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/pkg/errors"
	"inet.af/netaddr"
)

type Updater interface {
	Update(record string, ips []netaddr.IP) error
}

type updater struct {
	client *route53.Route53
	zoneId string
	ttl    int
}

func NewUpdater(zoneId string, ttl int) (Updater, error) {
	sess := session.Must(session.NewSession())
	svc := route53.New(sess)

	return &updater{
		client: svc,
		zoneId: zoneId,
		ttl:    ttl,
	}, nil
}

func (u *updater) Update(record string, ips []netaddr.IP) error {
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
	changes := []*route53.Change{}
	if len(ipv4s) > 0 {
		changes = append(changes, &route53.Change{
			Action: aws.String(route53.ChangeActionUpsert),
			ResourceRecordSet: &route53.ResourceRecordSet{
				Name:            aws.String(record),
				ResourceRecords: ipv4s,
				TTL:             aws.Int64(int64(u.ttl)),
				Type:            aws.String(route53.RRTypeA),
			},
		})
	}
	if len(ipv6s) > 0 {
		changes = append(changes, &route53.Change{
			Action: aws.String(route53.ChangeActionUpsert),
			ResourceRecordSet: &route53.ResourceRecordSet{
				Name:            aws.String(record),
				ResourceRecords: ipv6s,
				TTL:             aws.Int64(int64(u.ttl)),
				Type:            aws.String(route53.RRTypeAaaa),
			},
		})
	}
	if len(changes) == 0 {
		return nil
	}

	_, err := u.client.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: changes,
		},
		HostedZoneId: aws.String(u.zoneId),
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
