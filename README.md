# Time Clock API

This repo holds simple and robust micro services that handle the basic functions of a time clock system.

## Services

There are a handful of services that are provided.

### Admin

This service handles the creation of, and changes to users
and employees. Also, approving/rejecting requests made by other
users will be handled here.

Every endpoint requires a bearer token with the `user:admin` scope,
validated against the auth service's gRPC introspection service via the
`pkg/authclient` middleware. The auth service's `AUTH_GRPC_ADDR`
(e.g. `localhost:9090`) must be set for the admin service to reach it.

### Audit

This service handles fetching and recording any changes made to,
or the creation of pertinant data in the database.

### Auth

This service handles user login and the user token to be utilized 
by the other services for authentication.

Tokens issued by this service are opaque: they're cryptographically random
values that carry no embedded claims (unlike a JWT). Other services can't
validate a token on their own; they must ask this service to validate it
for them. Token introspection is exposed two ways, both modeled after
RFC 7662 (OAuth 2.0 Token Introspection):

- `POST /introspect` (HTTP/JSON) &mdash; convenient for ad-hoc debugging,
  non-Go consumers, or anything outside the cluster. Requires `PORT`.
- `TokenIntrospection/Introspect` (gRPC, see `proto/auth/v1`) &mdash; the
  preferred transport for internal, latency-sensitive service-to-service
  calls like the Admin service's per-request auth check. Requires
  `GRPC_PORT`.

Both report whether a token is active along with its scope, kind, and the
user it belongs to, and share the same underlying logic. The
`pkg/authclient` package provides ready-to-use clients for both
(`Client` and `GRPCClient`) plus HTTP middleware (`RequireToken`/
`RequireScope`) built on top of either.

A token's scope is resolved at login time: every authenticated user
implicitly carries the `default` scope, plus whatever has been explicitly
granted to them (e.g. `user:admin`). Scopes are granted/revoked via
`POST /users/{id}/scopes` and `DELETE /users/{id}/scopes/{scope}`, which
are themselves protected by the `user:admin` scope. The very first admin
must be granted directly via a database migration (see
`migrations/auth/002_User-Scopes-Table.sql`), since there's no admin yet
to grant it through the API.

### Base

This service holds the core buisiness logic of the time clock API.
This includes making punches, requesting time off, requesting
to update existing punches or times off entries, and updating a
users own information (e.g. family name, pronouns, gender, etc...)

### Reports

This service will handle the generating of reports need for payrole,
or for any other needs.
