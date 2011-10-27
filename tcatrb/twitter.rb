require 'rubygems'
require 'yajl/http_stream'
require 'thread'
require 'agent'

module Twitter
  class Stream
    def initialize config, chan
      @config = config
      @chan = chan
      go{self.process}
    end
    
    def process
      # uri = URI.parse("https://#{@config['username']}:#{@config['password']}@stream.twitter.com/1/statuses/sample.json")
      uri = URI.parse("http://localhost:8001/")
      Yajl::HttpStream.get(uri, :symbolize_keys => true) do |hash|
        @chan << hash unless hash == nil || hash[:user] == nil
      end
    end
    
  end
end