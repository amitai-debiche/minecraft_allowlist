from typing import Final
import os
from dotenv import load_dotenv
from discord import Intents, Client, Message






load_dotenv()
TOKEN: Final[str] = os.getenv('DISCORD_TOKEN', '')
if TOKEN is '':
    raise ValueError("DISCORD_TOKEN env variable has not been set")

intents: Intents = Intents.default()
intents.message_content = True
client: Client = Client(intents=intents)
channel = client.get_channel(12345)


async def send_message(username: str, message: str) -> None:
    if not username:
        print("Message empty, no username submittedf")
        return

    try:
        output: str = f"Request for {username} to be added to the-styx!"
        if message is not None:
            output += f"\nUser included the following message: `{message}`"
##        await channel.send(output)
    except Exception as e:
        print(e)


https://stackoverflow.com/questions/65090668/reaction-check-with-discord-py
