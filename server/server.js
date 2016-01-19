var express = require('express');
var bodyParser = require('body-parser');
var morgan = require('morgan');
var app = express();
var path = require('path');


app.use(morgan('dev'));
app.use(bodyParser.urlencoded({ extended: false }));
app.use(bodyParser.json());

/* Update endpoint */
app.get('/goble/latest', function (req, res) {
	// todo: finish writing this
});

app.get('/goble/latest/check', function (req, res) {
	var version = req.query.version;

	// no new version
	res.status(304).end();

	// new version
	res.send({
		md5: '',// TODO: implement checksum support on client
		version: '0.0.1-development',
		location: '',
		// TODO: include HMAC of data(-HMAC) with data
	});
});

/* API endpoint */
app.post('/profiles/:user/:profileId', function (req, res) {
	console.log('profile posted:');
	console.log(JSON.stringify(req.body, null, 4));
	res.send('1');// return current version number of profile
});

app.head('/profiles/:user/:profileId', function (req, res) {
	//If-Modified-Since: version
});

app.get('/profiles/:user/:profileId', function (req, res) {
	var name = req.params.user + '/' + req.params.profileId;
	console.log('profile request:', name);
	res.json({
		rev: 1,
		name: name,
		mods: {},
	});
});

app.get('/mods/:user/:mod/v', function (req, res) {
	res.json(['0.1.8', '1.0.0']);
});

app.get('/mods/:user/:mod/v/:version/meta', function (req, res) {
	if (req.params.version === 'latest') {
		req.params.version = '1.0.0';
	}

	console.info('params ', req.params);
	var result = {
		name: req.params.user + '/' + req.params.mod,
		version: req.params.version
	};

	console.info('meta ', result);
	res.json(result);
});

app.get('/mods/:user/v/:mod/:version', function (req, res) {
	if (req.params.version === 'latest') {
		req.params.version = '1.0.0';
	}
	
	console.info('params ', req.params);

	res.sendFile(path.join(__dirname, '/candyland2.smod'), function (err) {
		if (err) {
			console.error(err);
			return;
		}
	});
});

app.listen(6231);
console.log('server listening on 6231');