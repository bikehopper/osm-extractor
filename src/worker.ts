import { Worker, NativeConnection } from '@temporalio/worker';
import * as activities from './activities';

const temporalUrl = process.env.TEMPORAL_URL || 'localhost:7233';

async function run() {
  const connection = await NativeConnection.connect({
    address: temporalUrl
  });
  const worker = await Worker.create({
    workflowsPath: require.resolve('./workflows'),
    activities,
    taskQueue: 'schedules',
    connection
  });

  await worker.run();
}

run().catch((err) => {
  console.error(err);
  process.exit(1);
});
