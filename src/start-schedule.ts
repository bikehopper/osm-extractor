import { Connection, Client, ScheduleOverlapPolicy } from '@temporalio/client';
import { extract } from './workflows.js';

async function run() {
  const client = new Client({
    connection: await Connection.connect(),
  });

  // https://typescript.temporal.io/api/classes/client.ScheduleClient#create
  const schedule = await client.schedule.create({
    action: {
      type: 'startWorkflow',
      workflowType: extract,
      args: [],
      taskQueue: 'schedules',
    },
    scheduleId: 'extract-osm-cutouts-schedule',
    policies: {
      catchupWindow: '1 day',
      overlap: ScheduleOverlapPolicy.ALLOW_ALL,
    },
    spec: {
      intervals: [{ every: '1 day' }],
      jitter: '2m'
    },
  });

  console.log(`Started schedule '${schedule.scheduleId}'.

The reminder Workflow will run the Worker every day.

You can now run:
  npm run schedule.delete
  `);

  await client.connection.close();
}

run().catch((err) => {
  console.error(err);
  process.exit(1);
});
