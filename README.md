# Minecraft Allowlist(Whitelist) Automation 

[![grassblock](https://cdn3.emoji.gg/emojis/grassblock.png)](https://emoji.gg/emoji/grassblock)
-------------------------------------------------------------------------------------------
This ReadME needs some cleanup but for now some basic guidelines:

## Setup
Currently setup is a bit annoying as no automated scripts/docker YET

1. Setup the server agent by installing agent.go, (binary & install script coming soon)
    - Change your env settings to fit your needs

2. Now setup your backend
   - You will need to create a API_Key that will be shared between the frontend, the discord bot and the backend for authorization to the backend
   - You will also need to create a postgres database yourself, automation has not been added for this yet (that will be done asap)

3. Setup discord bot
  - Fill out env, to api token, you will need to go to your bot dashboard and copy token

4. Finally launch backend
  - This is probably done with the most ease using nginx as your engine, a docker container to serve the frontend will be added soon



## Coming Soon
- Setup script for server agent w/systemd configuration for persistence
- Container to support backend + setup for database
- Setup script for discord(debating to have this as container or not)
- Docker Container for Frontend nginx configuration
- Docker Compose/setup script to launch step 2-4 automatically (step 1 will still have to be done on vm running MC server tmux session)
- Making README's Pretty/More useful :)
  
