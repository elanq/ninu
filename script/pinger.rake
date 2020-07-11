desc 'Pings'
task :ping do
  require 'net/http'

  if ENV['PING_URL']
    uri = URI(ENV['PING_URL'])
    response = Net::HTTP.get_response(uri)
    puts "PING with status code #{response.code}"
  end
end

