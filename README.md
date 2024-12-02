# CSE-416

## Installation Instructions:

1. Clone GitHub Repository
    
    ```
    git clone https://github.com/Mayukh-Banik/CSE-416.git
    ```

2. Install Docker (Restart System if necessary, follow install prompts)
    
        https://docs.docker.com/engine/install/

3. Run docker startup command with admin privileges

    Linux/MacOs

        sudo docker compose up --build

    Windows - Powershell as Administrator

        docker compose up --build

4. The website will be on port 3000 on local system, with server located on port 8000.
This will apply to the docker image as well as on the host machine.

5. Database folder will be located on project's root directory with database only being accesible to the server which will have a reference to it in it's "\database" directory.

6. For server run:
   ```
   sudo docker compose up --build --detach
   sudo docker exec ipfs_server cp swarm.key Shared
   npm run dev
   ```

### Notes


Database folder will be populated with a lot of stuff, can be ignored for now, need to figure out database solution


### Electron Desktop Installation Instructions

Follow these steps to set up the Electron desktop application:
```
    cd client
    npm install
    npm run run-build
```


### notes for joyce
set up bootstrap node locally
copy multiaddr and peerid and use that as bootstrap addr in application-layer/dht