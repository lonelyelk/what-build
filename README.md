# what-build

A project created as an excercise to learn some Go for the Greater Good. Idea conceved as a help to
coordinate QA process.

## Installation

To be defined.

## Usage

```
  what-build [command]

Available Commands:
  find        find a build of a project
  help        Help about any command

Flags:
  -h, --help   help for what-build

Use "what-build [command] --help" for more information about a command.
```

## Config

When running for the first time, the tool looks for configuration in **~/.what-build.yaml** and uses
it to access AWS Parameter Store and fetch the configuration. It is assumed that the user has AWS
credentials configured. The path to SSM parameter and the region are queried if missing in local config.
The parameter is expected to contain a json string with projects and builds information:

```json
{
    "settings": {
        "per_page": 20,
        "max_offset": 200
    },
    "projects": [
        {
            "name": "proj1",
            "circleci_url": "https://circleci.com/api/v1.1/project/...",
            "circleci_token": "token_for_project"
        },
        {
            "name": "proj2",
            "circleci_url": "https://circleci.com/api/v1.1/project/...",
            "circleci_token": "token_for_project"
        }
    ],
    "builds": [
        {
            "name": "build1",
            "search_build_parameters": {
                "SOME_ENV": "superprod",
                "SOME_VAR": "true"
            }
        },
        {
            "name": "build2",
            "search_build_parameters": {
                "SOME_ENV": "QA999",
                "SOME_OTHER_VAR": "false"
            }
       }
    ]
}
```
