'use strict';

const express = require('express');
const morgan = require('morgan');
const path = require('path');
const app = express();
const opn = require('opn');

const Web3 = require('web3');
const EthereumTx = require('ethereumjs-tx');
app.use(morgan('dev'));
app.use(express.static(path.resolve(__dirname, '..', 'public')));

const port = process.env.PORT || 8080;
app.listen(port, function () {
	console.log(`Server is listening port ${port}`);
	opn('http://localhost:8080/index.html');
});
