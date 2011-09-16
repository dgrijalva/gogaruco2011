require 'fiber'

limit = 1000
produce = Fiber.new do 
  a, b = 0, 1
  Fiber.yield 0
  while b < limit
    Fiber.yield b
    print "Producer sent #{b}\n"
    a, b = b, a+b
  end
end

def consume chan
  while true
    i = chan.resume
    print "Consumer received #{i}\n"
  end
rescue
  # Resuming a dead fiber will raise an exception.
  # Just catching the exception is the lazy way out.
end

consume produce