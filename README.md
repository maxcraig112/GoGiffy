# GoGiffy
The original Giffy Bot however written in Golang and using the Google Cloud Platform.

## Giffy-Bot
A multipurpose repos designed to allow the ~~manipulation~~, tagging, archiving and retrieval of gifs.

This repos provides abstraction of these commands through the discord bot Giffy.

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