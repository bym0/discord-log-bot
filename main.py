import discord
import os

intents = discord.Intents.default()
intents.voice_states = True

client = discord.Client(intents=intents)

DISCORD_TOKEN = os.environ['LOGBOT_TOKEN']
DISCORD_CHANNEL = int(os.environ['LOGBOT_CHANNEL'])

@client.event
async def on_ready():
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
