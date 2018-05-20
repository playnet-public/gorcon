## GoRcon
[![Go Report Card](https://goreportcard.com/badge/github.com/playnet-public/gorcon)](https://goreportcard.com/report/github.com/playnet-public/gorcon)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/513590eff4e54095a25b66bf65bd1323)](https://www.codacy.com/app/playnet/gorcon?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=playnet-public/gorcon&amp;utm_campaign=Badge_Grade)
[![Build Status](https://travis-ci.org/playnet-public/gorcon.svg?branch=master)](https://travis-ci.org/playnet-public/gorcon)
[![Docker Repository on Quay](https://quay.io/repository/playnet/gorcon/status "Docker Repository on Quay")](https://quay.io/repository/playnet/gorcon)
[![GitHub license](https://img.shields.io/badge/license-AGPL-blue.svg)](https://raw.githubusercontent.com/playnet-public/gorcon/master/LICENSE)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/playnet-public/gorcon)
[![Join Discord at https://discord.gg/dWZkR6R](https://img.shields.io/badge/style-join-green.svg?style=flat&label=Discord)](https://discord.gg/dWZkR6R)

GoRcon is an abstraction layer to map common rcon functionality from various games to a unified api (grpc/rest) while also providing local management functionality to make running game servers easier.
This allows the use as a rcon connector for both client tools and cloud solutions like the PlayNet BanList.

## Why

At PlayNet it is one of our core motivation to make development in the gaming area better and open up new ways for the community.
Developing and offering standards with and for the community as well as building tools required by it or our own products is key to that vision.

Most Rcon protocols are either badly documented or simply bad to handle. While most offer the same functionality it is quite hard to reuse features and gaming communities rarely cooperate on this.
We see ourselves in duty here to both build a standardized protocol everybody can freely use, as well as bringing the gaming community together as a whole to work on projects like this.

Our primary goal at first is to provide a basic server tool for the ArmA game series which is capable enough of replacing BEC while developing a featureset rich enough to allow us building the community banlist we have in mind with PlayNet BanList.
Implementation of other rcon protocols like Minecraft and Valve are on the roadmap but we got to set priorities due to limited resources.

This project is a rewrite of GoRcon-ArmA which was the first proof-of-concept for the idea.

## Concepts

As pretty obvious from looking at the repository this project is realized in Golang. Mainly for it's portability to almost any platform as well as it's great language features and codestyle.

The project consists of a main application responsible for handling the cli and holding everything together, while game specific implementations reside in their respective locations.
While nothing is set in stone yet, the application is meant to consist of various parts:
* CLI for starting an managing the application
* Process Manager for starting and watching a defined (game) process
* Rcon Connection for communicating with the game servers
* Scheduler for invoking functions based on scripts and timed logic
* API Endpoints for configuring the application as well as invoking functions provided by other parts

## Coding and Style

Coding is done using pull requests and code reviews. Master is locked.
Our code is always checked by Travis using `make test check` therefor all Golang rules on syntax and formating have to be met for pull requests to be merged.
While this might incur more work for possible contributors, we see the code produced here as production critical once finished and therefor strive for high code quality.

The team is developing this mostly using TDD and BDD. If you don't know what this is, we recommend this [video](https://www.youtube.com/watch?v=uFXfTXSSt4I) for starters or to [get in touch](https://discord.gg/dWZkR6R) and let us introduce you.

Please do reasonable commit sizes.

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
If you are developing or already finished solutions that rely on rcon in any way, we would be happy to talk to you for both gaining insights as well as looking for options to collaborate and involve your project into other endeavors.

We are always looking for active contributors, team members and partner projects sharing our vision.
Easiest way of reaching us is via [Discord](https://discord.gg/dWZkR6R).

See you soon,
the PlayNet Team

## Attributions

* [Kolide for providing `kit`](https://github.com/kolide/kit)

## License

This project's license is located in the [LICENSE file](LICENSE).
When forking or copying this project or it's parts we ask you to redirect people back to this repository for keeping the community together and to help getting contributions to make this project better everyday.