import { proxyActivities } from '@temporalio/workflow';
import type * as activities from './activities';

const { extractOsmCutouts, uploadOsmCutouts, copyOsmCutouts } = proxyActivities<typeof activities>({
  startToCloseTimeout: '5 minute',
});

export async function extract(): Promise<void> {
  await extractOsmCutouts();
  await uploadOsmCutouts();
  await copyOsmCutouts();
}
