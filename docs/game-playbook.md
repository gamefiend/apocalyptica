# Games and Playbook Format

You can add your own *Powered By the Apocalypse* game to Apocalyptica's list. 

Apocalyptica reads its list of games out of `data/game`.

Each game has its own directory. containing files:

- `game.json`: contains game metadata like Name, Author, and a list of playbooks.
- `basic.json`: Basic moves for the game are always in this file.
- `<playbook>.json`: 1 or more files, each with moves specific to playbooks for the game.

a skeleton structure that can be copied and applied is in `data/example/game/example_game`. Copy that directory into `data/game`, rename and modify as you need.

# Contributing a Game

We would love to support many PbtA games out of the box! More importantly, though, we want to honor the wishes of creators. Please file an issue with "[GAME] My Game" in the title **before** submitting a PR and tell us:

- info about the game (Name, publisher, short description)
- if there is permission from the creator of the game to use it, with some proof.

Thanks in advance!

