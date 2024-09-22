from typing import Final
import os
from dotenv import load_dotenv
import discord
from discord import Intents, Client, Message
from discord.ext import commands
import requests

load_dotenv()
TOKEN: Final[str] = os.getenv('DISCORD_TOKEN', '')
API_KEY: Final[str] = os.getenv('API_KEY', '')
TRUSTED_ROLE_NAME: str = os.getenv("TRUSTED_ROLE_NAME", 'trusted') # defualt role to trusted

if TOKEN is '':
    raise ValueError("DISCORD_TOKEN env variable has not been set")
if API_KEY is '':
    raise ValueError("API_KEY env variable has not been set")


intents: Intents = Intents.default()
intents.message_content = True
intents.reactions = True
intents.members = True
channel = None
bot_messages = {}
bot = commands.Bot(command_prefix='$', intents=intents)

@bot.command()
async def set_channel(ctx):
    global channel
    channel = ctx.channel
    await ctx.send("SETUP: Channel has been set for receiving bot authorization messages")

async def send_message(username: str, message: str):
    global channel, bot_messages
    if not username:
        print("Message empty, no username submittedf")
        return

    try:
        output: str = f"Request for {username} to be added to the-styx!"
        if message is not None:
            output += f"\nUser included the following message: `{message}`"
        if channel != None:
            sent_message = await channel.send(output)
            bot_messages[sent_message.id] = username

    except Exception as e:
        print(e)

@bot.event
async def on_reaction_add(reaction, user): 
    global bot_messages
    if reaction.message.id not in bot_messages:
        return

    def check(user):
        role = discord.utils.get(user.guild.roles, name=TRUSTED_ROLE_NAME)

        return role in user.roles

    if check(user):
        username = bot_messages[reaction.message.id]

        try:
            response = requests.post("https://api.minecraft.the-styx.net/approveUsername", json={"username": username})
            if response.status_code == 200:
                print(f"Successfully sent {username} to the backend.")
            else:
                print(f"Failed to send {username}. Status code: {response.status_code}")
        except Exception as e:
            print(f"Error while sending username to backend: {e}")

#bot.event
async def on_ready() -> None:
    print(f'{bot.user} is now running!')


def main() -> None:
    bot.run(TOKEN)

if __name__ == '__main__':
    main()
