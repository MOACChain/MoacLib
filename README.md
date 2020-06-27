[![Documentation](https://godoc.org/github.com/MOACChain/MoacLib?status.svg)](http://godoc.org/github.com/MOACChain/MoacLib)

## Go MoacLib dev

Golang implementation of the MoacLib library used in MOAC project.

[![API Reference]](https://godoc.org/github.com/MOACChain/MoacLib)

Automated builds are available for stable releases and the unstable master branch.
Binary archives are published at https://github.com/MOACChain/moac-core/releases/.

## Test the package

To test the packge, simply run:
go test
under each src directory.

## Building the package

Building the package requires a Go (version 1.7 or later) and some vendors.
You can install them using your favourite package manager.
Once the dependencies are installed, run

    go install

the package should be installed under $GOPATH/package


## Contribution

Thank you for considering to help out with the source code! We welcome contributions from
anyone on the internet, and are grateful for even the smallest of fixes!

If you'd like to contribute to MoacLib, please fork, fix, commit and send a pull request
for the maintainers to review and merge into the main code base. If you wish to submit more
complex changes though, please check up with the core devs first on [our github](https://github.com/MOACChain/MoacLib)
to ensure those changes are in line with the general philosophy of the project and/or get some
early feedback which can make both your efforts much lighter as well as our review and merge
procedures quick and simple.

Please make sure your contributions adhere to our coding guidelines:

 * Code must adhere to the official Go [formatting](https://golang.org/doc/effective_go.html#formatting) guidelines (i.e. uses [gofmt](https://golang.org/cmd/gofmt/)).
 * Code must be documented adhering to the official Go [commentary](https://golang.org/doc/effective_go.html#commentary) guidelines.
 * Pull requests need to be based on and opened against the `master` branch.
 * Commit messages should be prefixed with the package(s) they modify.
   * E.g. "mc, rpc: make trace configs optional"

Please see the [Developers' Guide](https://github.com/MOACChain/MoacLib/wiki/Developers'-Guide)
for more details on configuring your environment, managing project dependencies and testing procedures.

## License

The MoacLib library is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html), also
included in our repository in the `COPYING.LESSER` file.

