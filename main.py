import os
import discord
from discord import app_commands as commands

intents = discord.Intents.default()
intents.voice_states = True

client = discord.Client(intents=intents)
tree = commands.CommandTree(client)

DISCORD_TOKEN = os.environ['LOGBOT_TOKEN']
DISCORD_CHANNEL = int(os.environ['LOGBOT_CHANNEL'])

@tree.command(name = "help", description = "Show this help message.")
async def first_command(interaction):
    await interaction.response.send_message("This bot actually has no commands, please just setup the env-Variables for the DISCORD_CHANNEL & DISCORD_TOKEN!")

@client.event
async def on_ready():
    await tree.sync()
    print('Logged in as {0.user}'.format(client))

@client.event
async def on_voice_state_update(member, before, after):
    if before.channel != after.channel:
        if before.channel is not None:
            log_channel = client.get_channel(DISCORD_CHANNEL) # replace with your own text channel ID
            await log_channel.send(f'{member.mention} has left voice channel <#{before.channel.id}>.')
        if after.channel is not None:
            log_channel = client.get_channel(DISCORD_CHANNEL) # replace with your own text channel ID
            await log_channel.send(f'{member.mention} has joined voice channel <#{after.channel.id}>.')

client.run(DISCORD_TOKEN)
