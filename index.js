var AWS = require('aws-sdk');
var s3 = new AWS.S3();
var http = require('https');


const APIKEY = "3v8qBPX0ePTaOmbm";
const hostURL = "api.songkick.com";
const artistsURL = "/api/3.0/users/{username}/artists/tracked.json?apikey=" + APIKEY;
const concertsURL = "/api/3.0/artists/{artist_id}/calendar.json?apikey=" + APIKEY;

exports.handler = async (event, context, callback) => {
  var bucketName = "concertnotifier";
  var keyName = "users.json"

  return await readFile(bucketName, keyName, readFileContent, onError);
};

async function readFile(bucketName, filename, onFileContent, onError) {
  var params = {
    Bucket: bucketName,
    Key: filename
  };
  var content = "Not null";
  const req = s3.getObject(params);
  var err, data = await req.promise();
  if (!err)
      content = onFileContent(filename, data.Body.toString());
    else {
      console.log(err);
      content = { statusCode: 500, body: "Couldn't read users file" }
    }
  return content;
}

async function readFileContent(filename, content) {
  const users = JSON.parse(content);
  var response = { statusCode: 200, body: []};
  for (var user in users) {
    let userConcerts = await parseUser(users[user].user, users[user].email);
    response.body = response.body.concat({"user": users[user].user, "concerts": userConcerts});
  }
  return response;
}

async function parseUser(username, email) {
  const artists = await getAllArtists(username);
  var concerts = [];
  //for (var artist in artists) {
  for (var artist = 0; artist <= 10; artist++) {
    concerts = concerts.concat(await getAllConcerts(artists[artist]));
  }
  return concerts;
}

async function getAllArtists(username) {
  let userUrl = artistsURL.replace("{username}", username);
  let artists = await get({"method": "GET", "host": hostURL, "path": userUrl});
  artists = artists["resultsPage"]["results"]["artist"];
  let artistsID = [];
  for (var artist in artists) {
    artistsID = artistsID.concat(artists[artist]["id"]);
  }
  return artistsID;
}

async function getAllConcerts(artist) {
  let artistUrl = concertsURL.replace("{artist_id}", artist)
  let concerts = await get({"method": "GET", "host": hostURL, "path": artistUrl});
  concerts = concerts["resultsPage"]["results"]["event"];
  return concerts;
}

function onError(err) {
  console.log('error: ' + err);
}

function get(options) {
  return new Promise(((resolve, reject) => {
    const request = http.request(options, (response) => {
      response.setEncoding('utf8');
      let returnData = '';

      if (response.statusCode < 200 || response.statusCode >= 300) {
        return reject(new Error(`${response.statusCode}: ${response.req.getHeader('host')} ${response.req.path}`));
      }

      response.on('data', (chunk) => {
        returnData += chunk;
      });

      response.on('end', () => {
        resolve(JSON.parse(returnData));
      });

      response.on('error', (error) => {
        reject(error);
      });
    });
    // request.write(postData)
    request.end();
  }));
}
