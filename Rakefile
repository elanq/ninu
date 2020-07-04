require 'standalone_migrations'
require 'tasks/standalone_migrations'

StandaloneMigrations::Configurator.environments_config do |env|
  env.on 'production' do
    uri = ENV['DATABASE_URL']
    if uri
      db = URI.parse(uri)
      return {
        adapter: db.scheme == 'postgres' ? 'postgresql' : db.scheme,
        host: db.host,
        username: db.username,
        password: db.password,
        database: db.path[1..-1],
        encoding: 'utf8'
      }
    end

    nil
  end
end
StandaloneMigrations::Tasks.load_tasks
