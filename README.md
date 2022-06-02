# mypass: a command line tool for storing credentials for website

mypass is a command line tool built using Go which provides credentials management features. The goal of this project is to provide an easy way to manage password from the terminal instead of using 3rd party applications or storing all your credentials in a notepad. (like me)

## Installation

Navigate to project directory and run:

```bash
$ go install
```

## Commands

### add

Add new credential:

```bash
$ mypass add -u <username> -p <password> 
```

Example:

```bash
$ mypass add Github -u Username123 -p Password123
```

### list

List all credentials:

```bash
$ mypass list all
```

List website credential:

```bash
$ mypass list <Website or Application>
```

Example:

```bash
$ mypass list Github
```

### delete

Delete credential:

```bash
$ mypass delete Github
```




