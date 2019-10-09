# Focal point reverse proxy

Focal is a reverse proxy that allows for dynamic loading of route targets by configuration file behind JWT authentication. Currently Focal implements PAM authentication locally against the system on which it is running, but it would be trivial to adjust to support any other form of authentication.

## Config File
Included in the repository is an `.example.directory.yml` file, but the structure of the file is as follows:

```
- name: name_of_route
  upstream: http://site.you.want.to.proxy
```

## Environment Variables
* `TOKEN_SECRET` : The secret used for encrypying / decrypting your tokens
* `LISTEN_PORT`: The port on which the application will load and listen
* `DIRECTORY_FILE` : Specifies the location of the yml file to load at start