'use strict';

const path = require('path');

exports.handler = (e, c, cb) => {
  const req = e.Records[0].cf.request;
  if (req.uri.length > 0 && req.uri.charAt(req.uri.length - 1) === '/') {
    req.uri += "index.html";
  } else if (!path.extname(req.uri)) {
    req.uri += "/index.html";
  }
  cb(null, req);
  return true;
};
