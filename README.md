# q3errord
REST endpoint for recording execptional events

###  error using slack endpoint: https://hooks.slack.com/services/xxx: x509: failed to load system roots and no roots provided
Fix: put root certs in the image.  From an ubuntu host, you can do this thusly:
    docker run -d  -v '/etc/ssl/certs/ca-certificates.crt:/etc/ssl/certs/ca-certificates.crt'  -p 8081:9999 tompscanlan/q3errord /q3errord --port 9999 --slack-webhook https://hooks.slack.com/services/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
