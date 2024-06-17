package osm_extractor_workflow

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
)

// OsmExtractor
func Create() {
	ctx := context.Background()
	temporalClient, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal Client", err)
	}
	defer temporalClient.Close()

	// Create Schedule and Workflow IDs
	scheduleID := "extract-osm-cutouts-schedule"
	workflowID := "extract-osm-cutouts"
	catchupDuration, _ := time.ParseDuration("12h")
	jitterDuration, _ := time.ParseDuration("2m")
	// intervalDuration, _ := time.ParseDuration("24h")

	spec := client.ScheduleSpec{
		Calendars: []client.ScheduleCalendarSpec{
			{DayOfWeek: []client.ScheduleRange{{Start: 0, End: 6}}},
		},
		Jitter:       jitterDuration,
		TimeZoneName: "US/Pacific",
	}
	// Create the schedule.
	scheduleHandle, err := temporalClient.ScheduleClient().Create(ctx, client.ScheduleOptions{
		ID:            scheduleID,
		CatchupWindow: catchupDuration,
		Overlap:       enums.SCHEDULE_OVERLAP_POLICY_CANCEL_OTHER,
		Spec:          spec,
		Action: &client.ScheduleWorkflowAction{
			ID:        workflowID + uuid.New().String(),
			Workflow:  OsmExtractor,
			TaskQueue: "osm-extractor",
		},
	})
	if err != nil {
		log.Fatalln("Unable to create schedule", err)
	}
	log.Println("Schedule created", "ScheduleID", scheduleID)
	_, _ = scheduleHandle.Describe(ctx)
}
