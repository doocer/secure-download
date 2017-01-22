# Secure Download

A little go server with secure download support.


## Run the binary

Run in your teminal:

    $ SECRET_KEY=your-secret-key secure-download [folder]

+ **SECRET_KEY**: A secret token for generate signature
+ **PORT**: HTTP server port
+ **SIGNATURE_IP**: Add client IP into signature

A full example:

    $ SECRET_KEY=your-secret-key SIGNATURE_IP=yes PORT=8700 secure-download [folder]


## Client Signature

A request URL looks like:

```
http://example.com/path/download.bin?e=1485015757&s=be3347c7
```

Parameter `e` is expire time, `s` is signature.

```
signature = md5(${SECRET_KEY} + ${URLPATH} + $expires [ + $remote_addr ])[8:16]
```

Example when not including client IP:

1. secret key: secret
2. expires time: 1485015757
3. request path: /book/hello.epub

The base string is `secret/book/hello.epub1485015757`, its md5 hexdigest is
`ba68c8d4c8d6dde327925df16cea8b17`. Slice hexdigest with `[8:16]`, therefore,
the request URL is:

```
http://example.com/book/hello.epub?e=1485015757&s=c8d6dde3
```
