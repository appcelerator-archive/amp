# Contributing to AMP

Want to contribute to AMP? We think you're awesome! This page should help you get started.

## Topics

* [Code of Conduct](#code-of-conduct)
* [Reporting security issues](#reporting-security-issues)
* [Reporting other issues](#reporting-other-issues)
* [Quick contribution tips and guidelines](#quick-contribution-tips-and-guidelines)
* [Coding style](#coding-style)

## Code of Conduct

Before becoming involved, please read our [CODE-OF-CONDUCT](CODE-OF-CONDUCT.md).

## Reporting security issues

The project maintainers take security seriously. If you discover a security issue, please bring it to their
attention right away!

Please **DO NOT** file a public issue. Instead, send your report privately to
[security@appcelerator.io](mailto:security@appcelerator.io).

## Reporting other issues

Helping to provide detailed reports on issues is a great way to contribute to the project and it
is highly appreciated.

Before reporting an issue, however, please check that it hasn't already been filed in our
[issue tracker](https://github.com/appcelerator/amp/issues). If you find one already open,
please feel free to add any extra information you feel adds value in helping to resolve it.
You can use the "subscribe" button to get update notifications. However, please do *not*
leave "+1" type of comments that only clutter the discussion and don't help resolve the
issue itself. GitHub now provides the ability to add [reactions](https://github.com/blog/2119-add-reactions-to-pull-requests-issues-and-comments)
to issues and comments, if you want to show your enthusiasm (or lack thereof) for anything
in particular.

### When reporting an issue

Please include:

  * The output of `docker version`.
  * The output of `docker info`.
  * The output of `amp --version`.
  * The output of `amp --config`.
  * Steps to reproduce the problem (if possible and applicable).

If sending lengthy output such as logs, it is preferred that you create a
[gist](https://gist.github.com) and provide the link. Remove any sensitive data
before posting anything (replace with "REDACTED").

## Quick contribution tips and guidelines

### Pull requests are always welcome

We welcome and appreciate contributors who want to help refactor code, fix bugs, and submit features.
With respect to submitting any improvements, these should first be documented as an [issue](https://github.com/appcelerator/amp/issues)
before work is started.

### Conventions

Fork the repository and make changes on your fork in a feature branch:

- If it's a bug fix branch, name it `XXXX-something` where `XXXX` is the number of
	the issue. 
- If it's a feature branch, create an enhancement issue to announce
	your intentions, and name it `XXXX-something` where `XXXX` is the number of the
	issue.

Submit unit tests for your changes. Go has a great test framework built in; use
it! Take a look at existing tests for inspiration. Run the full test
suite (`make test`) on your branch before submitting a pull request.

Update the documentation when creating or modifying features. Test your
documentation changes for clarity, concision, and correctness, as well as a
clean documentation build. For a good set of style and grammar conventions,
see Docker's guide [here](https://docker.github.io/opensource/doc-style/).

Write clean code. Universally formatted code promotes ease of writing, reading,
and maintenance. Always run `gofmt -s -w file.go` on each changed file before
committing your changes. Most editors have plug-ins that do this automatically.

Pull request descriptions should be as clear as possible and include a reference
to all the issues that they address.

Commit messages must start with a capitalized and short summary (max. 50 chars)
written in the imperative, followed by an optional, more detailed explanatory
text which is separated from the summary by an empty line.

Code review comments may be added to your pull request. Discuss, then make the
suggested modifications and push additional commits to your feature branch. Post
a comment after pushing. New commits show up in the pull request automatically,
but the reviewers are notified only when you comment.

Pull requests must be cleanly rebased on top of master without multiple branches
mixed into the PR.

Before you make a pull request, squash your commits into logical units of work
using `git rebase -i` and `git push -f`. A logical unit of work is a consistent
set of patches that should be reviewed together: for example, upgrading the
version of a vendored dependency and taking advantage of its now available new
feature constitute two separate units of work. Implementing a new function and
calling it in another file constitute a single logical unit of work. The very
high majority of submissions should have a single commit, so if in doubt: squash
down to one.

Include documentation changes in the same pull request so that a revert would
remove all traces of the feature or fix.

Include an issue reference like `Closes #XXXX` or `Fixes #XXXX` in commits that
close an issue. Including references automatically closes the issue on a merge.

Please do not add yourself to the `AUTHORS` file, as it is regenerated regularly
from the Git history.

Please see the [Coding Style](#coding-style) section for further guidelines.

### Merge approval

Maintainers use LGTM (Looks Good To Me) in comments on the code review to
indicate acceptance.

A change requires LGTMs from an absolute majority of the maintainers of each
component affected. For example, if a change affects `docs/` and `registry/`, it
needs an absolute majority from the maintainers of `docs/` AND, separately, an
absolute majority of the maintainers of `registry/`.

For more details, see the [MAINTAINERS](MAINTAINERS) page.

### Sign your work

The sign-off is a simple line at the end of the explanation for the patch. Your
signature certifies that you wrote the patch or otherwise have the right to pass
it on as an open-source patch. The rules are pretty simple: if you can certify
the below (from [developercertificate.org](http://developercertificate.org/)):

The nice thing about this method is that contributors do not need to execute
an individual Contributor License Agreement (CLA) to participate. This is a
great idea that originated with the Linux Kernel Project is and is used by the
Docker project itself.

```
Developer Certificate of Origin
Version 1.1

Copyright (C) 2004, 2006 The Linux Foundation and its contributors.
1 Letterman Drive
Suite D4700
San Francisco, CA, 94129

Everyone is permitted to copy and distribute verbatim copies of this
license document, but changing it is not allowed.

Developer's Certificate of Origin 1.1

By making a contribution to this project, I certify that:

(a) The contribution was created in whole or in part by me and I
    have the right to submit it under the open source license
    indicated in the file; or

(b) The contribution is based upon previous work that, to the best
    of my knowledge, is covered under an appropriate open source
    license and I have the right under that license to submit that
    work with modifications, whether created in whole or in part
    by me, under the same open source license (unless I am
    permitted to submit under a different license), as indicated
    in the file; or

(c) The contribution was provided directly to me by some other
    person who certified (a), (b) or (c) and I have not modified
    it.

(d) I understand and agree that this project and the contribution
    are public and that a record of the contribution (including all
    personal information I submit with it, including my sign-off) is
    maintained indefinitely and may be redistributed consistent with
    this project or the open source license(s) involved.
```

Then you just add a line to every git commit message:

    Signed-off-by: Joe Smith <joe.smith@email.com>

Use your real name (sorry, no pseudonyms or anonymous contributions.)

If you set your `user.name` and `user.email` git configs, you can sign your
commit automatically with `git commit -s`.

### Becoming a maintainer

The procedures for adding new maintainers are explained in the
[MAINTAINERS](MAINTAINERS) document.

Keep in mind that being a maintainer is a time investment. Make sure you
will have time to make yourself available. You don't have to be a
maintainer to be a valued and appreciated contributor!

## Coding Style

Unless explicitly stated, we follow all coding guidelines from the Go
community. While some of these standards may seem arbitrary, they somehow seem
to result in a solid, consistent codebase.

It is possible that the code base does not currently comply with these
guidelines. We are not looking for a massive PR that fixes this, since that
goes against the spirit of the guidelines. All new contributions should make a
best effort to clean up and make the code base better than they left it.
Obviously, apply your best judgement. Remember, the goal here is to make the
code base easier for humans to navigate and understand. Always keep that in
mind when nudging others to comply.

The rules:

1. All code should be formatted with `gofmt -s`.
2. All code should pass the default levels of
   [`golint`](https://github.com/golang/lint).
3. All code should follow the guidelines covered in [Effective
   Go](http://golang.org/doc/effective_go.html) and [Go Code Review
   Comments](https://github.com/golang/go/wiki/CodeReviewComments).
4. Comment the code. Tell us the why, the history and the context.
5. Document _all_ declarations and methods, even private ones. Declare
   expectations, caveats and anything else that may be important. If a type
   gets exported, having the comments already there will ensure it's ready.
6. Variable name length should be proportional to its context and no longer.
   `noCommaALongVariableNameLikeThisIsNotMoreClearWhenASimpleCommentWouldDo`.
   In practice, short methods will have short variable names and globals will
   have longer names.
7. No underscores in package names. If you need a compound name, step back,
   and re-examine why you need a compound name. If you still think you need a
   compound name, lose the underscore.
8. No utils or helpers packages. If a function is not general enough to
   warrant its own package, it has not been written generally enough to be a
   part of a util package. Just leave it unexported and well-documented.
9. All tests should run with `go test` and outside tooling should not be
   required. No, we don't need another unit testing framework. Assertion
   packages are acceptable if they provide _real_ incremental value.
10. Even though we call these "rules" above, they are actually just
    guidelines. Since you've read all the rules, you now know that.

If you are having trouble getting into the mood of idiomatic Go, we recommend
reading through [Effective Go](https://golang.org/doc/effective_go.html). The
[Go Blog](https://blog.golang.org) is also a great resource. Drinking the
kool-aid is a lot easier than going thirsty.

