require 'rubygems'
require 'eventmachine'
require 'em-websocket'
require 'thread'
require 'agent'
require 'json'
require '../tcatrb/twitter'

chan = Agent::Channel.new(name: "updates", type: Hash)
connections = []

config = Yajl::Parser.parse(File.read("../twitter/account.json"))
Twitter::Stream.new(config, chan)

go do 
  while entry = chan.receive
    # puts "EVENT" + entry.to_json
    connections.each {|c| c.send({
      "Text" => entry[:text], 
      "ImageURL" => entry[:user][:profile_image_url],
      "Username" => entry[:user][:screen_name]
      # "user" => entry[:user]
    }.to_json)}
  end
end

EventMachine::WebSocket.start(:host => "0.0.0.0", :port => 3100) do |ws|
  ws.onopen    { connections << ws; puts "New connection: #{connections.size} total" }
  ws.onclose   { connections.delete(ws); puts "Connection closed: #{connections.size} total" }
end
