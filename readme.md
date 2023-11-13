# MiniCas
### This is a very basic mock of the cas login flow implementing login and serviceValidate endpoints intended for dev use. I made it because getting apereo to run sucks and the login flow itself is very simple

requirements:
- run the server in a directory with config.yaml. it should include the following
```yaml
user: some-user-id # this will be returned by serviceValidate as json
port: 8082 # server port :)
```