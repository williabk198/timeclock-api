// package auth contains the micro service that authenticates users and issues, validates,
// and revokes the opaque tokens used by the other services to authenticate requests.
//
// Tokens issued by this service are opaque: they are cryptographically random strings
// that carry no embedded claims (unlike a JWT). Other services can't validate a token on
// their own; they must call this service's token introspection endpoint, which reports
// whether a token is active along with its associated scope, kind, and subject.
package auth
