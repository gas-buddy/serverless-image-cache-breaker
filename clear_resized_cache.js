const AWS = require('aws-sdk');

const S3 = new AWS.S3({
  signatureVersion: 'v4',
});

const RESIZE_IMAGE_REGEX = /^(\d{1,5}|auto)x(\d{1,5}|auto)\/.*/;
const BUCKET = 'gasbuddy-dynamic-images';

async function deleteResizedImages(fileKey) {
  let marker;
  let result;
  let total = 0;
  let commonPrefixes = [];
  try {
    do {
      const listRequest = {
        Bucket: BUCKET,
        Delimiter: '/',
        Marker: marker,
        MaxKeys: 10,
      }

      result = await S3.listObjects(listRequest).promise();
      commonPrefixes = commonPrefixes.concat(result.CommonPrefixes);
      marker = result.NextMarker;
    } while (result.IsTruncated);
    let i = 0;
    let prefixChunks = [];
    const chunkSize = 50; // how many images to delete at a time
    while (i < commonPrefixes.length) {
      prefixChunks.push(commonPrefixes.slice(i, i += chunkSize));
    }
    await Promise.all(prefixChunks.map(async (prefixes) => {
      const resizedKeys = prefixes
      .map(p => `${p.Prefix}${fileKey}`)
      .filter(p => RESIZE_IMAGE_REGEX.test(p))
      .map(Key => ({ Key }));
      if (resizedKeys.length > 0) {
        const deleteRequest = {
          Bucket: BUCKET,
          Delete: {
            Objects: resizedKeys,
          },
        };
        const result = await S3.deleteObjects(deleteRequest).promise();
        result.Deleted = result.Deleted.length;
        total += result.Deleted;
        console.log(`result = ${JSON.stringify(result)}`);
      }
    }));
    console.log(`Deleted ${total}`);
  } catch (err) {
    console.log(err);
  }
  console.log('Done');
}

(async function () {
  await deleteResizedImages(process.argv[2]);
}());
