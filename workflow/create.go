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
	hostPort := GetEnv("TEMPORAL_URL", "localhost:7233")
	log.Println(hostPort)
	temporalClient, err := client.Dial(client.Options{
		HostPort: hostPort,
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

	spec := client.ScheduleSpec{
		Calendars: []client.ScheduleCalendarSpec{
			{DayOfWeek: []client.ScheduleRange{{Start: 0, End: 6}}},
		},
		Jitter:       jitterDuration,
		TimeZoneName: "US/Pacific",
	}
	// Create the schedule.
	_, err = temporalClient.ScheduleClient().Create(ctx, client.ScheduleOptions{
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
	if err != nil && err.Error() != "schedule with this ID is already registered" {
		log.Fatalln("Unable to create schedule", err)
	}
	log.Println("Schedule created", "ScheduleID", scheduleID)
}
