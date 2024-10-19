# CSE-416

## Installation Instructions:

### Electron Desktop Installation Instructions

Follow these steps to set up the Electron desktop application:
```
git clone 
```
1. Clone the repository and cd into it
```
https://github.com/Mayukh-Banik/CSE-416.git
cd ./CSE-416
```

2. Navigate to the client directory, install dependencies, and build the React app:
    ```
    cd client
    npm install
    npm install electron
    npm run build
    ```
    This will install all necessary client dependencies and create a `build` folder with the `index.html` file inside it.

3. Go back to the root directory, install additional dependencies, and launch the Electron app:
    ```
    cd ..
    npm install
    npm run build
    npm run electron
    ```
