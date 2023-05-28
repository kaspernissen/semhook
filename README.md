# semhook
Combining [Semgrep](https://github.com/returntocorp/semgrep) and [Starhook](https://github.com/fatih/starhook) to do on demand scanning of multiple repositories

Semhook is a web application.
It serves the following endpoints:

```
/scan
```

```
/sync
```

```
/repo
```

### /scan

POST request.
Accepts a file containing the rule you want to run on all the repositories

### /sync

GET request. 
Ensures all repositories are up to date

## Getting started

[Semgrep](https://github.com/returntocorp/semgrep) and [Starhook](https://github.com/fatih/starhook) must be available on the host.

Starhook with a configuration that syncs the repositories you want to run tests agains.

Set the environment variable `SEMHOOK_REPO_ROOT=<rRepositories Directory>` from the output of `$ starhook config show`.


## TODO

- Consider if options should be given to the scan
- Initializing starhook remotely
- list starhook configs and switch between them
- cache results?


Pieced together with the help of ChatGPT <3