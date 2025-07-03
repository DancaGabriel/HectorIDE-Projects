require 'rack/test'
require 'rspec'
require 'json'
require_relative '../backend/app'  # Adjust path as needed

ENV['RECIPE_BOOK_DATABASE_URI'] = 'sqlite:///:memory:' # Use an in-memory database for testing

module RSpecMixin
  include Rack::Test::Methods
  def app() Sinatra::Application end
end

# For RSpec 2.x and 3.x
RSpec.configure do |config|
  config.include RSpecMixin
  config.before(:each) do
    # Create a new app for each test to ensure a clean state
    @app = create_app('config.Config') # Assuming config.py exists and is needed
    @app.config['TESTING'] = true
    @db = @app.extensions['sqlalchemy'].db # Access the db instance
    @db.create_all()  # Create tables
  end

  config.after(:each) do
      @db.session.remove() # Close the connection
      @db.drop_all() # drop the database after each test
  end

  # Helper method to parse JSON responses
  def json_response
    JSON.parse(last_response.body)
  end

  # Helper method to create a recipe.
  def create_recipe(title, ingredients, instructions)
    post '/', {
      title: title,
      ingredients: ingredients,
      instructions: instructions
    }.to_json, { 'CONTENT_TYPE' => 'application/json' }
  end

  def get_recipe(recipe_id)
    get "/#{recipe_id}"
  end

  def update_recipe(recipe_id, title, ingredients, instructions)
    put "/#{recipe_id}", {
      title: title,
      ingredients: ingredients,
      instructions: instructions
    }.to_json, { 'CONTENT_TYPE' => 'application/json' }
  end

  def delete_recipe(recipe_id)
    delete "/#{recipe_id}"
  end
end