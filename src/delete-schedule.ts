import { Connection, Client } from '@temporalio/client';

const temporalUrl = process.env.TEMPORAL_URL || 'localhost:7233';

async function run() {
  const client = new Client({
    connection: await Connection.connect({
      address: temporalUrl
    }),
  });

  const handle = client.schedule.getHandle('extract-osm-cutouts-schedule');
  await handle.delete();

  console.log(`Schedule is now deleted.`);
}

run().catch((err) => {
  console.error(err);
  process.exit(1);
});
