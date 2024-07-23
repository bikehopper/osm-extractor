package osm_extractor_workflow

import (
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/workflow"
)

// OsmExtractor executes on the given schedule
func OsmExtractor(ctx workflow.Context) error {
	so := &workflow.SessionOptions{
		CreationTimeout:  time.Minute,
		ExecutionTimeout: 20 * time.Minute,
	}
	sessionCtx, err := workflow.CreateSession(ctx, so)
	if err != nil {
		return err
	}
	defer workflow.CompleteSession(sessionCtx)

	ao := workflow.ActivityOptions{
		TaskQueue:           "osm-extractor",
		StartToCloseTimeout: 10 * time.Second,
	}
	activitySessionCtx := workflow.WithActivityOptions(sessionCtx, ao)

	err = workflow.ExecuteActivity(activitySessionCtx, ExtractOsmCutoutsActivity).Get(ctx, nil)
	if err != nil {
		workflow.GetLogger(ctx).Error("Executing activity 'ExtractOsmCutoutsActivity' failed", "Error", err)
		return err
	}

	err = workflow.ExecuteActivity(activitySessionCtx, UploadOsmCutoutsActivity).Get(ctx, nil)
	if err != nil {
		workflow.GetLogger(ctx).Error("Executing activity 'UploadOsmCutoutsActivity' failed", "Error", err)
		return err
	}

	var copyOsmCutoutsResult LatestExractsObjects
	err = workflow.ExecuteActivity(activitySessionCtx, CopyOsmCutouts).Get(ctx, &copyOsmCutoutsResult)
	if err != nil {
		workflow.GetLogger(ctx).Error("Executing activity 'CopyOsmCutouts' failed", "Error", err)
		return err
	}

	cwo := workflow.ChildWorkflowOptions{
		WorkflowExecutionTimeout: 10 * time.Minute,
		TaskQueue:                "REAPLCE",
		WorkflowID:               "REPLACE" + "-" + uuid.New().String(),
	}
	childCtx := workflow.WithChildOptions(ctx, cwo)

	for _, extract := range copyOsmCutoutsResult.ExtractObjects {
		err = workflow.ExecuteChildWorkflow(childCtx, "REPLACE", extract).Get(childCtx, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
