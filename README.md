
# LinkedIn clone as a webservice

### Techs
- GoLang as the language
- GraphQL as API service
- Fuseki Jena (Apache) as SPARQL server and database
- QLGen to generate resolvers and handle graphql requests
- gorilla, sessions and gin to manage the webserver and sessions
- auth via OAuth with google as provider
- knakk as sparql driver

make sure you have a Fuseki server running:
```
fuseki-server --port=3030 --mem /link_krec
```

