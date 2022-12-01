package scheduler

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func findScheduleByGroupAndName(ctx context.Context, conn *scheduler.Client, groupName, scheduleName string) (*scheduler.GetScheduleOutput, error) {
	in := &scheduler.GetScheduleInput{
		Name:      aws.String(scheduleName),
		GroupName: aws.String(groupName),
	}
	out, err := conn.GetSchedule(ctx, in)
	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			return nil, &resource.NotFoundError{
				LastError:   err,
				LastRequest: in,
			}
		}

		return nil, err
	}

	if out == nil || out.Arn == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out, nil
}

func findScheduleGroupByName(ctx context.Context, conn *scheduler.Client, name string) (*scheduler.GetScheduleGroupOutput, error) {
	in := &scheduler.GetScheduleGroupInput{
		Name: aws.String(name),
	}
	out, err := conn.GetScheduleGroup(ctx, in)
	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			return nil, &resource.NotFoundError{
				LastError:   err,
				LastRequest: in,
			}
		}

		return nil, err
	}

	if out == nil || out.Arn == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out, nil
}
