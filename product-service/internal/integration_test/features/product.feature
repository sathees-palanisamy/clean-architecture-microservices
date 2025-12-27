Feature: Product Management
  As a store administrator
  I want to manage products and their stock levels
  So that I can ensure inventory accuracy

  Scenario: Successful product creation
    When I create a product with SKU "PROD-001", name "Premium Headphones", and price 99.50
    Then the product should be successfully saved

  Scenario: Reserving stock for a product
    Given a product exists with SKU "STOCK-01", name "Keyboard", price 50.0, and stock 10
    When I reserve 4 units of stock for this product
    Then the reservation should be successful
    And the available stock should be 6

  Scenario: Failing to reserve excessive stock
    Given a product exists with SKU "STOCK-02", name "Monitor", price 300.0, and stock 2
    When I reserve 5 units of stock for this product
    Then the reservation should fail with "insufficient stock"
