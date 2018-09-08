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
