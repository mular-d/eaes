# EAES

A simple web server to simulate the resource required to run the EAES API and retrieve student results from https://eaes.et/.

## Usage

1. Seed the PostgreSQL database with sample data.
2. Start the web server.
3. Perform a performance test using a tool like `hey`.

### Seed the Database

Before running the web server, make sure to seed the PostgreSQL database with sample data. You can use the provided `seed.go` file to insert random data into the `results` table.

```bash
go run seed.go
```

### Start the Web Server

To start the web server, run the following command:

```bash
go run main.go
```

The server will start listening on `http://localhost:8080`.

### Perform a Performance Test

To perform a performance test, you can use the `hey` CLI tool. Here's an example command to send 100,000 requests with 10,000 concurrent connections:

```bash
hey -n 100000 -c 10000 http://localhost:8080/result?id=1
```

Make sure to replace `http://localhost:8080/result?id=1` with the actual URL of your endpoint.

### Performance Testing Result

An example performance testing result is available in the `scripts/performance_result.txt` file in the repository. You can refer to this file to see the performance metrics and results obtained during the test.

## Note

Please note that this is a simulated server and does not actually connect to the EAES API. It is intended for testing and demonstration purposes only.

Make sure you have PostgreSQL installed and running, and update the database connection details in the `config.go` file.
