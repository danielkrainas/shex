var express = require('express');
var bodyParser = require('body-parser');
var morgan = require('morgan');
var app = express();
var path = require('path');


app.use(morgan('dev'));
app.use(bodyParser.urlencoded({ extended: false }));
app.use(bodyParser.json());

app.post('/profiles/:user/:profileId', function (req, res) {
	console.log('profile posted:');
	console.log(JSON.stringify(req.body, null, 4));
	res.send('1')
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

app.get('/mods/:user/:mod/:version/meta', function (req, res) {
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

app.get('/mods/:user/:mod/:version', function (req, res) {
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