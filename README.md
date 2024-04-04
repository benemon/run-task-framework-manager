# Run Task Framework Manager

This project is a command-line tool written in Go that generates scaffolding to assist with the implementation of Terraform Run Tasks in a number of languages. The goal is to remove the initial hurdles associated with creating a custom Run Task by taking care of the basic integration with [Terraform Cloud](https://developer.hashicorp.com/terraform/cloud-docs/workspaces/settings/run-tasks) and/or [Terraform Enterprise](https://developer.hashicorp.com/terraform/enterprise/workspaces/settings/run-tasks).

Whilst there are a number of [Certified Run Tasks](https://registry.terraform.io/browse/run-tasks) within the Terraform Ecosystem, there is often signifcant value derived from being able to integrate custom tools and process specific to an individual use case at key stages in the Terrform workflow where such publically available Run Tasks do not suffice.
 
## Usage

```bash
$ rtfm -name <run task name> -dir <working directory in which the scaffolding will be generated> -language <the runtime language>
```

Depending on the language you chose, this will generate a number of files in the target directory. At a bare minimum, this will include:

* A basic Run Task scaffold in that language, that when used will simply return a valid response to the enterprise Terraform platform.
* A Containerfile that can be built and run with the OCI-compliant container runtime of your choice to ensure consistency and ease of build / deployment. This is currently based on Red Hat's Universal Base Image, and is tested with Docker and Podman. You can of course modify this to suit your requirements if you choose.

## Supported Languages

* Go

## Proposed Languages

* Python
* Node.js
