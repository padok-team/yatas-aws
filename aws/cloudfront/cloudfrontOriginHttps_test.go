package cloudfront

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/padok-team/yatas/plugins/commons"
)

func TestCheckIfOriginHTTPSTLS12(t *testing.T) {
	custom := func(policy types.OriginProtocolPolicy, protocols ...types.SslProtocol) types.Origin {
		return types.Origin{
			Id: aws.String("alb"),
			CustomOriginConfig: &types.CustomOriginConfig{
				OriginProtocolPolicy: policy,
				OriginSslProtocols:   &types.OriginSslProtocols{Items: protocols},
			},
		}
	}
	s3 := types.Origin{Id: aws.String("s3"), S3OriginConfig: &types.S3OriginConfig{}}
	dist := func(viewer types.ViewerProtocolPolicy, origins ...types.Origin) types.DistributionSummary {
		return types.DistributionSummary{
			Id:                   aws.String("test"),
			DefaultCacheBehavior: &types.DefaultCacheBehavior{ViewerProtocolPolicy: viewer},
			Origins:              &types.Origins{Items: origins},
		}
	}
	redirect := types.ViewerProtocolPolicyRedirectToHttps

	tests := []struct {
		name string
		d    types.DistributionSummary
		want string
	}{
		{"https-only TLS1.2", dist(redirect, custom(types.OriginProtocolPolicyHttpsOnly, types.SslProtocolTLSv12)), "OK"},
		{"match-viewer with HTTPS viewers", dist(redirect, custom(types.OriginProtocolPolicyMatchViewer, types.SslProtocolTLSv12)), "OK"},
		{"S3 origin only", dist(redirect, s3), "OK"},
		{"http-only", dist(redirect, custom(types.OriginProtocolPolicyHttpOnly, types.SslProtocolTLSv12)), "FAIL"},
		{"match-viewer with HTTP viewers", dist(types.ViewerProtocolPolicyAllowAll, custom(types.OriginProtocolPolicyMatchViewer, types.SslProtocolTLSv12)), "FAIL"},
		{"https-only allowing TLSv1.1", dist(redirect, custom(types.OriginProtocolPolicyHttpsOnly, types.SslProtocolTLSv11, types.SslProtocolTLSv12)), "FAIL"},
		{"one bad origin among good ones", dist(redirect, s3, custom(types.OriginProtocolPolicyHttpsOnly, types.SslProtocolTLSv12), custom(types.OriginProtocolPolicyHttpOnly)), "FAIL"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := commons.CheckConfig{Queue: make(chan commons.Check, 1)}
			CheckIfOriginHTTPSTLS12(cc, []types.DistributionSummary{tt.d}, "AWS_CFT_006")
			if got := (<-cc.Queue).Results[0]; got.Status != tt.want {
				t.Errorf("CheckIfOriginHTTPSTLS12() = %v, want %s", got, tt.want)
			}
		})
	}

	// a behavior other than the default accepting HTTP also breaks match-viewer
	d := dist(redirect, custom(types.OriginProtocolPolicyMatchViewer, types.SslProtocolTLSv12))
	d.CacheBehaviors = &types.CacheBehaviors{Items: []types.CacheBehavior{{ViewerProtocolPolicy: types.ViewerProtocolPolicyAllowAll}}}
	cc := commons.CheckConfig{Queue: make(chan commons.Check, 1)}
	CheckIfOriginHTTPSTLS12(cc, []types.DistributionSummary{d}, "AWS_CFT_006")
	if got := (<-cc.Queue).Results[0]; got.Status != "FAIL" {
		t.Errorf("match-viewer with HTTP cache behavior = %v, want FAIL", got)
	}
}
