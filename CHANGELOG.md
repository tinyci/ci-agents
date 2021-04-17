## 0.3.0 -- Fri 16 Apr 2021 03:04:27 AM PDT

- hooksvc now listens on 0.0.0.0, not localhost
- independent branches & "do not merge": you can now select specific branches
  to not be merged against the default branch. Just configure these in the
  default branch's tinyci.yml. See types/task.go's RepoConfig struct for more
  information.
  - additionally, you can get this behavior by specifying the "do not merge"
    flag in the same place. It applies to all branches and PRs.
  - currently, only the overlay-runner supports this functionality.
- `tinycli` now has a splash of color in certain spots. Feedback please! Set
  `TINYCI_NOCOLOR` to something in your environment to turn it off.
- per-run environment properties. Add to each run in task.ymls to utilize.
- Using golang 1.16
- (optional!) privileged mode for runners.
- unified `tinyci` binary for launching all services.
- Many internal build improvements & upgraded dependencies

## 0.2.6 -- Thu 16 Jul 2020 05:06:39 AM UTC

This release contains many small bug fixes to the underlying components:

- hooksvc
- errors framework

It also supports the notion of resource constraints which are a part of the
RunSettings protocol now. New runners will need to address these values if they
want to use them.

The vendor tree was pruned. Use go 1.14 and modules now.

## 0.2.5 -- Mon 06 Jul 2020 05:20:25 PM UTC

0.2.5 allows changes to come from non-master branches

## 0.2.0 -- Mon Nov 04 09:11:42 PDT 2019

0.2.0 represents a lot of bugfixes and a few large features:

- Submissions! Submissions group your tasks and runs into a single item that
  corresponds for the submission of the test in the first place. Submissions
  come from:
  - A POST to hooksvc from github (pull request, branch push, stuff like that)
  - Manual submissions in the UI or through tinycli
- Auth has been broken out into its own service, in preparation for our multi-service journey.
  - An additional, but unused service called the `reposvc` is present in the release. Ignore it.
- The UI has received a significant overhaul as well; several changes were made
  to facilitate that. Please see [ci-ui CHANGELOG.md](https://github.com/tinyci/ci-ui/blob/master/CHANGELOG.md)
  for more info.
- GRPC services now gracefully stop, attempting to finish their transactions
  before terminating.
- Certain classes of transactional / distributed bugs were squashed due to
  changes in use of golang contexts.
- Pull Request IDs are now called "Ticket IDs" in the UI, and are reflected in the UI.
- Now on golang 1.13 and targeting ubuntu disco for installations
- Run data now carries what host it ran on if it has entered a running state.

## 0.1.1 -- Sat Jul 13 10:27:26 PDT 2019

Just a patch release intended to fix some things with the 0.1.0 launch.

- Update tinycli to use certs if necessary/desired
- Clear session cookies when they're invalid

## 0.1.0 -- Wed Jul 3 11:13:35 PDT 2019

Initial release!
