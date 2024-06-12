import { log } from '@temporalio/activity';
import { resolve, join, parse } from 'path'
import { S3Client, PutObjectCommand, CopyObjectCommand } from "@aws-sdk/client-s3";
import config from './config.json';
import { execa } from 'execa';

const outputDir = process.env.OUTPUT_DIR || '/mnt/output';
const pbfPath = process.env.PBF_PATH || '/mnt/input/latest.osm.pbf';
const bucket = process.env.BUCKET || 'osm-extracts';
const s3EndpointUrl = process.env.S3_ENDPOINT_URL || '';

export async function extractOsmCutouts(): Promise<void> {
  log.info('extracting OSM cutouts');
  
  // not try/catching because I want errors to fail the workflow
  await execa({
    all: true
  })`osmium extract -d ${outputDir} -c ./config.json ${pbfPath}`;
  
  log.info('finished extracting.');
}

export async function uploadOsmCutouts(): Promise<void> {
  log.info('uploading extracts to bucket');
  const client = new S3Client({ 
    endpoint: s3EndpointUrl,
    credentials: {
      accessKeyId: process.env.ACCESS_KEY_ID || '',
      secretAccessKey: process.env.SECRET_ACCESS_KEY || '',
    },
    logger: log,
    maxAttempts: 3,
  });
  
  for (const extract of config.extracts) {
    log.info(`uploading ${extract.output}`);
    const srcPath = resolve(outputDir, extract.output);
    
    // upload the date stamped file
    const destPathDated = join(extract.directory, getDatedFileName(extract.output));
    const s3PutObjectInput = {
      Bucket: bucket,
      Body: srcPath,
      Key: destPathDated
    };
    await client.send(new PutObjectCommand(s3PutObjectInput));
    
    // copy the date stamped file to "latest.osm.pbf"
    const destPathLatest = join(extract.directory, getLatestFileName(extract.output));
    const copyObjectInput = {
      Bucket: bucket,
      CopySource: destPathDated,
      Key: destPathLatest
    }
    await client.send(new CopyObjectCommand(copyObjectInput));
  }

  log.info(`finished uploading extracts`);
}

function getDatedFileName(fileName:string, date=new Date()): string {
  const fileNameParts = parse(fileName);
  // YYYY-mm-dd
  const dateString = [date.getUTCFullYear(),date.getUTCMonth()+1,date.getUTCDate()].join('-');
  return `${fileNameParts.name}-${dateString}-${fileNameParts.ext}`;
}

function getLatestFileName(fileName:string): string {
  const fileNameParts = parse(fileName);
  return `${fileNameParts.name}-latest-${fileNameParts.ext}`;
}