[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/new/template?template=https%3A%2F%2Fgithub.com%2Fmorgangallant%2Fhealthchecker&envs=ENDPOINT%2CDISCORD_URL%2CSECONDS&optionalEnvs=SECONDS&ENDPOINTDesc=HTTP+Endpoint&DISCORD_URLDesc=Discord+Webhook+URL&SECONDSDesc=Number+of+Seconds+between+Checks&SECONDSDefault=30)

Healthchecker is a simple application to periodically send HTTP requests to an endpoint.
If the endpoint returns an error, the application is marked as unhealthy and a notification
is sent to Discord. Healthchecker will continue to issue HTTP requests, and will also send
a Discord notification once the application is back online. There are two required environment
variables (ENDPOINT, DISCORD\_URL), and one optional (SECONDS).
