# Contributing

This document describes our mode or collaboration and the [technologies](#technologies) needed to write code.

## Pull Requests

We use git for version control. If you make any changes be sure to open up a new branch beforehand like so: `git checkout -b <00-feature-branch>` Naming convention for feature branches is TicketNumber-short-description.
After commiting your changes push them to the github repository and open up a [pull request](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request). Make sure to assign reviewers, version number, tags. Newcomers to Git will find a tutorial explaining basic commands and local workflow [here](https://www.shellhacks.com/git-basic-workflow/). Before you start working on the project, make sure to [clone](https://docs.github.com/en/repositories/creating-and-managing-repositories/cloning-a-repository) the remote repository.

## Versioning

We use [SemVer](https://semver.org/) for assigning Numbers to Versions.

## Technologies

In order to work at this project it is necessary to have following depedencies installes locally

__Backend__:

* Golang >= 1.17
* docker + docker-compose

__Frontend__:

* Node Package Manager `npm` (install via your operating system's package manager or via [nvm](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm))
* Node.js Runtime >= v 12 LTS ((install via your operating system's package manager or via [nvm](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm)))
* Angular CLI: `npm install @angular/cli`