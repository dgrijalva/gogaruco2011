require './twitter'
require 'thread'
require 'agent'

chan = Agent::Channel.new(name: "updates", type: Hash)

config = Yajl::Parser.parse(File.read("../twitter/account.json"))
Twitter::Stream.new(config, chan)

while update = chan.receive do
  puts "#{update[:user][:screen_name]}: #{update[:text].gsub("\n", '')}"
end