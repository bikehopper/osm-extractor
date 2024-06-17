package main

import (
	"log"

	osm_extractor "github.com/bikehopper/osm-extractor/workflow"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "osm-extractor", worker.Options{})

	w.RegisterWorkflow(osm_extractor.OsmExtractor)
	w.RegisterActivity(osm_extractor.ExtractOsmCutoutsActivity)
	w.RegisterActivity(osm_extractor.UploadOsmCutoutsActivity)
	w.RegisterActivity(osm_extractor.CopyOsmCutouts)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
