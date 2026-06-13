# Posts REST API CLI

## Install database schema

bash posts-api install 

## Start the server

bash posts-api start 

Alias:

bash posts-api serve 

## Stop the server

bash posts-api stop 

## Examples

Install the schema and start the API:

bash posts-api install posts-api start 

The API listens on http://localhost:8200 by default and exposes endpoints for:

- Platforms
- Users
- Communities
- Posts

## Commands

| Command | Description |
|-----------|-------------|
| install | Creates or updates the database schema |
| start | Starts the REST API server |
| serve | Alias for start |
| stop | Stops the running server |