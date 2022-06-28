# mypass: a command line tool for storing credentials for website

mypass is a command line tool built using Go which provides credentials management features. It is inspired by https://github.com/ejcx/passgo, but there are a few differences. The main difference is that that mypass uses Argon2id as the hashing algorithm for the master password instead of scrypt. Functionalities wise, mypass also allow users to store both their username and password. Another main difference in the functionality is that users are required to authenticate with their master password even when updating their username and password which isn't the case for passgo. The goal of this project is to provide an easy way to manage password from the terminal instead of using 3rd party applications or storing all your credentials in a notepad. (like me) 

## Installation

Navigate to project directory and run:

```bash
$ go install
```

## Commands

### Initialize

Initialize mypass:

```bash
$ mypass init
```

### Add

Add credentials for finance/ocbc:

```bash
$ mypass add <site name>
```

Example:

```bash
$ mypass add finance/ocbc
```

### Show

Show all site credentials:

```bash
$ mypass show
```

Show specific site credentials:

```bash
$ mypass show <site name>
```

Example:

```bash
$ mypass show finance/ocbc
```

### Edit

Edit credential:

```bash
$ mypass edit
```

Example:

```bash
$ mypass edit finance/ocbc
```




