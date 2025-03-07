package elb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/padok-team/yatas-aws/logger"
)

type LoadBalancerAttributes struct {
	LoadBalancerArn  string
	LoadBalancerName string
	LoadBalancerType types.LoadBalancerTypeEnum
	Listeners        []types.Listener
	Output           *elasticloadbalancingv2.DescribeLoadBalancerAttributesOutput
}

func GetLoadBalancersAttributes(s aws.Config, loadbalancers []types.LoadBalancer) []LoadBalancerAttributes {
	svc := elasticloadbalancingv2.NewFromConfig(s)
	var loadBalancerAttributes []LoadBalancerAttributes
	for _, loadbalancer := range loadbalancers {
		describeLoadBalancerAttributesInput := &elasticloadbalancingv2.DescribeLoadBalancerAttributesInput{
			LoadBalancerArn: loadbalancer.LoadBalancerArn,
		}
		resultDescribeLoadBalancerAttributes, err := svc.DescribeLoadBalancerAttributes(context.TODO(), describeLoadBalancerAttributesInput)
		if err != nil {
			logger.Logger.Error(err.Error())
			// return empty struct
			return []LoadBalancerAttributes{}
		}

		var listeners []types.Listener
		describeListenersInput := &elasticloadbalancingv2.DescribeListenersInput{
			LoadBalancerArn: loadbalancer.LoadBalancerArn,
		}
		resultDescribeListeners, err := svc.DescribeListeners(context.TODO(), describeListenersInput)
		if err != nil {
			logger.Logger.Error(err.Error())
			// return empty struct
			return []LoadBalancerAttributes{}
		}
		listeners = resultDescribeListeners.Listeners
		loadBalancerAttributes = append(loadBalancerAttributes, LoadBalancerAttributes{
			LoadBalancerArn:  *loadbalancer.LoadBalancerArn,
			LoadBalancerName: *loadbalancer.LoadBalancerName,
			LoadBalancerType: loadbalancer.Type,
			Listeners:        listeners,
			Output:           resultDescribeLoadBalancerAttributes,
		})
	}
	return loadBalancerAttributes
}

func GetElasticLoadBalancers(s aws.Config) []types.LoadBalancer {
	svc := elasticloadbalancingv2.NewFromConfig(s)
	var loadBalancers []types.LoadBalancer
	input := &elasticloadbalancingv2.DescribeLoadBalancersInput{
		PageSize: aws.Int32(100),
	}
	result, err := svc.DescribeLoadBalancers(context.TODO(), input)
	if err != nil {
		logger.Logger.Error(err.Error())
		// return empty struct
		return []types.LoadBalancer{}
	}
	loadBalancers = append(loadBalancers, result.LoadBalancers...)
	for {
		if result.NextMarker != nil {
			input.Marker = result.NextMarker
			result, err = svc.DescribeLoadBalancers(context.TODO(), input)
			if err != nil {
				logger.Logger.Error(err.Error())
				// return empty struct
				return []types.LoadBalancer{}
			}
			loadBalancers = append(loadBalancers, result.LoadBalancers...)
		} else {
			break
		}
	}
	return loadBalancers
}
