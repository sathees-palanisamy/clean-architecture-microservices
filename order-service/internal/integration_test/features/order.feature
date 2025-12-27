Feature: Order Processing
  As a customer
  I want to place an order for a product
  So that I can receive the item I want

  Scenario: Successful order creation
    Given a product exists with ID 1, name "Wireless Mouse", price 25.0, and stock 10
    When I create an order for product ID 1 with quantity 2 for user 101
    Then the order should be successfully created
    And the product stock should be 8

  Scenario: Insufficient stock for order
    Given a product exists with ID 2, name "Laptop", price 1200.0, and stock 1
    When I create an order for product ID 2 with quantity 5 for user 102
    Then the order creation should fail with "insufficient stock"

  Scenario: Product not found during order creation
    Given product ID 999 does not exist
    When I create an order for product ID 999 with quantity 1 for user 103
    Then the order creation should fail with "product not found"

  Scenario: Successfully retrieving an order by ID
    Given a product exists with ID 1, name "Wireless Mouse", price 25.0, and stock 10
    And an order exists with ID 500 for product ID 1 and user 101
    When I retrieve the order with ID 500
    Then I should receive the order details for ID 500

  Scenario: Cancelling a pending order
    Given a product exists with ID 1, name "Wireless Mouse", price 25.0, and stock 10
    And an order exists with ID 600 for product ID 1 and user 101
    When I cancel the order with ID 600
    Then the order status should be "CANCELLED"
