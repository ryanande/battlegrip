# Welcome to battlegrip 👋

![Version](https://img.shields.io/badge/version-0.0.1-blue.svg?cacheSeconds=2592000)
![Prerequisite](https://img.shields.io/badge/golang-%5E1.17-blue)
[![License: MIT](https://img.shields.io/github/license/ryanande/battlegrip)](https://github.com/ryanande/battlegrip/blob/master/LICENSE)
[![Twitter: ryanande](https://img.shields.io/twitter/follow/ryanande.svg?style=social)](https://twitter.com/ryanande)

Battlegrip is a simple companion library to add a little help to your user experience when using [Cobra](https://github.com/spf13/cobra) for Go CLI interactions.

The core value proposition is streamlining the creation of multiflag commands in a simple web base UI.

When incorporated into your Go CLI, the library, once commanded, will perform the following;

* launch a basic local web server,
* launch your default browser, navigating to the server,

All CLI commands will be displayed, along with ech commands flags in a smooth form based web page. From there you can create your terminal commands in a clean web form + copy paste feature set.

## Features

* Easy 1 line setup
* Clean meta driven web page
* Simple form based style command builder

## Install

Grab the lib,

```sh
go get github.com/ryanande/battlegrip
```

Import the library into your root command module,

```go
import (
    "github.com/ryanande/battlegrip"

    ...
)
```

Where you register your commands, simply register the battlegrip command,

```go
rootCmd.AddCommand(battlegrip.UICmd) // battlegrip setup
```

Now, the `ui` command should be available from your assembly,

```sh
./myassembly ui
```

## Author

👤 **Ryan Anderson**

* Twitter: [@ryanande](https://twitter.com/ryanande)
* Github: [@ryanande](https://github.com/ryanande)

## 🤝 Contributing

Contributions, issues and feature requests are welcome!

Feel free to check [issues page](https://github.com/ryanande/battlegrip/issues). You can also take a look at the [contributing guide](https://github.com/ryanande/battlegrip/blob/master/CONTRIBUTING.md).

## Show your support

Give a ⭐️ if this project helped you!

## 📝 License

Copyright © 2021 [Ryan Anderson](https://github.com/ryanande).

This project is [MIT](https://github.com/ryanande/battlegrip/blob/master/LICENSE) licensed.
