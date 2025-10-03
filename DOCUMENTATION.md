# DevOps - Group 08

Hello World GitHub Actions!

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

- [ ] K8S oder K3S Hosting
- [ ] ArgoCD for Deployments
- [ ] Credential Vault (Hashicorp Vault?)
- [ ] SSH Reverse Tunnel
- [ ] Configuration as Code (TerraForm + Ansible?) -> Desaster Recovery

**Extensions**

- [ ] New feature with toggle (feature flag)
- [ ] Database with automated backup
- [ ] Database with schema change
