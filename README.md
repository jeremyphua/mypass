# password-cli: a command line tool for storing credentials for website

password-cli is a command line tool built using Go which provides credentials management features. The goal of this project is to provide an easy way to manage password from the terminal instead of using 3rd party applications or storing all your credentials in a notepad. (like me)

## Installation

Navigate to project directory and run:

```bash
$ go install
```

## Commands

### add

Add new credential:

```bash
$ password-cli add -u <username> -p <password> 
```

Example:

```bash
$ password-cli add Github -u Username123 -p Password123
```

### list

List all credentials:

```bash
$ password-cli list all
```

List website credential:

```bash
$ password-cli list <Website or Application>
```

Example:

```bash
$ password-cli list Github
```

### delete

Delete credential:

```bash
$ password-cli delete Github
```




