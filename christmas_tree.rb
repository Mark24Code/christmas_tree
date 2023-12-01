#!/usr/bin/env ruby
# christmas_tree.rb
#
# author: Mark24
# mail: mark.zhangyoung@gmail.com
# github: https://github.com/Mark24Code
# version: 1.0
#
#
# LICENSE: MIT
#
# Copyright (c) 2023-forever Mark24Code

# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

module FontStylePatch
  RESET_COLOR = "\e[0m" #重置所有颜色和样式
  COLORS = {
    # https://ss64.com/nt/syntax-ansi.html
    black:"\e[30m",
    white:"\e[97m",
    red:"\e[31m",
    green:"\e[32m",
    yellow:"\e[33m",
    blue:"\e[34m",
    magenta:"\e[35m",
    cyan:"\e[36m",
    gray:"\e[90m",
    light_gray:"\e[37m",
    light_red:"\e[91m",
    light_green:"\e[92m",
    light_yellow:"\e[93m",
    light_blue:"\e[94m",
    light_magenta:"\e[95m",
    light_cyan:"\e[96m",

    bold:"\e[1m",
    underline:"\e[4m",
    nounderline:"\e[24m",
    reversetext:"\e[7m",
  }
  COLORS.keys.each do |color_name|
    define_method(color_name) do
      return "#{COLORS[color_name]}#{self}#{RESET_COLOR}"
    end
  end
end


class String
  include FontStylePatch
end

class ChristmasTree
  CHAR_SPACE = ' '
  def initialize(layer_count = 10)
    @layer_count = layer_count
    @scale = 3
    @layer = []
    @buffer = []
    @light_char = 'o'
    @light_colors = [
      :white,
      :red,
      :green,
      :yellow,
      :blue,
      :magenta,
      :cyan,
      :gray,
      :light_gray,
      :light_red,
      :light_green,
      :light_yellow,
      :light_blue,
      :light_magenta,
      :light_cyan,
    ]
    @speed = 3 # 3 times per second
  end

  def tree_crown
    (1..@layer_count).each do |count|
      total = 2*count-1
      index_arr = []
      cache = []

      total.times.each do |i|
        index_arr.push(i)
        cache << "*".green
      end

      random_count = (total * 0.3).to_i
      random_count.times {
        choose_index = index_arr.sample
        cache[choose_index] = @light_char.__send__(@light_colors.sample)
      }
      @layer << cache.join("")
    end
  end

  def tree_trunk(height=2)
    height.times {
      @layer << "mWm".white
    }
  end

  def footer_text
    @layer << "MARRY CHRISTMAS".yellow
    code_text = "CODE".split("").map {|c| c.bold.__send__(@light_colors.sample)}.join("")
    text = "And lots of ".yellow + code_text + " in 2024".yellow
    @layer << text
  end

  def border(height=4)
    height.times {
      @layer << CHAR_SPACE
    }
  end

  def canvas
    border(3)
    tree_crown
    tree_trunk
    footer_text
    border
    @buffer = @layer.clone
  end

  def line_process(line)
    half_space = CHAR_SPACE * ((@width - real_width(line)) / 2)
    "#{half_space}#{line}#{half_space}"
  end

  def real_width(text)
    text2 = text.clone
    text2 = text2.gsub(/\e\[(\d+)m/,'')
    return text2.length
  end

  def draw_buffer
    @width = (@buffer.map {|l| real_width(l)}).max * @scale

    screen = nil
    screen_buffer = "\r" + @buffer.map {|line| line_process(line)}.join("\n")
    printf screen_buffer
    STDOUT.flush
  end

  def clean_screen
    # system("clear")
    print("\033[H\033[J")
  end

  def draw
    while true
      clean_screen
      canvas
      draw_buffer
      sleep 1.0/@speed
    end
  end
end

ChristmasTree.new.draw
