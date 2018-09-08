# GOAuth
OAuth2 server

## Goal

We regularly make small applications. Every new application has to have an user
database, a register/login flow etc. We decided to setup a SSO (Single Sign On),
storing our users for every application in a single database and having a
pluggable HTML component we can include in every app to handle sign ups and sign
ins. 

This was also an occasion to discover AWS Lambdas, AWS API Gateway, OAuth and
the GitHub API.

## Flow

1. Display a HTML file with a `Login with GitHub` button
2. Call `https://github.com/login/oauth/authorize` on click
3. GitHub calls back on our AWS API Gateway, acting as a proxy to a lambda
4. The `handlecallback` lambda gets triggered with a `code`
5. Use this `code` to get an access token at `https://github.com/login/oauth/access_token`
6. Use this access token to get the authenticated user at `https://api.github.com/user`
7. Store the user info in our database

## Resources used

[Spec](https://tools.ietf.org/html/rfc6749)  
[Building a Basic Auth Server](https://medium.com/google-cloud/understanding-oauth2-and-building-a-basic-authorization-server-of-your-own-a-beginners-guide-cf7451a16f66)  
[Go Lambdas](https://github.com/eawsy/aws-lambda-go)

[GitHub API doc](https://developer.github.com/v3/)  
[GitHub OAuth flow](https://developer.github.com/apps/building-oauth-apps/authorizing-oauth-apps/)

GitHub user API : https://api.github.com/users/:username or 
https://api.github.com/user with the token in scope

## Installation

**Requirements**

You need to install Go to run this repo.  
You also need a PostgreSQL, local or remote (we used https://elephantsql.com)

**Environment**

Setup the following environment variables:

- `DATABASE_HOST`: database URL
- `DATABASE_PORT`: database port
- `DATABASE_USERNAME`: database login
- `DATABASE_PASSWORD`: database password
- `DATABASE_DATABASE`: database name
- `GH_ID`: application ID (found at https://github.com/settings/developers)
- `GH_SECRET`: application secret (same)

**Database Migrations**

Once everything is set up, you can create the database structure by running
the migrations:

```
go run database/migrations/migrate.go database/migrations/queries
```

You can provide folders or files to the executable, if you want to run only the
first two migrations, for example, you can do:

```
go run database/migrations/migrate.go database/migrations/queries/0CreateUserTable.sql database/migrations/queries/1AddBasicColumns.sql
```

or for a shorter syntax:

```
cd database/migrations/queries
go run ../migrate.go 0CreateUserTable.sql 1AddBasicColumns.sql
```

**Build**

To build all the lambdas, execute `./build.sh`.  
To build specific lambdas, add arguments: `./build.sh lambda1 lambda2`.

Linux/amd64 executables will be built in `./bin`, and the ready-for-deployment
`.zip` file will be put in `./dist`.


>NOTE: The lambdas are meant to be served with a proxy API Gateway method.
