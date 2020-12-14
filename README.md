Inclusion Scanner
---

### LexisNexis Inclusion & Diversity
Inclusion & Diversity are integral to enhancing our culture by attracting, developing, retaining, promoting and rewarding talent to drive our organization’s mission of advancing the Rule of Law around the world.  We are committed to eliminating systemic racism in our legal systems and to being an advocate for Black lives.

### Taking Action
Inclusion Scanner provides users with the ability to scan for racist and non-inclusive terms inside project files then analyze the results using a dashboard so they can be corrected with alternative terms.

### Terms
The project provides a set of racist and potentially non-inclusive terms. Beyond the default terms, users can add or remove terms in the configuration to be included in the scan.

---

## Quick Start

### Prerequisites
- Clone this project
- Install [Golang](https://golang.org/doc/install)
- Install [Docker](https://docs.docker.com/get-docker/)
- Install [docker-compose](https://docs.docker.com/compose/install/)

### Configure
The "config.yml" file is used to specify the ElasticSearch host URL, terms to scan for, file extensions to include, and directories to exclude.
These steps are optional as a set of sensible defaults are provided. 

1. Edit the elasticsearch section to point to your host. If you don't have an ElasticSearch host then you can use the docker-compose.yaml to set one up in later steps.
```yaml
elasticsearch:
  url: http://localhost:9200
```

2. Change the name of the index that will be used in elasticsearch.
```yaml
scanners:
  repositories:
    # NOTE: a timestamp formatted as '-YYYY-MM-DD' is appended to the index
    database_index: inclusion-scanner-repositories
```

3. Add or remove terms from the scanner by editing the "terms" section in config.yml:
```yaml
terms:
  - "master"
  - "slave"
  - "blacklist"
  - "whitelist"
  ... ( and more ) ...
```
Note: terms are case-insensitive

4. Add or remove file extensions to include in the scan.
```yaml
included_extensions:
  - ".yaml"
  - ".yml"
  - ".json"
  ... ( and more ) ...
```

5. Add or remove excluded directories in the scan.
```yaml
excluded_directories:
  - ".git"
  - ".eggs"
  - "build"
  ... ( and more ) ...
```

### Start the data sink and data visualization
ElasticSearch and Kibana will start up.
```shell
docker-compose up
```

### Build & Run
Using a new terminal, go into the source directory to build the binary then run the command with the path to a directory containing repositories.
The scanning process will take several minutes, depending on the number of files in your repositories.

To scan only your local repository provide the path to the project.
```shell
./inclusion-scanner ~/myproject
```

Or to scan a directory containing several projects specify the parent directory.
```shell
./inclusion-scanner ~/myrepositories
```

Details will be output while the scan is in progress displaying the total number of scanned documents can viewed.
```text
File /Users/me/projects/myproject/README.md contains non-inclusive term(s): master
File /Users/me/projects/myproject/main.go contains non-inclusive term(s): blacklist
... ( and more ) ...
```

### Visualize The Results
If the docker-compose.yaml file was used then Kibana is available at `localhost:5601` by default.
Go to Kibana and view the index named `inclusion-scanner-repositories-YYYY-MM-DD`. This index name can be changed in config.yml.

See `VISUALIZATION.md` for more details.

#### Teardown
```shell
docker-compose down
```

### Note about Filepaths

When running the inclusion scanner, it is best to put the executable in your path and run the executable from the project root directly, however, if this is not possible then use absolute paths instead of relative paths.
---

### MIT License
Copyright 2020 Lexis Nexis
