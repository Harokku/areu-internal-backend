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

`DATA_TABLE: {string} <- Content data file root`

## Routes

### Auth

### Docs

Root `/docs`

GET `docs/ <- Get all documents info`

GET `docs/:id <- Get single document info by passed DB id`

GET `docs/recent/:num <- Get most recent {num} documents`

GET `docs/serveById/:id <- Download file by passed DB id`

### Content

Root `/content`

GET `content/ <- Get content index and links name`

GET `content/:link <- Get single content by passed link name`
