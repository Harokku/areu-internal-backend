# areu-internal-backend

SOREU Laghi Internal infrastructure backend

...wip

---

## Configs

### Environment file

You must define Env variables, actually read from local env

#### Parameters

`PORT: {number} <- server port number`
`SECRET: {string} <- JWT sign secret`
`JWT_EXPIRE: {string} <- JWT expire as number + time mod, ex: 24h`
`DATABASE_URL: {string} <- DB connection url`
`DOC_ROOT: {string} <- Document share root`

