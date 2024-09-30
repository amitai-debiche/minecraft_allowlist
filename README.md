# Minecraft Allowlist(Whitelist) Automation 

[![grassblock](https://cdn3.emoji.gg/emojis/grassblock.png)](https://emoji.gg/emoji/grassblock)
-------------------------------------------------------------------------------------------
This ReadME needs some cleanup but for now some basic guidelines:
See frontend setup at [https://minecraft.the-styx.net](https://minecraft.the-styx.net)

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



## Updates / Coming Soon
- Setup script for server agent &check;
- Container to support backend + setup for database &check;
- Setup docker container for discord &check;
- Docker Container for Frontend nginx configuration &check;
- Make client web app more efficient/ image not so slow to load
- Add gateway functionality, if username exists and wasn't recent request, send another discord message
- Expand discord bot functionality to account for specific reactions
- Docker Compose/setup script to launch step 2-4 automatically (step 1 will still have to be done on vm running MC server tmux session)
- Making README's Pretty/More useful :)
- Add a check on agent to see if user was successfully added with whitelist command(sometimes this gave error player doesn't exists, not sure why, then manually entering name worked)
  
