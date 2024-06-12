import { proxyActivities } from '@temporalio/workflow';
import type * as activities from './activities.js';

const { extractOsmCutouts, uploadOsmCutouts } = proxyActivities<typeof activities>({
  startToCloseTimeout: '1 minute',
});

export async function extract(): Promise<void> {
  await extractOsmCutouts();
  await uploadOsmCutouts();
}
