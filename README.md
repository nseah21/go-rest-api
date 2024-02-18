# go-rest-api

### Instructions for running the server
- Ensure that you have [Golang](https://go.dev/doc/install) installed: `go version`
- Ensure that you have [psql](https://www.postgresql.org/download/) installed: `<Windows key> + <type 'psql'>`
- Pull a copy of the source code using:
  - `git clone git@github.com:nseah21/go-rest-api.git` via `ssh`, 
  - or `git clone https://github.com/nseah21/go-rest-api.git` via `https`,
  - or manually by downloding the ZIP file
- In a terminal window, make sure you are inside the root directory `go-rest-api/`
- To test the server:
  1. Modify the `.env` file with your credentials for `psql`
  2. Start up your `psql` shell using `<Windows key> + <type 'psql'> + <Enter key>` and key in your credentials as in the previous step
  3. Create a database by entering `CREATE DATABASE <database-name>;` 
  4. Change into your database by entering `\c <database-name>`
  5. Copy and paste the definitions block from `./internal/db/init.sql` into your `psql` shell
  6. Optionally, you can also copy and paste the seed data block from `./internal/db/init.sql` into your `psql` shell
  7. In a terminal, verify that you are still in the root directory `go-rest-api/`
  8. Run `go mod tidy` followed by `go run ./cmd/server`
  9. Verify that the server is running on `localhost:8080`

- To run unit tests, run `go test -v ./...` in the root directory `go-rest-api/`
- Alternatively, you can run `go test ./...` for an overview of the test packages