# NXG utils

Nxg utils to build nxg links

# What is NXG?

Is a new way to build Usenet message ID's that allows to share NZB files without the need of an actual .nzb file. Instead a link with the NXG header is shared and the client can use that to download all associated files from usenet.

See [NXG](https://github.com/Tensai75/nxg-upper/tree/main?tab=readme-ov-file#advantages-of-the-nxg-header)

## Development Setup

To set up the project for development, follow these steps:

1. Clone the repository:

```sh
git clone https://github.com/javi11/nxg.git
cd nxg
```

2. Install dependencies:

```sh
go mod download
```

3. Run tests:

```sh
make test
```

4. Lint the code:

```sh
make lint
```

5. Generate mocks and other code:

```sh
make generate
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request. See the [CONTRIBUTING.md](CONTRIBUTING.md) file for details.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
