# wrike-me-up
Never gonna give you up, never gonna let you down... but it will log your hours into Wrike!

# Get authorization_code

https://login.wrike.com/oauth2/authorize/v4?client_id={client_id}&response_type=code

- Grab the code from the URL

# Get the access_token & refresh token

curl -X POST -d "client_id={client_id}&client_secret={client_secret}&grant_type=authorization_code&code={authorization_code}" https://login.wrike.com/oauth2/token