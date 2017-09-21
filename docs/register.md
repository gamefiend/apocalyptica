# Registering Apocalyptica
> **IN PROGRESS! DO NOT USE!!!** Still testing out some things.

## Overview

This document covers instructions for registering apps with Discord. Looking for information on how to deploy? Go to ['Deploying Apocalyptica'](deploying.md).

Interested in test driving Apocalyptica before setting up your own? [Go to the Apocalyptica Reference install](https://apocalyptica.social-fiction.net) and invite it to your server.

## Steps to Register Apocalyptica in Discord.

Sign in to Discord and go to the [Developer page, 'My Apps'](https://discordapp.com/developers/applications/me). 
![Registering a New App](img/register_part_1.png)

Select 'New App'.

![Registration Details](img/register_part_2.png)

Give your instance of apocalyptica a name, description, and picture.

![Creating a Bot User](img/register_part_3.png)

Select 'Creat a Bot User', then on the following screen:

![Retrieving ID, Client Information](img/register_part_4.png)

Click 'click to reveal' and copy the token string. You will want to save this in your `env_discord.sample` file in the variable `DISCORD_TOKEN`.

Above the 'Create a Bot User' information is the client ID of the application. Copy that ID into `env_discord.sample` file in the variable `CLIENT_ID`.

When you have retrieved that information, save the changes.

![Save Changes](img/register_final.png)


## Invite Apocalyptica to a Server

Now that you have an app registered, you need to invite it to your server. To do that you need an authorization link.

The authorize link looks like this:
`https://discordapp.com/oauth2/authorize?&client_id=<CLIENTID>&scope=bot&permissions=0`

which will take you to a page like this:

![Authorize App](img/register_authorize.png)

Choose one of your servers from the list, and choose 'Authorize' to invite Apocalyptica to the server!
