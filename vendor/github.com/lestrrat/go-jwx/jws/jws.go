// Package jws implements the digital signature on JSON based data
// structures as described in https://tools.ietf.org/html/rfc7515
//
// If you do not care about the details, the only things that you
// would need to use are the following functions:
//
//     jws.Sign(payload, algorithm, key)
//     jws.Verify(encodedjws, algorithm, key)
//
// To sign, simply use `jws.Sign`. `payload` is a []byt