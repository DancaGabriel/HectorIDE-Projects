const API_BASE_URL = 'http://localhost:5000/api/recipes'; // Adjust if your backend runs on a different port

/**
 * Fetches all recipes from the API.
 * @returns {Promise<Array>} A promise that resolves to an array of recipe objects.
 */
async function getAllRecipes() {
  try {
    const response = await fetch(API_BASE_URL);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Error fetching recipes:', error);
    throw error; // Re-throw the error for handling in the component
  }
}

/**
 * Fetches a single recipe by its ID.
 * @param {number} id The ID of the recipe to fetch.
 * @returns {Promise<object>} A promise that resolves to a recipe object.
 */
async function getRecipeById(id) {
  try {
    const response = await fetch(`${API_BASE_URL}/${id}`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error(`Error fetching recipe with ID ${id}:`, error);
    throw error;
  }
}

/**
 * Creates a new recipe.
 * @param {object} recipe The recipe object to create.
 * @returns {Promise<object>} A promise that resolves to the created recipe object.
 */
async function createRecipe(recipe) {
  try {
    const response = await fetch(API_BASE_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(recipe),
    });
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Error creating recipe:', error);
    throw error;
  }
}

/**
 * Updates an existing recipe.
 * @param {number} id The ID of the recipe to update.
 * @param {object} recipe The updated recipe object.
 * @returns {Promise<object>} A promise that resolves to the updated recipe object.
 */
async function updateRecipe(id, recipe) {
  try {
    const response = await fetch(`${API_BASE_URL}/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(recipe),
    });
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error(`Error updating recipe with ID ${id}:`, error);
    throw error;
  }
}

/**
 * Deletes a recipe by its ID.
 * @param {number} id The ID of the recipe to delete.
 * @returns {Promise<void>} A promise that resolves when the recipe is deleted.
 */
async function deleteRecipe(id) {
  try {
    const response = await fetch(`${API_BASE_URL}/${id}`, {
      method: 'DELETE',
    });
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    // No content is expected on successful DELETE, but handle potential responses.
    if (response.status !== 204) { // 204 No Content is the standard for successful DELETE
        return await response.json(); // Handle cases where the server returns a JSON response on DELETE
    }
    return; // Resolve with void if the DELETE was successful and no content is expected.
  } catch (error) {
    console.error(`Error deleting recipe with ID ${id}:`, error);
    throw error;
  }
}

export { getAllRecipes, getRecipeById, createRecipe, updateRecipe, deleteRecipe };