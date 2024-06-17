package osm_extractor_workflow

import (
	"context"
	"log"

	"go.temporal.io/sdk/client"
)

// OsmExtractor
func Resume() {
	ctx := context.Background()
	hostPort := GetEnv("TEMPORAL_URL", "localhost:7233")
	temporalClient, err := client.Dial(client.Options{
		HostPort: hostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal Client", err)
	}
	defer temporalClient.Close()

	// Create Schedule and Workflow IDs
	scheduleID := "extract-osm-cutouts-schedule"

	// Create the schedule.
	scheduleHandle := temporalClient.ScheduleClient().GetHandle(ctx, scheduleID)
	err = scheduleHandle.Unpause(ctx, client.ScheduleUnpauseOptions{
		Note: "The Schedule resumed via CLI",
	})
	if err != nil {
		log.Fatalln("Unable to resume schedule", err)
	}
	log.Println("Schedule resumed", "ScheduleID", scheduleID)
	_, _ = scheduleHandle.Describe(ctx)
}
