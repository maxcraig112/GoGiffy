
# GoGiffy
A refactor of the original [Giffy Bot](https://github.com/maxcraig112/Giffy-bot) however written in Golang and utilising the Google Cloud Platform.

[Discord Invite Link](https://discord.com/oauth2/authorize?client_id=1220642492173778996&permissions=2684472384&scope=bot+applications.commands)

Please note that this bot is a work in progress and therefore current functionality is limited. Feel free to report any bugs or potential features you would like to see.
## Giffy-Bot
A multipurpose repos designed to allow the manipulation, tagging, archiving and retrieval of gifs.

This repos provides abstraction of these commands through the discord bot Giffy.

![Jif Gif Season](https://c.tenor.com/oylHwLtwhbsAAAAC/gif-jif.gif)

## Commands

### Gif Searching
| Command  | Description |
| ------------- | ------------- |
| /search tag1,tag2,...  | Searches through all currently stored gifs for the particular tags specified. This function will automatically add or remove pluralisation in order to widen the search window. Interaction is available through buttons that let you browse the returned gifs.  |

### Debugging Tool
| Command  | Description |
| ------------- | ------------- |
| /giffy  | A link to github and a general description of the bot  |
| /stats  | Returns interesting information regarding the state of the bots database  |



## FAQ
### How to run the bot?

You may have some trouble running the bot for yourself given it's reliance on GCP and some particular database names, however here are a list of some prerequisites if you want to give it a try.
1. Have a Google Cloud Platform account
    - Access to the `Cloud Vision API`, `Bigquery` and a `Bucket`
2. add your bot API key under a `token.txt` file
3. add your Cloud Vision API Key under a `visionAPIKey.txt` file
4. Modify the [table names](https://github.com/maxcraig112/GoGiffy/blob/main/src/bigquery.go#L14-L17) to correspond to your bigquery table
    - Your table schemas should correspond with the structs in [bigquery.go](https://github.com/maxcraig112/GoGiffy/blob/main/src/bigquery.go#L24-L39)
5. Pray

If you get everything setup, you can run the bot using `task run`

### Why use Golang?
Idk
### Why use Google Cloud Platform
What's with all the questions?
