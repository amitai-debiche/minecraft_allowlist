from typing import Final
import os
from dotenv import load_dotenv
import discord
from discord import Intents, Client, Message
from discord.ext import commands
import requests






load_dotenv()
TOKEN: Final[str] = os.getenv('DISCORD_TOKEN', '')
if TOKEN is '':
    raise ValueError("DISCORD_TOKEN env variable has not been set")

intents: Intents = Intents.default()
intents.message_content = True
client: Client = Client(intents=intents)
channel = client.get_channel(12345)

bot = commands.Bot(command_prefix='$', intents=intents)

@bot.command()
async def set_channel(ctx):
    channel = ctx.channel
    await ctx.send("SETUP: Channel has been set for receiving bot authorization messages")



async def send_message(username: str, message: str) -> None:
    if not username:
        print("Message empty, no username submittedf")
        return

    try:
        output: str = f"Request for {username} to be added to the-styx!"
        if message is not None:
            output += f"\nUser included the following message: `{message}`"
        sent_message = await channel.send(output)
        message_id = sent_message.id 
        reaction = await bot.wait_for

    except Exception as e:
        print(e)

@bot.event
async def on_message(msg): 
    def check(reaction, user):
        role = discord.utils.get(msg.guild.roles, name="trusted")
        return (role in reaction.roles)

    if msg.author != bot.user or "SETUP" in msg.content:
        return 
    await bot.wait_for("reaction_add", check=check)
    ## SEND REQUEST BACK TO BACKEND


    r.POST("/approveUsername", handleUserApproval)

##https://stackoverflow.com/questions/65090668/reaction-check-with-discord-py
