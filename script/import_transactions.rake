require 'pg'
require 'pry-byebug'
require 'csv'

def connect_pg!
  if ENV['DATABASE_URL']
    uri = ENV['DATABASE_URL']
    db = URI.parse(uri)
    args = {
      host: db.host,
      port: db.port,
      dbname: db.path[1..-1],
      user: db.user,
      password: db.password
    }
  else
    args = {
      dbname: 'ninu',
      host: ENV['PG_HOST'],
      port: 5432,
      user: ENV['PG_USERNAME'],
      password: ENV['PG_PASSWORD']
    }
  end

  PG.connect(args)
end

def help_message
  "please specify filename argument. Example rake import:transactions ./transactions/june2020-july2020-transactions.csv"
end

namespace :import do
  task transactions: :environment do
    fileloc = ARGV[1]
    unless fileloc.present?
      puts help_message
      exit
    end

    pg_client = connect_pg!

    pg_client.prepare('stmt', "INSERT INTO transactions (date, category, amount) VALUES ($1, $2, $3)" )

    CSV.table(fileloc, {col_sep: ';', }).each do |row|
      amount = row.fetch(:pengeluaran).tr('^0-9', '').to_i
      date = Date.parse(row.fetch(:tanggal))
      values = [
        date.to_s,
        row.fetch(:detail),
        amount
      ]
      pg_client.exec_prepared('stmt', values)
      puts "#{row.fetch(:detail)} inserted"
    rescue StandardError => e
      puts e.message
    end

  end
end
