# DevOps - Group 08 [![sync](https://github.com/ruegerj/devops/actions/workflows/sync-docs.yaml/badge.svg)](https://github.com/ruegerj/devops/actions/workflows/sync-docs.yaml)

Members:

* Andreas Affentranger
* Jan Rüger
* Laurin Scholtysik
* Mirco Stadelmann

VMs:

* srv-001.devops.ls.eee.intern
* srv-019.devops.ls.eee.intern
* srv-022.devops.ls.eee.intern
* srv-023.devops.ls.eee.intern

## Table of contents

- [Project Idea & Tools](#project-idead--tools)
- [Miscellaneous](#miscellaneous)
  - [GitLab Docs Sync](#gitlab-docs-sync)


## Project Idea & Tools

**Application**

Dummy app that displays some information fetched from a private REST api.

- [ ] Frontend - SvelteKit
- [ ] Backend - Go REST API
- [ ] (Database - PostgreSQL) optional as extension

**DevOps Features**

- [ ] CI Pipeline (GitHub Actions)
    - [ ] Build
    - [ ] Lint
    - [ ] Test
    - [ ] Static Code Analysis (Snyk?)

- [ ] CD Pipeline (Github Actions)
    - [ ] Build
    - [ ] Docker Image Build
    - [ ] Image Push (GitHub Container Registry)
    - [ ] (ArgoCD Sync)

- [ ] K8S or K3S Hosting
- [ ] ArgoCD for Deployments
- [ ] Credential Vault (Hashicorp Vault?)
- [ ] SSH Reverse Tunnel
- [ ] Configuration as Code (TerraForm + Ansible?) -> Desaster Recovery

**Extensions**

- [ ] New feature with toggle (feature flag)
- [ ] Database with automated backup
- [ ] Database with schema change

## Miscellaneous

### GitLab Docs Sync

The main entry point for the course & grading is a repository hosted by the [Switch GitLab](https://gitlab.switch.ch/hslu/edu/bachelor-computer-science/devops/25hs01/g08/g08-documentation/).
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
