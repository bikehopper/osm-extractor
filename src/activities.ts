import { log, Context } from '@temporalio/activity';
import { resolve, join } from 'path'
import { S3Client, CopyObjectCommand } from '@aws-sdk/client-s3';
import { Upload } from '@aws-sdk/lib-storage';
import config from './config.json';
import { createReadStream } from 'node:fs';

const outputDir = process.env.OUTPUT_DIR || '/mnt/output';
const pbfPath = process.env.PBF_PATH || '/mnt/input/latest.osm.pbf';
const bucket = process.env.BUCKET || 'osm-extracts';
const s3Region = process.env.S3_REGION || 'us-east-1';
const s3EndpointUrl = process.env.S3_ENDPOINT_URL || '';

const s3Client = new S3Client({ 
  endpoint: s3EndpointUrl,
  region: s3Region,
  credentials: {
    accessKeyId: process.env.ACCESS_KEY_ID || '',
    secretAccessKey: process.env.SECRET_ACCESS_KEY || '',
  },
  logger: log,
  maxAttempts: 3,
  forcePathStyle: true
});

export async function extractOsmCutouts(): Promise<void> {
  log.info('extracting OSM cutouts');
  const {execa} = await import('execa');
  // not try/catching because I want errors to fail the workflow
  await execa({
    all: true
  })`osmium extract -d ${outputDir} -c /app/lib/config.json ${pbfPath} --overwrite`;
  
  log.info('finished extracting.');
}

export async function uploadOsmCutouts(): Promise<void> {
  log.info('uploading extracts to bucket');
  const ctx = Context.current();
  const scheduledDate = new Date(ctx.info.scheduledTimestampMs);
  
  for (const extract of config.extracts) {
    log.info(`uploading ${extract.output}`);
    const readable = createReadStream(resolve(outputDir, extract.output));

    // upload the date stamped file
    const destPathDated = join(extract.directory, getDatedFileName(extract.output, scheduledDate));
    const datedObjectInput = {
      Bucket: bucket,
      Body: readable,
      Key: destPathDated
    };
    const uploadDatedFile = new Upload({
      client: s3Client,
      params: datedObjectInput,
    });
    await uploadDatedFile.done();
  }

  log.info(`finished uploading extracts`);
}

export async function copyOsmCutouts() {
  log.info('Creating "latest" copies of extracts');
  const ctx = Context.current();
  const scheduledDate = new Date(ctx.info.scheduledTimestampMs);

  for (const extract of config.extracts) {
    const destPathDated = join(extract.directory, getDatedFileName(extract.output, scheduledDate));
    const destPathLatest = join(extract.directory, getLatestFileName(extract.output));
    const latestObjectInput = {
      Bucket: bucket,
      CopySource: join(bucket, destPathDated),
      Key: destPathLatest
    }
    await s3Client.send(new CopyObjectCommand(latestObjectInput));
  }
  log.info('Created "latest" copies of extracts.');
}

function getDatedFileName(fileName:string, date=new Date()): string {
  // YYYY-mm-dd
  const dateString = [date.getUTCFullYear(),date.getUTCMonth()+1,date.getUTCDate()].join('-');
  return `${dateString}-${fileName}`;
}

function getLatestFileName(fileName:string): string {
  return `latest-${fileName}`;
}