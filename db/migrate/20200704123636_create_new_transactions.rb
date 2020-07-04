class CreateNewTransactions < ActiveRecord::Migration[5.2]
  def self.up
    create_table :transactions do |t|
      t.date :date
      t.string :category
      t.integer :amount

      t.timestamps
    end

    add_index :transactions, :date
    add_index :transactions, :category
  end

  def self.down
    drop_table :transactions
  end
end
