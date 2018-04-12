# Crypto Bot


Make sure mongodb is installed & running:

* Install
`apt-get install mongodb`
`mkdir /data`
`mkdir /data/db`

* Start
`sudo service mongodb start`

* Check if is really running
```
sudo service mongodb status
● mongodb.service - An object/document-oriented database
   Loaded: loaded (/lib/systemd/system/mongodb.service; enabled; vendor preset: enabled)
   Active: active (running) since Mon 2017-09-18 15:37:02 CEST; 7h ago
     Docs: man:mongod(1)
 Main PID: 15149 (mongod)
    Tasks: 20 (limit: 4915)
   Memory: 36.8M
      CPU: 2min 7.772s
   CGroup: /system.slice/mongodb.service
           └─15149 /usr/bin/mongod --unixSocketPrefix=/run/mongodb --config /etc/mongodb.conf

sep 18 15:37:02 systemd[1]: Started An object/document-oriented database.
```

## Format code
`make gofmt`

## Test
`make test`

## Start
`make start`


# API KEYS
* To create/manage api keys, use:  `gocryptobot@gmail.com`