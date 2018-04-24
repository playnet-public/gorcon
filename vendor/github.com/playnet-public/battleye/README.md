## BattlEye - Golang Package
[![Go Report Card](https://goreportcard.com/badge/github.com/playnet-public/battleye)](https://goreportcard.com/report/github.com/playnet-public/battleye)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/513590eff4e54095a25b66bf65bd1323)](https://www.codacy.com/app/playnet/battleye?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=playnet-public/battleye&amp;utm_campaign=Badge_Grade)
[![Build Status](https://travis-ci.org/playnet-public/battleye.svg?branch=master)](https://travis-ci.org/playnet-public/battleye)
[![GitHub license](https://img.shields.io/badge/license-AGPL-blue.svg)](https://raw.githubusercontent.com/playnet-public/battleye/master/LICENSE)
[![Join Discord at https://discord.gg/dWZkR6R](https://img.shields.io/badge/style-join-green.svg?style=flat&label=Discord)](https://discord.gg/dWZkR6R)

The BattlEye package is implementing the [BattlEye Protocol](https://www.battleye.com/downloads/BERConProtocol.txt) for use in other projects.

## Status

This package is being seen as feature complete and changes should not occur.
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