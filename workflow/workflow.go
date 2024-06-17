package osm_extractor

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

// ScheduleWorkflow executes on the given schedule
func OsmExtractor(ctx workflow.Context) error {

	workflow.GetLogger(ctx).Info("Schedule workflow started.", "StartTime", workflow.Now(ctx))

	ao := workflow.ActivityOptions{
		TaskQueue:           "osm-extractor",
		StartToCloseTimeout: 30 * time.Second,
	}
	ctx1 := workflow.WithActivityOptions(ctx, ao)

	err := workflow.ExecuteActivity(ctx1, ExtractOsmCutoutsActivity).Get(ctx, nil)
	if err != nil {
		workflow.GetLogger(ctx).Error("Executing activity 'ExtractOsmCutoutsActivity' failed", "Error", err)
		return err
	}

	err = workflow.ExecuteActivity(ctx1, UploadOsmCutoutsActivity).Get(ctx, nil)
	if err != nil {
		workflow.GetLogger(ctx).Error("Executing activity 'UploadOsmCutoutsActivity' failed", "Error", err)
		return err
	}

	err = workflow.ExecuteActivity(ctx1, CopyOsmCutouts).Get(ctx, nil)
	if err != nil {
		workflow.GetLogger(ctx).Error("Executing activity 'CopyOsmCutouts' failed", "Error", err)
		return err
	}

	return nil
}
