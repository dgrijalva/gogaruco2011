require 'thread'
require 'rubygems'
require 'agent'

def produce chan, limit 
  a, b = 0, 1
  chan << 0
  while b < limit
    chan << b
    print "Producer sent #{b}\n"
    a, b = b, a+b
  end
  chan.close
end

def consume chan, done, n
  until chan.closed?
    sleep(rand(100) * 0.001)
    i = chan.receive
    print "Consumer #{n} received #{i}\n"
  end
rescue
ensure
  done << 1
end

chan = Agent::Channel.new(name: 'fib', type: Integer)
done = Agent::Channel.new(name: 'complete', type: Integer)
go{produce chan, 1000}
10.times{|i| go{consume(chan, done, i)}}
10.times{done.receive}