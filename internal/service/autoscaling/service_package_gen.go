// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package autoscaling

import (
	"context"

	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  DataSourceGroup,
			TypeName: "aws_autoscaling_group",
		},
		{
			Factory:  DataSourceGroups,
			TypeName: "aws_autoscaling_groups",
		},
		{
			Factory:  DataSourceLaunchConfiguration,
			TypeName: "aws_launch_configuration",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourceAttachment,
			TypeName: "aws_autoscaling_attachment",
		},
		{
			Factory:  ResourceGroup,
			TypeName: "aws_autoscaling_group",
		},
		{
			Factory:  ResourceGroupTag,
			TypeName: "aws_autoscaling_group_tag",
		},
		{
			Factory:  ResourceLifecycleHook,
			TypeName: "aws_autoscaling_lifecycle_hook",
		},
		{
			Factory:  ResourceNotification,
			TypeName: "aws_autoscaling_notification",
		},
		{
			Factory:  ResourcePolicy,
			TypeName: "aws_autoscaling_policy",
		},
		{
			Factory:  ResourceSchedule,
			TypeName: "aws_autoscaling_schedule",
		},
		{
			Factory:  ResourceLaunchConfiguration,
			TypeName: "aws_launch_configuration",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.AutoScaling
}

var ServicePackage = &servicePackage{}
