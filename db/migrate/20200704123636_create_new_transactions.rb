class CreateNewTransactions < ActiveRecord::Migration
  def self.up
    create_table :transactions do |t|
      t.string :date
      t.string :category
      t.integer :amount

      t.timestamps
    end

    add_index :transactions, :category
  end

  def self.down
    drop_table :transactions
  end
end
