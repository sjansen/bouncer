'use strict';

import {
    createRemoteJWKSet
} from 'jose/jwks/remote'
import {
    jwtVerify
} from 'jose/jwt/verify'
import {
    JWTExpired
} from 'jose/util/errors'
import {
    parse
} from "cookie";

import {
    JWKS_ENDPOINT,
    PUBLIC_PREFIXES,
    ROOT_IS_PUBLIC
} from "./config"

const JWKS = createRemoteJWKSet(JWKS_ENDPOINT);
const rootIsPublic = ROOT_IS_PUBLIC || PUBLIC_PREFIXES.has("/")
PUBLIC_PREFIXES.delete("/")

exports.handler = async (e, c, cb) => {
    const request = e.Records[0].cf.request;
    const host = request.headers.host[0].value;
    const uri = request.uri;

    if (uri.startsWith('/b/')) {
        console.log(host, uri);
        cb(null, request);
        return;
    }

    if (rootIsPublic && (uri === "/" || uri == "/favicon.ico")) {
        console.log(host, uri, "Root is public.");
        cb(null, request);
        return;
    }

    for (let prefix of PUBLIC_PREFIXES) {
        if (uri.startsWith(prefix)) {
            console.log(host, uri, "Public prefix matched.");
            cb(null, request);
            return;
        }
    }

    if (!request.headers.cookie) {
        console.log(host, uri, "Authentication required.")
        refreshToken(uri, cb);
        return;
    }

    const cookies = request.headers.cookie.reduce(
        (reduced, header) => Object.assign(reduced, parse(header.value)), {}
    );
    if (!cookies.auth_token) {
        console.log(host, uri, "Missing auth_token.")
        refreshToken(uri, cb);
        return;
    }

    try {
        const {
            payload,
            protectedHeader
        } = await jwtVerify(
            cookies.auth_token, JWKS, {
                issuer: 'bouncer',
                audience: 'bouncer-authz',
                subject: host,
            }
        )
        //console.log(protectedHeader);
        //console.log(payload);
        console.log(host, uri, "Access granted.")
        cb(null, request);
        return;
    } catch (err) {
        if (err instanceof JWTExpired) {
            console.log(host, uri, "Refreshing auth_token...")
            refreshToken(uri, cb);
            return;
        }
    }

    console.log(host, uri, "Access denied.")
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
                value: 'text/plain; charset=utf-8'
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
