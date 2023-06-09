![Build](https://github.com/otrabalhador/web-crawler/actions/workflows/build.yaml/badge.svg)

# Recursive mirroring web crawler

## What does this do?

It is a command line tool that will accept a url and download its content to a destination directory. For each url inside the content of the root page, the web crawler will extract the content and save it to the destination directory in a hierarchical way. It is recursive and will only end when all the pages are fetched or when the iteration limit is reached.

## How to use it?

- `go build ./cmd/web-crawler && ./web-crawler --root=<ROOT_URL> --destination=<DESTINATION_FOLDER>`

- Example: `go build ./cmd/web-crawler && ./web-crawler --root=https://martinfowler.com/ --destination=`

## How this repository is organized?

- There is a cmd and internal folder
- The CLI is located at the cmd and the crawler components are located at internal

## Application

### Simplified flow

- The crawler (responsible for the flow) will receive the root url
- Calls url
- With the page of response
  - Save it to the destination directory
  - Extract the containing urls
- For each url, it will call the orchestrator again

### Main components

#### Crawler

- Execute flow

#### WebClient

- Call the url and return content page

#### Repository

- Persist the pages
- Return a list of already saved pages 

#### Extractor

- Extract urls given a page page


## Development ROADMAP

- [x] Create a repo with README on github
- [x] Plan initial structure of application
- [x] Add CI for build, unit test, linter
- [x] Create crawler using interfaces of WebClient, Repository and Extractor
- [x] Add recursion
- [x] Fix circular dependency problem
- [x] Implement way of resuming work
- [X] Persist using hierarchy
- [x] Create CLI using dummy dependencies of crawler
- [x] Add logging
- [x] Implement WebClient
- [x] Implement Repository
- [x] Implement Extractor
- [ ] Deal with errors
- [ ] Add concurrency
