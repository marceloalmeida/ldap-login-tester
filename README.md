# LDAP Authentication CLI

This repository contains a CLI tool for authenticating against an LDAP server. The tool is built using Go and leverages the `cobra` library for command-line interface creation.

## Features

- Authenticate against an LDAP server using provided credentials.
- Search and list groups associated with the authenticated user.
- Supports multiple environments through `.env` files.

## Prerequisites

- Go 1.24.1 or later
- An LDAP server for authentication

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/marceloalmeida/ldap-auth-test.git
    cd ldap-auth-test
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

## Configuration

Create a `.env` file in the root directory with the following content:

### Examples
- [FreeIPA Demo Server client configuration](./.env_freeipa)
- [Forum Systems LDAP Test Server client configuration](./.env_forumsys)


## Usage

### Download from releases

#### Execute download binary

```sh
< download location >/ldap-login-tester
```

#### Allow application to run on MacOS

```sh
xattr -r -d com.apple.quarantine < download location >/ldap-login-tester
```

### Building localy

To run the CLI tool, use the following command:

```sh
go run main.go
```

You will be prompted to enter your username and password. The tool will then authenticate against the LDAP server and list the groups associated with the authenticated user.

#### Building the Project

To build the project, run:

```sh
go build -o ldap-auth
```

This will create an executable named `ldap-auth` in the root directory.


## License

This project is licensed under the MPL-2.0 license. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.


## Thanks

Special thanks to:
- [FreeIPA Demo Server](https://www.freeipa.org/page/Demo) for providing a test environment
- [Forum Systems LDAP Test Server](https://www.forumsys.com/2022/05/10/online-ldap-test-server/) for their public LDAP test server
