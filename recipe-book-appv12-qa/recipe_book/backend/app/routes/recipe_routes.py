from flask import Blueprint, jsonify, request
from ..services import recipe_service
from ..exceptions import RecipeNotFoundException, InvalidRecipeDataException

recipe_routes = Blueprint('recipes', __name__)

@recipe_routes.route('/', methods=['GET'])
def get_all_recipes():
    """
    Retrieve all recipes.
    """
    recipes = recipe_service.get_all_recipes()
    return jsonify([recipe.to_dict() for recipe in recipes]), 200

@recipe_routes.route('/<int:recipe_id>', methods=['GET'])
def get_recipe(recipe_id):
    """
    Retrieve a specific recipe by ID.
    """
    try:
        recipe = recipe_service.get_recipe_by_id(recipe_id)
        return jsonify(recipe.to_dict()), 200
    except RecipeNotFoundException:
        return jsonify({'message': 'Recipe not found'}), 404

@recipe_routes.route('/', methods=['POST'])
def create_recipe():
    """
    Create a new recipe.
    """
    data = request.get_json()
    try:
        recipe = recipe_service.create_recipe(data)
        return jsonify(recipe.to_dict()), 201
    except InvalidRecipeDataException as e:
        return jsonify({'message': str(e)}), 400
    except Exception as e:
        return jsonify({'message': 'Failed to create recipe', 'error': str(e)}), 500


@recipe_routes.route('/<int:recipe_id>', methods=['PUT'])
def update_recipe(recipe_id):
    """
    Update an existing recipe.
    """
    data = request.get_json()
    try:
        recipe = recipe_service.update_recipe(recipe_id, data)
        return jsonify(recipe.to_dict()), 200
    except RecipeNotFoundException:
        return jsonify({'message': 'Recipe not found'}), 404
    except InvalidRecipeDataException as e:
        return jsonify({'message': str(e)}), 400
    except Exception as e:
        return jsonify({'message': 'Failed to update recipe', 'error': str(e)}), 500


@recipe_routes.route('/<int:recipe_id>', methods=['DELETE'])
def delete_recipe(recipe_id):
    """
    Delete a recipe.
    """
    try:
        recipe_service.delete_recipe(recipe_id)
        return jsonify({'message': 'Recipe deleted'}), 200
    except RecipeNotFoundException:
        return jsonify({'message': 'Recipe not found'}), 404
    except Exception as e:
        return jsonify({'message': 'Failed to delete recipe', 'error': str(e)}), 500