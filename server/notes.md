# TODO

# Known Tech Debt
- User info for websocket messaging is retrieved at the moment the WS connection is established. This means future changes to user info won't reflect on the messages.
- No data layer on the API. There is code rewrite for fetching data from de database.
- No database interface.
- No migration system for database
