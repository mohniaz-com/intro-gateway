# Intro Gateway

A lightweight Discord application designed to grant a specific role to users who send a message in a designated channel. This project was developed primarily for use in the [852 Resurgence Discord server](https://discord.gg/852r).

## Requirements

- Go 1.23 or above
- A Discord bot token with the appropriate permissions
- Docker (optional for containerized deployment)
- GNU Make (optional for a simplified development CLI)

## Environment Variables

Define the following variables in a .env file in the root directory or in the system environment:

- TOKEN: Discord bot token.
- GUILD_ID: ID of the guild (server).
- CHANNEL_ID: ID of the channel to monitor.
- ROLE_ID: Role ID to assign.
- LOG_ID: Channel ID for logging bot actions.
- JOIN_MSG: Direct message to send to users when they join the guild.

## Compilation

1. Clone the repository:

```
git clone https://github.com/mohniaz-com/intro-gateway.git
cd intro-gateway
```

2. Install dependencies:

```
go mod tidy
```

3. Build the binary:

```
make build
```

## Usage

Run the program locally:

```
make run
```

Or start the service in a Docker container:

```
make image
make start
```

To stop and remove the container:

```
make stop
```

## Endpoints

• Health Check: Accessible at `/health`. Returns a 200 OK response if the service is healthy.

## Logging

All actions are logged to the channel specified in the LOG_ID environment variable.

_Syntax <3_
