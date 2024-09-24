from typing import Final
import os
from dotenv import load_dotenv
import discord
from discord import Intents, Client, Message
from discord.ext import commands
import requests
import asyncio
from flask import Flask, request, jsonify

load_dotenv()
TOKEN: Final[str] = os.getenv('DISCORD_TOKEN', '')
API_KEY: Final[str] = os.getenv('API_KEY', '')
TRUSTED_ROLE_NAME: str = os.getenv("TRUSTED_ROLE_NAME", 'trusted') # defualt role to trusted
BACKEND_URL = os.getenv("BACKEND_URL", 'localhost')

if TOKEN == '':
    raise ValueError("DISCORD_TOKEN env variable has not been set")
if API_KEY == '':
    raise ValueError("API_KEY env variable has not been set")


intents: Intents = Intents.default()
intents.message_content = True
channel = None
bot_messages = {}
bot = commands.Bot(command_prefix='$', intents=intents)


app = Flask(__name__)

@bot.command()
async def set_channel(ctx):
    global channel
    channel = ctx.channel
    print(channel)
    await ctx.send("SETUP: Channel has been set for receiving bot authorization messages")

async def send_message_to_channel(username: str, message: str):
    global channel, bot_messages
    if not username:
        print("Message empty, no username submittedf")
        return False, "No username provided"

    try:
        output: str = f"Request for {username} to be added to the-styx!"
        if message is not None:
            output += f"\nUser included the following message: `{message}`"
        if channel != None:
            sent_message = await channel.send(output)
            bot_messages[sent_message.id] = username
        else:
            return False, "channel not declared"


        return True, None
    except Exception as e:
        return False, str(e)


@app.route('/send_message', methods=['POST'])
def send_message():
    print("send message called")
    auth_header = request.headers.get('Authorization')
    if auth_header is None or auth_header != f'Bearer {API_KEY}':
        return jsonify({"error": "Unauthorized"}), 403

    data = request.json
    if data is None:
        return jsonify({"error": "Invalid or missing JSON"}), 400

    username = data.get('username')
    message = data.get('message')

    if not message or not username:
        return jsonify({"error": "Invalid input"}), 400

    # using asyncio to send from synchronous to async function
    success, error = asyncio.run_coroutine_threadsafe(
            send_message_to_channel(username, message), bot.loop).result()
    if success:
        print("success")
        return jsonify({"message": "Message sent successfully"}), 200
    else:
        print(error)
        return jsonify({"error": error}), 500


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
            response = requests.post(
                f"{BACKEND_URL}/approveUsername",
                json={"username": username, "status": True},
                headers={'Authorization': f'Bearer {API_KEY}'}
            )

            if response.status_code == 200:
                print(f"Successfully sent {username} to the backend.")
            else:
                print(f"Failed to send {username}. Status code: {response.status_code}")
        except Exception as e:
            print(f"Error while sending username to backend: {e}")

#bot.event
async def on_ready() -> None:
    print(f'{bot.user} is now running!')

def run_flask_app():
    app.run(port=5000, debug=False)

def main() -> None:
    from threading import Thread
    flask_thread = Thread(target=run_flask_app)
    flask_thread.start()

    bot.run(TOKEN)

if __name__ == '__main__':
    main()
