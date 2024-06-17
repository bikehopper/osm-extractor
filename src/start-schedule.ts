import { Connection, Client, ScheduleOverlapPolicy, ScheduleAlreadyRunning} from '@temporalio/client';
import { extract } from './workflows';

const temporalUrl = process.env.TEMPORAL_URL || 'localhost:7233';

async function run() {
  const client = new Client({
    connection: await Connection.connect({
      address: temporalUrl
    }),
  });

  // https://typescript.temporal.io/api/classes/client.ScheduleClient#create
  try {
    const schedule = await client.schedule.create({
      action: {
        type: 'startWorkflow',
        workflowType: extract,
        args: [],
        taskQueue: 'osm-exttractor',
      },
      scheduleId: 'extract-osm-cutouts-schedule',
      policies: {
        catchupWindow: '1 day',
        overlap: ScheduleOverlapPolicy.CANCEL_OTHER,
      },
      spec: {
        intervals: [{ every: '1 day' }],
        jitter: '2m'
      },
    });
    console.log(`Started schedule '${schedule.scheduleId}'.`);
  } catch (error) {
    if (!(error instanceof ScheduleAlreadyRunning)) {
      throw error;
    }
  } finally {
    await client.connection.close();
  }
}

run().catch((err) => {
  console.error(err);
  process.exit(1);
});
