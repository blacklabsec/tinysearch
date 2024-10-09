# TinySearch CLI

A command-line tool that uses the [Google Programmable Search API](https://developers.google.com/custom-search/v1/overview) to query web pages and return relevant URLs. This is made for performing a couple quick searches and getting the URLs

## Features

- Query a specific domain or a list of domains.
- Add custom search queries to refine the results.
- Automatically handles pagination for up to 100 results per query.
- Configurable maximum results per query.
- Uses a configuration file to store API keys securely.

## Installation
```
go install -v github.com/blacklabsec/tinysearch@1.0.1
```

or

1. **Clone the repository**:
   ```sh
   git clone git@github.com:blacklabsec/tinysearch.git
   cd tinysearch
   ```

2. **Build the binary**:
   ```sh
   go build -o tinysearch tinysearch.go
   ```

## Configuration:
   Create a configuration file at `~/.config/tinysearch/config.json` with the following format:
   ```json
   {
     "cx": "YOUR_CX_HERE",
     "api_key": "YOUR_API_KEY_HERE"
   }
   ```

   Replace `YOUR_CX_HERE` and `YOUR_API_KEY_HERE` with your Google Programmable Search Engine CX and API Key respectively.

## Usage

### Command-Line Options

- **`-u`**: Domain or URL to search (e.g., `-u example.com`).
- **`-l`**: File containing a list of domains or URLs to search (e.g., `-l domains.txt`).
- **`-q`**: Additional query to search with the domain (e.g., `-q "security news"`).
- **`-m`**: Maximum number of results to query (default 100) (e.g., `-m 50`).
- **`-c`**: Path to the configuration file (default: `~/.config/tinysearch/config.json`).

### Examples

1. **Search a specific domain**:
   ```sh
   ./tinysearch -u example.com
   ```

2. **Search a list of domains**:
   ```sh
   ./tinysearch -l domains.txt -q "vulnerabilities"
   ```

3. **Limit the maximum number of results**:
   ```sh
   ./tinysearch -u example.com -m 50
   ```


## License

This project is licensed under the MIT License.
