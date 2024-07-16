# Simple Bank Project

This project is a simple banking application built in Go. It demonstrates the use of Go's standard `database/sql` package along with PostgreSQL for database interactions. The project structure includes a testing setup that leverages environment variables for database configuration.

## Getting Started

To get started with the Simple Bank Project, follow these steps:

### Prerequisites

- Go (1.15 or later)
- PostgreSQL
- A `.env` file with your database configuration

### Installation

1. Clone the repository to your local machine:

```bash
git clone https://github.com/yourusername/simplebankproject.git
cd simplebankproject
```

2. Ensure PostgreSQL is running and create a database for the project.

3. Create a `.env` file in the root of the project with the following configuration:

```bash
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
DB_PORT=5432
DB_HOST=localhost
DB_SSLMODE=disable
```

4. Install the project dependencies:

```bash
go mod tidy
```

Running Tests
```bash
go test -v ./...
```
this will run all the tests in the project

## Contributing
Contributions to the Simple Bank Project are welcome. Please ensure to follow the standard Go coding guidelines and add tests for new features.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.
