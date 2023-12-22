# ChristmasTree

christmas tree cli

# dependency

System dependencies:

- Ruby 3.0+
  - MacOS: `brew install ruby`
  - Linux: `apt install ruby`
- Ncurses:
  - MacOS: `brew install ncurses`
  - Linux: `apt install libncurses5-dev` or `apt install libncursesw5-dev`

Gem dependencies:
  below Ruby 3.0  need install gems manually:
  
  - `gem install curses`

# Execution

## Way 1: use GEM remote run


`gem exec christmas_tree`


if you meet path trouble, use absolute brew ruby path:

`/usr/local/opt/ruby/bin/gem exec christmas_tree`


## Way 2: download it and execute on local machine

Enter project directory: `cd <project_path>`

Execute:

`./christmas_tree.rb`


# One more thing

Add your name

``./christmas_tree.rb --merry_to <your name>``

# preview

![img](./demo.png)
