package osm_extractor_workflow

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func Worker() {
	hostPort := GetEnv("TEMPORAL_URL", "localhost:7233")
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: hostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "osm-extractor", worker.Options{})

	w.RegisterWorkflow(OsmExtractor)
	w.RegisterActivity(ExtractOsmCutoutsActivity)
	w.RegisterActivity(UploadOsmCutoutsActivity)
	w.RegisterActivity(CopyOsmCutouts)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
