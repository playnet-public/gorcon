## BattlEye - Golang Package
[![Go Report Card](https://goreportcard.com/badge/github.com/playnet-public/battleye)](https://goreportcard.com/report/github.com/playnet-public/battleye)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/513590eff4e54095a25b66bf65bd1323)](https://www.codacy.com/app/playnet/battleye?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=playnet-public/battleye&amp;utm_campaign=Badge_Grade)
[![Build Status](https://travis-ci.org/playnet-public/battleye.svg?branch=master)](https://travis-ci.org/playnet-public/battleye)
[![GitHub license](https://img.shields.io/badge/license-AGPL-blue.svg)](https://raw.githubusercontent.com/playnet-public/battleye/master/LICENSE)
[![Join Discord at https://discord.gg/dWZkR6R](https://img.shields.io/badge/style-join-green.svg?style=flat&label=Discord)](https://discord.gg/dWZkR6R)

The BattlEye package is implementing the [BattlEye Protocol](https://www.battleye.com/downloads/BERConProtocol.txt) for use in other projects.

## Status

### Version 2

The previous version of this package was refactored with a lot of breaking changes. The reason for this was lack of abstraction which made testing in other projects harder.

As we are using TDD/BDD and strafe towards a high level of coverage and overall code quality, we had to make some changes to this.
The most notable one is that we moved all functionality used by clients themselves into the `Protocol` interface.
From there we mostly copied the old implementation and created wrappers where necessary.
While we were at it, there were also some improvements to both the api and type safety.

Using the interface allows us to generate mocks from it to enable better testing in other projects without having to mess with real BE packets (see [mocks/](mocks/)).

This new version might still be subject to improvements and changes that could break the API (even further). We suggest either pinning the state you need by using dependency management tools, or to use the v1 tag if the new features are not important to you.

TODO's:
- All tests have just been ported to the new structure which results in some legacy test code no longer up to our standards. This should be addressed in the near future and might also offer a good introduction to the project for newcomers.

### Version 1

NOTE: Version 1 is the legacy implementation which works but lacks some abstraction as well as mocks to allow best use in other projects.
There is no change in features/functionality when upgrading but the overall api and signature was drastically changed.

The git tag v1 is being seen as feature complete and changes should not occur. 
If bugs in this implementation are found, please report those using GitHub Issues or by opening a Pull Request.
Whilst being feature complete a reasonable request to this might still make it into the code.

The code is covered with tests by 100% and code quality is being ensured by running go test, vet, errcheck, fmt and lint against every commit.
Checks are run against Go 1.8.x, 1.9.x and 1.10.x.

## Dependencies
All dependencies inside this project are being managed by [dep](https://github.com/golang/dep) and are checked in.
After pulling the repository, it should not be required to do any further preparations aside from `make deps` to prepare the dev tools (once).

If new dependencies get added while coding, make sure to add them using `dep ensure --add "importpath"` and to check them into git.
We recommend adding your vendor changes in a separate commit to make reviewing your changes easier and faster.

## Testing
To run tests you can use:
```bash
make test
```

## Contributing

Feedback and contributions are highly welcome. Feel free to file issues, feature or pull requests.
If you are interested in using this project now or in a later stage, feel free to get in touch.

We are always looking for active contributors, team members and partner projects sharing our vision.
Easiest way of reaching us is via [Discord](https://discord.gg/dWZkR6R).

See you soon,
the PlayNet Team

## License

This project's license is located in the [LICENSE file](LICENSE).