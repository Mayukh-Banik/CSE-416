# Getting Started with Create React App

# CSE-416

## Installation Instructions:

1. Development Mode

    To run your app in development mode (with live reloading):
    
    ```
    npm start
    ```

2. Install Docker (Restart System if necessary, follow install prompts)
    
        https://docs.docker.com/engine/install/

3. Run docker startup command with admin privileges

    Linux/MacOs

        sudo docker compose up

    Windows - Powershell as Administrator

        docker compose up

4. The website will be on port 3000 on local system, with server located on port 8000.
This will apply to the docker image as well as on the host machine.

5. Database folder will be located on project's root directory with database only being accesible to the server which will have a reference to it in it's "\database" directory.

6. For server run:
   ```
   sudo docker compose up --build --detach
   sudo docker exec ipfs_server cp swarm.key Shared 
   ```
        

### Notes


Database folder will be populated with a lot of stuff, can be ignored for now, need to figure out database solution
