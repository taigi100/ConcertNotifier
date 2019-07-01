var AWS = require('aws-sdk');
var s3 = new AWS.S3();
var http = require('https');
var ses = new AWS.SES({
  region: 'eu-west-1'
});

const timisoara_lat = 45.7494;
const timisoara_lon = 21.2272;
const max_distance = 2000;

const APIKEY = "3v8qBPX0ePTaOmbm";
const hostURL = "api.songkick.com";
const artistsURL = "/api/3.0/users/{username}/artists/tracked.json?apikey=" + APIKEY;
const concertsURL = "/api/3.0/artists/{artist_id}/calendar.json?apikey=" + APIKEY;

exports.handler = async (event, context, callback) => {
  var bucketName = "concertnotifier";
  var keyName = "users.json"

  var retVal = await readFile(bucketName, keyName, readFileContent, onError);

  // return retVal;
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
    content = {
      statusCode: 500,
      body: "Couldn't read users file"
    }
  }

  return content;
}

async function readFileContent(filename, content) {
  const users = JSON.parse(content);
  var response = {
    statusCode: 200,
    body: []
  };
  for (var user in users) {
    let userConcerts = await parseUser(users[user].user, users[user].email);
    response.body = response.body.concat({
      "user": users[user].user,
      "concerts": userConcerts
    });
  }
  return response;
}

function sortF(a, b) {
  if (a["distance"] == b["distance"]) {
    return a["start"]["numdate"] - b["start"]["numdate"];
  }
  return a["distance"] - b["distance"];
}

async function parseUser(username, email) {
  const artists = await getAllArtists(username);
  var concerts = [];
  for (var artist in artists) {
    // for (var artist = 0; artist <= 10; artist++) {
    let aux = await getAllConcerts(artists[artist]);
    if (aux.length > 0)
      aux.sort(sortF);
    concerts.push(aux.slice(0, 2));
  }
  concerts = flatten(concerts);
  concerts.sort(sortF);
  let email_concerts = {
    "concerts": concerts
  };
  var eParams = {
    Template: "Concert_Notifier4",
    ConfigurationSetName: "test",
    Destination: {
      ToAddresses: [email]
    },
    Source: "taigi100@gmail.com",
    TemplateData: JSON.stringify(email_concerts)
  };

  var err, data = await ses.sendTemplatedEmail(eParams).promise();
  if (err) console.log(err);
  else {
    console.log("===EMAIL SENT===");
    console.log('EMAIL: ', email);

  }

  return concerts;
}

function flatten(arr) {
  return arr.reduce(function(flat, toFlatten) {
    return flat.concat(Array.isArray(toFlatten) ? flatten(toFlatten) : toFlatten);
  }, []);
}

async function getAllArtists(username) {
  let userUrl = artistsURL.replace("{username}", username);
  let artists = await get({
    "method": "GET",
    "host": hostURL,
    "path": userUrl
  });
  artists = artists["resultsPage"]["results"]["artist"];
  let artistsID = [];
  for (var artist in artists) {
    artistsID.push(artists[artist]["id"]);
  }
  return artistsID;
}

async function getAllConcerts(artist) {
  let artistUrl = concertsURL.replace("{artist_id}", artist);
  let concerts = await get({
    "method": "GET",
    "host": hostURL,
    "path": artistUrl
  });
  concerts = concerts["resultsPage"]["results"]["event"];
  let retVal = [];
  for (var concert in concerts) {
    let aux = {};
    aux["displayName"] = concerts[concert]["displayName"];
    aux["uri"] = concerts[concert]["uri"];
    aux["start"] = concerts[concert]["start"];
    aux["start"] = concerts[concert]["start"];
    aux["start"]["numdate"] = Date.parse(concerts[concert]["start"]["date"]);
    aux["location"] = concerts[concert]["location"];
    aux["distance"] = Math.ceil(haversine(timisoara_lat, timisoara_lon, aux["location"]["lat"], aux["location"]["lng"]) / 1000); //in KM
    if (aux["distance"] <= max_distance)
      retVal.push(aux);
  }
  return retVal;
}

function onError(err) {
  console.log('error: ' + err);
}

function haversine(lat1, lon1, lat2, lon2) {
  let R = 6371e3;
  let omega1 = toRadians(lat1);
  let omega2 = toRadians(lat2);
  let deltaomega = toRadians(lat2 - lat1);
  let deltaalpha = toRadians(lon2 - lon1);

  let a = Math.sin(deltaomega / 2) * Math.sin(deltaomega / 2) +
    Math.cos(omega1) * Math.cos(omega2) * Math.sin(deltaalpha / 2) * Math.sin(deltaalpha / 2);

  let c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));

  let d = R * c;
  return d;
}

function toRadians(Value) {
  /** Converts numeric degrees to radians */
  return Value * Math.PI / 180;
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
