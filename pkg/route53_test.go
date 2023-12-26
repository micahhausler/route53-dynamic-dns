package pkg

import (
	"net/netip"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	"github.com/google/go-cmp/cmp"
)

type MockR53 struct {
	route53iface.Route53API

	got *route53.ChangeResourceRecordSetsInput
}

func (m *MockR53) ChangeResourceRecordSets(input *route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error) {
	m.got = input
	return nil, nil
}

func TestUpdate(t *testing.T) {
	cases := []struct {
		name     string
		config   Config
		addrs    []netip.Addr
		expected *route53.ChangeResourceRecordSetsInput
	}{
		{
			name: "simple",
			config: Config{
				ZoneID: "testZone",
				Records: []string{
					"test",
				},
				RecordTypes: []RecordType{
					RecordTypeA,
				},
				TTL: 3600,
			},
			addrs: []netip.Addr{netip.MustParseAddr("1.1.1.1")},
			expected: &route53.ChangeResourceRecordSetsInput{
				HostedZoneId: aws.String("testZone"),
				ChangeBatch: &route53.ChangeBatch{
					Changes: []*route53.Change{
						{
							Action: aws.String("UPSERT"),
							ResourceRecordSet: &route53.ResourceRecordSet{
								Name: aws.String("test"),
								Type: aws.String("A"),
								TTL:  aws.Int64(3600),
								ResourceRecords: []*route53.ResourceRecord{
									{
										Value: aws.String("1.1.1.1"),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			updater := &Route53Updater{
				client:     &MockR53{},
				defaultTtl: 3600,
			}
			updater.Update(tc.config, tc.addrs, false)
			got := updater.client.(*MockR53).got
			if got == nil {
				t.Fatalf("got nil")
			}
			if *got.HostedZoneId != *tc.expected.HostedZoneId {
				t.Errorf("expected HostedZoneId %s, got %s", *tc.expected.HostedZoneId, *got.HostedZoneId)
			}
			if len(got.ChangeBatch.Changes) != len(tc.expected.ChangeBatch.Changes) {
				t.Errorf("expected %d changes, got %d", len(tc.expected.ChangeBatch.Changes), len(got.ChangeBatch.Changes))
			}

			if diff := cmp.Diff(tc.expected, got); diff != "" {
				t.Errorf("unexpected records (-want +got):\n%s", diff)
			}

		})
	}

}
