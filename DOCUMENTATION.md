# DevOps - Group 08 [![sync](https://github.com/ruegerj/devops/actions/workflows/sync-docs.yaml/badge.svg)](https://github.com/ruegerj/devops/actions/workflows/sync-docs.yaml)

Members:

* Andreas Affentranger
* Jan Rüger
* Laurin Scholtysik
* Mirco Stadelmann

VMs:

* devops-bastion (83.228.209.199)
* srv-001.devops.ls.eee.intern
* srv-019.devops.ls.eee.intern
* srv-022.devops.ls.eee.intern
* srv-023.devops.ls.eee.intern

## Table of contents

- [Project Idea & Tools](#project-idead--tools)
- [Prerequisites & Setup](#prerequisites--setup)
  - [Local Workspace](#local-workspace)
- [Applications](#applications)
  - [REST API](#rest-api)
  - [Web Frontend](#web-frontend)
- [Pipelines](#pipelines)
  - [Continous Integration - API](#continous-integration---api)
  - [Continous Integration - Web](#continous-integration---web)
- [Miscellaneous](#miscellaneous)
  - [GitLab Docs Sync](#gitlab-docs-sync)


## Project Idea & Tools

**Application**

Dummy app that displays some information fetched from a private REST api.

- [X] Frontend - SvelteKit
- [X] Backend - Go REST API
- [ ] (Database - PostgreSQL) optional as extension

**DevOps Features**

- [X] CI Pipeline (GitHub Actions)
    - [X] Build
    - [X] Lint
    - [X] Test
    - [X] Static Vulnerability Analysis (Snyk)

- [ ] CD Pipeline (Github Actions)
    - [X] Build
    - [X] Docker Image Build
    - [X] Image Push (GitHub Container Registry)
    - [ ] (ArgoCD Sync)

- [ ] K8S or K3S Hosting
- [ ] ArgoCD for Deployments
- [ ] Credential Vault (Hashicorp Vault?)
- [X] SSH Reverse Tunnel
- [ ] Configuration as Code (TerraForm + Ansible?) -> Desaster Recovery

**Extensions**

- [ ] New feature with toggle (feature flag)
- [ ] Database with automated backup
- [ ] Database with schema change

## Prerequisites & Setup

### Local Workspace

In order work on the applications & pipelines locally, one needs the following tools:

- [Go](https://go.dev/) (v1.25)
- [golangci-lint](https://golangci-lint.run/) (>= v2.5.0)
- [Task](taskfile.dev) (v3)
- [Docker](https://www.docker.com/)
- [pnpm](https://pnpm.io/) (v10.x)

Recommended but not mandatory:

- [act](https://github.com/nektos/act)
- [openssl](https://openssl-library.org/) (if you want to generate JWTs by hand)

For starting the application locally, testing, linting, or any other automations [Task](taskfile.dev) is used as a modern cross-platform Make
alternative. You can run task cmd's either in the repository root or in the dedicated subdirectories of the applications:

```bash
# start the api (cwd: ./)
task api:run

# start the api (cwd: ./api)
task run

# use: task --list-all for discovery of other cmd's
```

> All commands are written with unix systems in mind (Linux & Mac). If you are on Windows - try using a WSL instance if you are facing issues.

### Git workflow

The workflow should be pragmatic and enable productivity, however the following guardrails apply:

- every meaningful change should be done in a dedicated feature branch
- every feature branch should be either _squash-merged_ or _rebased_ back onto `main`
- every meaningful change should be peer-review via pull-request
- CI should be passing before a change gets onto `main`

### Docker Containers

The release pipeline will build, tag and push two docker images to the GitHub Container Registry. One for the Backend and one for the frontend.
During the build process the JWT_KEY and matching ACCESS_TOKEN are injected as build arguments into the respective docker images.

Note: JWT_KEY is defined as a github secret.

**Running the containers**
- Backend
  - Forward port 3000 to the host system.
- Frontend
  - Forward port 4173 to the host system.
  - Set API_BASE_URL to the backend URL (http://host.docker.internal:3000)
  - Note: http://localhost:3000 does not work as with this the container uses its INTERNAL localhost


You can run the containers with the following commands after loggin into ghcr.io with username and token:
```bash
docker pull ghcr.io/ruegerj/devops/api:latest
docker pull ghcr.io/ruegerj/devops/web:latest

docker run -d --name frontend -e API_BASE_URL=http://host.docker.internal:3000 -p 4173:4173 ghcr.io/ruegerj/devops/web:latest
docker run -d --name backend -p 3000:3000 ghcr.io/ruegerj/devops/api:latest
```

## Applications

The following diagram gives a brief overview on the application architecture. For more details see the corresponding sub sections.

```mermaid
---
title: Architecture Overview
---
flowchart LR
    browser["Browser"]
    svelte["SvelteKit"]
    api["Go API"]

    browser --> |Http| svelte
    svelte -- "REST (JWT)" --> api
    api -.-> |Protected content| svelte
    svelte -.-> |Http| browser
```

### REST API

The REST API is written in Go. It features one main endpoint holding a "secret" value which can only be accessed when authenticated.

**Configuration:**

All configuration options are supplied as environment variables. When running the app locally via _Task_, they are automatically
supplied from the `.env` file. One for development purposes can be generated using the following command:

```bash
task generate:env
```

The following variables need to be set for the app to function properly:

| Key     | Description                                                       | Default value (Docker) |
| ------- | ----------------------------------------------------------------- | ---------------------- |
| HOST    | holds the hostname which the server listens on                    | 0.0.0.0                |
| PORT    | holds the TCP port number the server listens on                   | 3000                   |
| JWT_KEY | holds the symmetric signing key used to verify the JWT signatures | _none_                 |

**Endpoints:**

- `GET /health` - Indicates if the application is up and running
  - Authentication: _none_
  - Sample response:
    ```json
    {
      "status": "healthy",
      "timestamp": 1759777840
    }
    ```
- `GET /api/secret` - Holds secret value, only to be accessed by authenticated requests
  - Authentication: _Bearer token_ (JWT, HS256 symmetrically signed)
  - Sample request:
    ```bash
    curl --request GET \
      --url http://localhost:3000/api/secret \
      --header 'authorization: Bearer <jwt>'
    ```
  - Sample response:
    ```json
    {
      "message":"Some text",
      "number": 27
    }
    ```

**Tests:**

Most of the app's code is tested using unit tests. They live right next to their implementation in `x_test.go` files and are usually scoped in
their own testing module.

End-to-end (e2e) tests are located under _tests/e2e_. Before they are executed, the complete app is bootstrapped as [testcontainer](https://testcontainers.com/).
To be able to distinguish these tests from the unit tests, they must have the [build tag](https://pkg.go.dev/go/build#hdr-Build_Constraints) `e2e`.

In order to run the tests locally, use the following commands:

```bash
task test # both unit & e2e
task test:unit # unit tests only
task test:e2e # e2e tests only
```

### Web Frontend

The frontend is built using [SvelteKit](https://svelte.dev/docs/kit/introduction) as meta framework. It renders the UI and makes
authenticated calls to the [REST Api](#rest-api) in order to fetch data.

**Configuration:**

All configuration options are supplied as environment variables. When running the app locally, they are automatically supplied from
the `.env` file. One for development purposes can be generated using the following command:

```bash
task generate:env
```

> Note: the `.env` file of the api must already be present before running the command above

The following variables need to be set for the app to function properly:

| Key          | Description                                                        | Default value (Docker) |
| ------------ | ------------------------------------------------------------------ | ---------------------- |
| NODE_ENV     | holds the environment type the server should run in                | production             |
| ORIGIN       | holds the URL the application should listen on                     | http://localhost:4173  |
| PORT         | holds the TCP port number the server listens on                    | 4173                   |
| API_BASE_URL | holds the base url pointing to the api instance                    | _none_                 |
| ACCESS_TOKEN | holds the JWT access token used for authentication against the api | _none_                 |

**Endpoints:**

- `GET /api/secret` - Proxy for the same endpoint provided by the api (request will be enriched with private access token)
  - Authentication: _none_
  - Sample request:
    ```bash
    curl --request GET \
      --url http://localhost:5173/api/secret
    ```
  - Sample response: see [API > _Endoints_](#rest-api)

**Routes:**

- `/` - Displays the UI of the application
  - Authentication: _none_

**Tests:**

The app's code is tested using unit tests on component level. However due to the very limited functionality which is located in
the frontend their amount is very limited. In the case that project grows the testing infrastructure ([vitest](https://vitest.dev/)
) can be built uppon and extended with additional tests.

End-to-end (e2e) tests can be found in the _e2e_ directory. Here [Playwright](https://playwright.dev/) is used to test the actual
site for multiple headless browsers. In order to test all functionality, a complete api instance gets spun up using a
[testcontainer](https://testcontainers.com/).

In order to run the tests locally, use the following commands:

```bash
task test # both unit & e2e
task test:unit # unit tests only
task test:e2e # e2e tests only
```

> Before you execute e2e tests for the first time, make sure to run this command:
> ```bash
> pnpm dlx playwright install
> ```

## Pipelines

### Continous Integration - API

```mermaid
---
title: CI API
---
flowchart LR
    trigger_push[/push 'main'/]
    trigger_pr[/PR -> 'main'/]
    build(build)
    lint(lint)
    test_unit(unit test)
    test_e2e(e2e test)
    scan(scan)

    trigger_push --> build
    trigger_pr --> build
    build --> lint
    build --> test_unit
    build --> test_e2e
    lint --> scan
    test_unit --> scan
    test_e2e --> scan
```

The pipeline runs for every _push_ and _pull request_ targeting the `main` branch, which is holding changes in the `api`
directory. It features the following steps:

- **build** - Builds the entire codebase in order to find any compilation errors.
- **lint** - Statically checks the codebase for potential quality flaws using [golangci-lint](https://golangci-lint.run/)
- **unit test** - runs all unit tests
- **e2e test** - runs all e2e tests against a testcontainer instance of the app (see _Tests_ section of [API](#rest-api))
- **scan** - runs static vulnerability scans using [Snyk](https://snyk.io/)
  -> results are uploaded as [GitHub code scanning alerts](https://docs.github.com/en/code-security/code-scanning/managing-code-scanning-alerts/about-code-scanning-alerts)

### Continous Integration - Web

```mermaid
---
title: CI Web
---
flowchart LR
    trigger_push[/push 'main'/]
    trigger_pr[/PR -> 'main'/]
    build(build)
    lint(lint)
    test_unit(unit test)
    test_e2e(e2e test)
    scan(scan)

    trigger_push --> build
    trigger_pr --> build
    build --> lint
    build --> test_unit
    build --> test_e2e
    lint --> scan
    test_unit --> scan
    test_e2e --> scan
```

The pipeline runs for every _push_ and _pull request_ targeting the `main` branch, which is holding changes in the `web`
directory. It features the following steps:

- **build** - Builds the entire codebase in order to find any compilation errors.
- **lint** - Statically checks the codebase for potential quality flaws using [eslint](https://golangci-lint.run/) & SvelteKit
linter. Additionaly it checks the format of every file using [prettier](https://prettier.io/).
- **unit test** - runs all unit tests using [vitest](https://vitest.dev)
- **e2e test** - runs all e2e tests using [Playwright](https://playwright.dev/) (see _Tests_ section of [Web](#web-frontent))
- **scan** - runs static vulnerability scans using [Snyk](https://snyk.io/)
  -> results are uploaded as [GitHub code scanning alerts](https://docs.github.com/en/code-security/code-scanning/managing-code-scanning-alerts/about-code-scanning-alerts)

## Infrastructure

### Network Access for Github Actions
On order to allow access to the Enterprise Lab VMs from the GitHub Action pipelines without the need of Pulse Secure VPN, a SSH reverse tunnel from `srv-001-devops.ls.eee.intern` to a bastion host (`devops-bastion` - 83.228.209.199) located at a public cloud provider was estabilished.

```mermaid
---
title: SSH Reverse-Tunnel
---
sequenceDiagram
    Github Actions Runner->devops-bastion: ssh -i id_devops -p 3333 labadmin@devops-bastion
    srv-001-devops.ls.eee.intern->devops-bastion: ssh -R 3333:localhost:22 debian@devops-bastion -i id_devops-bastion
```
## Miscellaneous

### GitLab Docs Sync

The main entry point for the course & grading is a repository hosted on the [Switch GitLab](https://gitlab.switch.ch/hslu/edu/bachelor-computer-science/devops/25hs01/g08/g08-documentation/).
As a fancy way to be able to view the documentation still in the GitLab repository, a GitHub Actions [workflow](/.github/workflows/sync-docs.yaml)
was created in order to syn  the contents of this file to the GitLab repo.

**Limitations:**

 - only the single file `DOCUMENTATION.md` is supported -> no support for standalone artifacts (e.g. images)
 - it uses a GitLab repository access token with a hard expiration date -> has to be checked and renewed manually (see [here](https://gitlab.switch.ch/hslu/edu/bachelor-computer-science/devops/25hs01/g08/g08-documentation/-/settings/access_tokens))
 - this workflow is _not_ concurrency safe (when multiple actions run concurrently the last to run wins) -> could lead to an incorrect state of the
 GitLab docs

**Considerations:**

The GitLab access token must have the scope `api` with the role _Maintainer_ in order to be able to update files on the `main` branch.
Thus the token grants one full permissions on most of the sensitive actions of the repository, which could lead to unpleasant outcomes when its
leaked. However since this token is on the repository-level and the GitLab documentation repository is not the source of truth (e.g. doesn't hold
any information which may be lost) it is considered an acceptable risk.

**Implementation:**

At it's core it simply uses the [Update File endpoint](https://docs.gitlab.com/api/repository_files/#update-existing-file-in-repository) of the
GitLab REST API.

- it runs on any commit on the `main` branch which alters the file `DOCUMENTATION.md`
- [jq](https://jqlang.org/) is used to construct the JSON body for the file update request
- the GitLab _project-id_ & _access-token_ are stored as Action Variable or Secret respectively
- [cURL](https://curl.se/) is used to perform the HTTP request to the GitLab API
