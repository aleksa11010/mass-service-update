# Harness YAML String Replacement Utility

## Overview
This utility tool is designed to facilitate the process of replacing specific strings within Service YAML.

## Features
- Fetch all projects and services from a specified Harness account.
- Replace specified target strings with replacement strings in Service YAML files.
- Update the modified YAML configurations back to the Harness platform.

## Prerequisites
- Go programming language
- Access to a Harness account with appropriate permissions
- Harness API Key

## Installation
Clone the repository to your local machine:

```bash
git clone [REPOSITORY_URL]
cd [REPOSITORY_DIRECTORY]
```

Ensure that Go is installed on your system. You can download and install Go from [https://golang.org/dl/](https://golang.org/dl/).

## Usage

### Setting up Command Line Arguments
The utility accepts the following command line arguments:
- `account`: Your Harness account ID.
- `api-key`: Your Harness API Key.
- `target`: The string you want to replace in the YAML files.
- `replacement`: The string to replace the target with.

### Running the Tool
Execute the program with the required arguments:

```bash
go run main.go -account="YOUR_ACCOUNT_ID" -api-key="YOUR_API_KEY" -target="TARGET_STRING" -replacement="REPLACEMENT_STRING"
```

### Example
To replace `"old-service-url"` with `"new-service-url"` in all services:

```bash
go run main.go -account="abcd1234" -api-key="xyz789" -target="old-service-url" -replacement="new-service-url"
```
