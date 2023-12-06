# ChristmasTree

christmas tree cli

# dependency

System dependencies:

- Ruby 3.0+
- Ncurses:
  - Macos: `brew install ncurses`
  - Linux: `apt install libncurses5-dev` or `apt install libncursesw5-dev`

Gem dependencies

Ruby below 3.0, need install gems by manually:

- Curses: `gem install curses`

# local run

`ruby christmas_tree.rb`

`./christmas_tree.rb`

## add your name

`./christmas_tree.rb --merry_to <your name>`

# remote run

## use GEM

`gem exec christmas_tree`

### add your name

`gem exec christmas_tree --merry_to <yourname>`

## use Curl

`ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Mark24Code/christmas_tree/main/christmas_tree.rb)"`

### add your name

`ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Mark24Code/christmas_tree/main/christmas_tree.rb)" -- --merry_to <your name>`

# preview

![img](./demo.png)
