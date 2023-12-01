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

require 'curses'
Curses.init_screen
Curses.noecho
Curses.curs_set(0)
COLORS_ORDER = [
  :black, :red, :green, :yellow, :blue, :magenta, :cyan, :white,
  :bright_black, :bright_red, :bright_green, :bright_yellow,
  :bright_blue, :bright_magenta, :bright_cyan, :bright_white
]

COLORS_MAP = {}
Curses.start_color
COLORS_ORDER.length.times { |i|
  Curses.init_pair(i, i, 0)
  COLORS_MAP[COLORS_ORDER[i]] = Curses.color_pair(i)
}

FONT_WEIGHT_MAP = {
  normal: Curses::A_NORMAL,
  bold: Curses::A_BOLD,
}

class Token
  attr_accessor :content, :font_weight, :color_pair
  def initialize(content, color_pair = nil, font_weight = nil)
    @content = content
    @font_weight = font_weight || :normal
    @color_pair = color_pair || :white
  end

  def length
    @content.length
  end

  def draw(&block)
    Curses.attron(COLORS_MAP[@color_pair] | FONT_WEIGHT_MAP[@font_weight] ){
      yield @content
    }
  end
end

class ChristmasTree
  CHAR_SPACE = ' '
  def initialize(layer_count = 10)
    @buffer_count = layer_count
    @scale = 3
    @buffer = []
    @light_char = 'o'
    @snow_char = '*'
    @light_colors = COLORS_ORDER
    @speed = 3 # 3 times per second
    @width = nil
  end

  def tree_crown
    (1..@buffer_count).each do |count|
      total = 2*count-1
      index_arr = []
      line = []

      total.times.each do |i|
        index_arr.push(i)
        line << Token.new("*", :green)
      end



      # snow
      snow_random_count = (total * 0.3).to_i
      snow_random_count.times {
        choose_index = index_arr.sample
        line[choose_index] = Token.new(@snow_char, :white)
      }

      # light
      light_random_count = (total * 0.3).to_i
      light_random_count.times {
        choose_index = index_arr.sample
        line[choose_index] = Token.new(@light_char, @light_colors.sample)
      }
      @buffer << line
    end
  end

  def tree_trunk(height=2)
    height.times {
      @buffer << [Token.new("mWm", :white)]
    }
  end

  def footer_text
    @buffer << [Token.new("MARRY CHRISTMAS", :yellow)]
    line = [Token.new("And lots of ", :yellow)]
    code_text = "CODE".split("").map {|c| Token.new(c, @light_colors.sample, :bold)}
    code_text.each do |code|
      line << code
    end
    line << Token.new(" in #{Time.now.year + 1}", :yellow)
    @buffer << line
  end

  def border(height=4)
    height.times {
      @buffer << [Token.new(CHAR_SPACE)]
    }
  end

  def draw_canvas
    @buffer = []
    border(3)
    tree_crown
    tree_trunk
    footer_text
    border(3)

    if !@width
      @width = (@buffer.map {|l| real_width(l)}).max * @scale
    end

    @buffer.each_with_index do |buffer_line,y|
      x = (@width - real_width(buffer_line)) / 2
      Curses.setpos(y, x)
      buffer_line.each do |token|
        token.draw do |text|
          Curses.addstr(text)
        end
      end
    end
  end

  def line_process(line)
    half_space = CHAR_SPACE * ((@width - real_width(line)) / 2)
    "#{half_space}#{line}#{half_space}"
  end

  def real_width(line)
    line.map {|token| token.length }.sum
  end


  def clean_screen
    Curses.clear
  end

  def draw
    clean_screen
    while true
      draw_canvas
      Curses.refresh
      sleep 1.0/@speed
    end
  end
end

begin
  ChristmasTree.new.draw
ensure
  Curses.close_screen
end
