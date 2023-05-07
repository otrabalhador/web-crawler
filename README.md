![Build](https://github.com/otrabalhador/web-crawler/actions/workflows/build.yaml/badge.svg)

# Recursirve mirroring web crawler

## What does this do?

It is a command line tool that will accept a url and download its content to a destination directory. For each url inside the content of the root page, the web crawler will extract the content and save it to the destination directory in a hierarchical way. It is recursive and will only end when all the pages are fetched or when the iteration limit is reached.

## How to use it?

TODO

## How this repository is organized?

TODO

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

### Common issues

#### Circular reference

TODO

#### Race condition on persistence

TODO


## Development ROADMAP

- [x] Create a repo with README on github
- [x] Plan initial structure of application
- [x] Add CI for build, unit test, linter
- [ ] Create crawler using interfaces of WebClient, Repository and Extractor
- [ ] Add logging
- [ ] Fix circular dependency problem
- [ ] Create CLI using dummy dependencies of crawler
- [ ] Implement WebClient
- [ ] Implement Repository
- [ ] Implement Extractor
- [ ] Implement way of resuming work
- [ ] Add concurrency
