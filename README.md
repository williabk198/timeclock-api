# Time Clock API

This repo holds simple and robust micro services that handle the basic functions of a time clock system.

## Services

There are a handful of services that are provided.

### Admin

This service handles the creation of, and changes to users
and employees. Also, approving/rejecting requests made by other
users will be handled here.

### Audit

This service handles fetching and recording any changes made to,
or the creation of pertinant data in the database.

### Auth

This service handles user login and the user token to be utilized 
by the other services for authentication

### Base

This service holds the core buisiness logic of the time clock API.
This includes making punches, requesting time off, requesting
to update existing punches or times off entries, and updating a
users own information (e.g. family name, pronouns, gender, etc...)

### Reports

This service will handle the generating of reports need for payrole,
or for any other needs.
