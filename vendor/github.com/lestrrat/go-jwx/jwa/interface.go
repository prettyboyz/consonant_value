package jwa

// KeyType represents the key type ("kty") that are supported
type KeyType string

// Supported KeyTypes
const (
	EC       KeyType = "EC"  // Elliptic Curve
	RSA      KeyType = "RSA" // RSA
	OctetSeq KeyType = "oct" // Octet sequence (used to represent symmetric keys)
)

// EllipticCurveAlgorithm represents the algorithms used for EC keys
type EllipticCurveAlgorithm string

// Supported EllipticCurveAlgorithms
const (
	P256 EllipticCurveAlgorithm = "P-256"
	P384 EllipticCurveAlgorithm = "P-384"
	P521 EllipticCurveAlgorithm = "P-521"
)

// SignatureAlgorithm represents the various signature algorithms
// as described in https://tools.ietf.org/html/rfc7518#section-3.1
type SignatureAlgorithm string

// Supported SignatureAlgorithms
const (
	NoSignature SignatureAlgorithm = "none"
	HS256       SignatureAlgorithm = "HS256" // HMAC using SHA-256
	HS384       SignatureAlgorithm = "HS384" // HMAC using SHA-384
	HS512       SignatureAlgorithm = "HS512" // HMAC using SHA-512
	RS256       SignatureAlgorithm = "RS256" // RSASSA-PKCS-v1.5 using SHA-256
	RS384       SignatureAlgorithm = "RS384" // RSASSA-PKCS-v1.5 using SHA-384
	RS512       SignatureAlgorithm = "RS512" // RSASSA-PKCS-v1.5 using SHA-512
	ES256       SignatureAlgorithm = "ES256" // ECDSA using P-256 and SHA-256
	ES384       SignatureAlgorithm = "ES384" // ECDSA using P-384 and SHA-384
	ES512       SignatureAlgorithm = "ES512" // ECDSA using P-521 and SHA-512
	PS256       SignatureAlgorithm = "PS256" // RSASSA-PSS using SHA256 and MGF1-SHA256
	