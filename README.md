# CSE-416

Requirements:
go, node, and python must be installed to PATH

## Steps to Install and Run:

1. From the root directory, run the following command:
```bash
npm run start
```
This will handle all builds and start Electron. Please note that this process may take some time.

2. To start the server separately, use:
```bash
npm run start:server
```

3.If you encounter issues running npm run start on macOS run:
```bash
cd client
npm run run-build-test
```
4. If the server does not start, navigate to the application-layer directory and run the following commands:
```bash
cd application-layer
go run main.go
go run fileAndProxy/testMain.go
```

