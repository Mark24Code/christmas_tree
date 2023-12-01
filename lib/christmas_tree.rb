# frozen_string_literal: true

require_relative "christmas_tree/version"

require 'optparse'
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


module ChristmasTree
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

  class Error < StandardError; end

  class Tree
    CHAR_SPACE = ' '
    def initialize(name = nil, layer_count = 10)
      @name = name
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


    def color_tokens(text)
      text.split("").map {|c| Token.new(c, @light_colors.sample, :bold)}
    end

    def footer_text
      @buffer << [Token.new("MARRY CHRISTMAS", :yellow)]

      # name
      if @name
        name = "@"+@name.strip
        line = []
        color_tokens(name).each do |code|
          line << code
        end

        @buffer << line
      end

      # code
      line = [Token.new("And lots of ", :yellow)]
      color_tokens("CODE").each do |code|
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
end
