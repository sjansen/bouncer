'use strict';

import { createRemoteJWKSet } from 'jose/jwks/remote'
import { jwtVerify } from 'jose/jwt/verify'
import { JWTExpired } from 'jose/util/errors'
import { parse } from "cookie";

const JWKS = createRemoteJWKSet(new URL('https://bouncer.tuntap.net/b/jwks/'))

exports.handler = async (e, c, cb) => {
  const request = e.Records[0].cf.request;
  const host = request.headers.host[0].value;
  const uri = request.uri;

  console.log(host, uri);
  if (uri === '/' || uri === '/favicon.ico' || uri.startsWith('/b/')) {
    cb(null, request);
    return;
  }

  if (!request.headers.cookie) {
    console.log("Authentication required...")
    refreshToken(uri, cb);
    return;
  }

  const cookies = request.headers.cookie.reduce(
    (reduced, header) => Object.assign(reduced, parse(header.value)),
    {}
  );
  if (!cookies.auth_token) {
    console.log("Missing auth_token...")
    refreshToken(uri, cb);
    return;
  }

  try {
    const { payload, protectedHeader } = await jwtVerify(
      cookies.auth_token, JWKS, {
        issuer: 'bouncer',
        audience: 'bouncer-authz',
        subject: host,
      }
    )
    //console.log(protectedHeader);
    //console.log(payload);
    cb(null, request);
  } catch (err) {
    if (err instanceof JWTExpired) {
      console.log("Refreshing auth_token...")
      refreshToken(uri, cb);
      return;
    }
  }

  accessDenied(cb);
};

function accessDenied(cb) {
  const response = {
    status: '403',
    statusDescription: 'Forbidden',
    headers: {
      'cache-control': [{
          key: 'Cache-Control',
          value: 'no-cache'
      }],
      'content-type': [{
          key: 'Content-Type',
          value: 'text/plain'
      }]
    },
    body: '403 Access Denied',
  };
  cb(null, response);
}

function refreshToken(uri, cb) {
  cb(null, {
    status: '302',
    statusDescription: 'Found',
    headers: {
        location: [{
            key: 'Location',
            value: `/b/redirect/?target=${encodeURIComponent(uri)}`,
        }],
    },
  });
}
