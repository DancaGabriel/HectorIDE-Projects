require 'rack/test'
require 'rspec'
require 'json'

# Assuming the Flask app is served by a WSGI server in app.py
require_relative '../../backend/app'

describe 'Recipe Routes API' do
  include Rack::Test::Methods

  def app
    # Create a test app instance using the create_app function.
    # Use an in-memory SQLite database for testing.
    test_config = {
      'SQLALCHEMY_DATABASE_URI' => 'sqlite:///:memory:',
      'SQLALCHEMY_TRACK_MODIFICATIONS' => false,
      'TESTING' => true
    }
    app = create_app
    app.config.update(test_config)

    # Create the database tables within the test application context.
    require_relative '../../backend/app/models/recipe'
    require_relative '../../backend/app/__init__'
    app.app_context.push
    db = app.extensions['sqlalchemy'].db
    db.create_all
    return app
  end

  before(:each) do
    # Clear the database before each test
    require_relative '../../backend/app/models/recipe'
    require_relative '../../backend/app/__init__'
    db = app.extensions['sqlalchemy'].db
    Recipe.delete_all # Assuming Recipe is accessible here.  Correct way to empty table for testing.
    db.session.commit

    # Example of seeding data - optional, but helpful for some tests
    @recipe1_data = { 'title' => 'Pasta Carbonara', 'ingredients' => "Pasta\nEggs\nBacon\nCheese", 'instructions' => "Cook pasta\nMix eggs and cheese\nFry bacon\nCombine all" }
    @recipe2_data = { 'title' => 'Chicken Stir-fry', 'ingredients' => "Chicken\nVegetables\nSoy sauce", 'instructions' => "Cook chicken\nStir-fry vegetables\nAdd soy sauce" }
  end

  after(:each) do
    require_relative '../../backend/app/models/recipe'
    require_relative '../../backend/app/__init__'
    db = app.extensions['sqlalchemy'].db
    db.session.remove
    db.drop_all # Drops tables at the end of the test suite.
  end

  describe 'GET /' do
    it 'returns an empty list when no recipes exist' do
      get '/'
      expect(last_response.status).to eq(200)
      expect(JSON.parse(last_response.body)).to eq([])
    end

    it 'returns a list of recipes when they exist' do
      post '/', @recipe1_data.to_json, { 'CONTENT_TYPE' => 'application/json' }
      post '/', @recipe2_data.to_json, { 'CONTENT_TYPE' => 'application/json' }

      get '/'
      expect(last_response.status).to eq(200)
      recipes = JSON.parse(last_response.body)
      expect(recipes.size).to eq(2)
      expect(recipes[0]['title']).to eq('Pasta Carbonara')
      expect(recipes[1]['title']).to eq('Chicken Stir-fry')
    end
  end

  describe 'GET /<id>' do
    it 'returns a recipe when it exists' do
      post '/', @recipe1_data.to_json, { 'CONTENT_TYPE' => 'application/json' }
      recipe = JSON.parse(last_response.body)
      recipe_id = recipe['id']

      get "/#{recipe_id}"
      expect(last_response.status).to eq(200)
      retrieved_recipe = JSON.parse(last_response.body)
      expect(retrieved_recipe['title']).to eq('Pasta Carbonara')
    end

    it 'returns a 404 when the recipe does not exist' do
      get '/999'
      expect(last_response.status).to eq(404)
      expect(JSON.parse(last_response.body)['message']).to eq('Recipe not found')
    end
  end

  describe 'POST /' do
    it 'creates a new recipe' do
      post '/', @recipe1_data.to_json, { 'CONTENT_TYPE' => 'application/json' }
      expect(last_response.status).to eq(201)
      recipe = JSON.parse(last_response.body)
      expect(recipe['title']).to eq('Pasta Carbonara')

      get '/'
      recipes = JSON.parse(last_response.body)
      expect(recipes.size).to eq(1)
    end

    it 'returns a 400 if the recipe data is invalid' do
      post '/', { 'title' => nil }.to_json, { 'CONTENT_TYPE' => 'application/json' } # Missing required fields
      expect(last_response.status).to eq(400)
    end
  end

  describe 'PUT /<id>' do
    it 'updates an existing recipe' do
      post '/', @recipe1_data.to_json, { 'CONTENT_TYPE' => 'application/json' }
      recipe = JSON.parse(last_response.body)
      recipe_id = recipe['id']

      put "/#{recipe_id}", { 'title' => 'Updated Pasta Carbonara' }.to_json, { 'CONTENT_TYPE' => 'application/json' }
      expect(last_response.status).to eq(200)
      updated_recipe = JSON.parse(last_response.body)
      expect(updated_recipe['title']).to eq('Updated Pasta Carbonara')

      get "/#{recipe_id}"
      retrieved_recipe = JSON.parse(last_response.body)
      expect(retrieved_recipe['title']).to eq('Updated Pasta Carbonara')
    end

    it 'returns a 404 if the recipe does not exist' do
      put '/999', { 'title' => 'Updated Pasta Carbonara' }.to_json, { 'CONTENT_TYPE' => 'application/json' }
      expect(last_response.status).to eq(404)
      expect(JSON.parse(last_response.body)['message']).to eq('Recipe not found')
    end

    it 'returns a 400 if the updated recipe data is invalid' do
       post '/', @recipe1_data.to_json, { 'CONTENT_TYPE' => 'application/json' }
       recipe = JSON.parse(last_response.body)
       recipe_id = recipe['id']
       put "/#{recipe_id}", { 'title' => nil }.to_json, { 'CONTENT_TYPE' => 'application/json' }
       expect(last_response.status).to eq(400)
    end
  end

  describe 'DELETE /<id>' do
    it 'deletes a recipe' do
      post '/', @recipe1_data.to_json, { 'CONTENT_TYPE' => 'application/json' }
      recipe = JSON.parse(last_response.body)
      recipe_id = recipe['id']

      delete "/#{recipe_id}"
      expect(last_response.status).to eq(200)
      expect(JSON.parse(last_response.body)['message']).to eq('Recipe deleted')

      get "/#{recipe_id}"
      expect(last_response.status).to eq(404) # Confirm deleted
    end

    it 'returns a 404 if the recipe does not exist' do
      delete '/999'
      expect(last_response.status).to eq(404)
      expect(JSON.parse(last_response.body)['message']).to eq('Recipe not found')
    end
  end
end