import React, { useState, useEffect } from 'react';
import RecipeList from './components/RecipeList';
import api from './services/api';
import './styles.css';

function App() {
  const [recipes, setRecipes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchRecipes = async () => {
      try {
        const data = await api.get('/recipes');
        setRecipes(data);
        setLoading(false);
      } catch (error) {
        console.error("Error fetching recipes:", error);
        setError("Failed to load recipes.");
        setLoading(false);
      }
    };

    fetchRecipes();
  }, []);

  const handleRecipeCreated = async (newRecipe) => {
    try {
      const createdRecipe = await api.post('/recipes', newRecipe);
      setRecipes([...recipes, createdRecipe]);
    } catch (error) {
      console.error("Error creating recipe:", error);
      setError("Failed to create recipe.");
    }
  };

    const handleRecipeUpdated = async (updatedRecipe) => {
        try {
            await api.put(`/recipes/${updatedRecipe.id}`, updatedRecipe);
            const updatedRecipes = recipes.map(recipe =>
                recipe.id === updatedRecipe.id ? updatedRecipe : recipe
            );
            setRecipes(updatedRecipes);
        } catch (error) {
            console.error("Error updating recipe:", error);
            setError("Failed to update recipe.");
        }
    };

  const handleRecipeDeleted = async (recipeId) => {
    try {
      await api.delete(`/recipes/${recipeId}`);
      const updatedRecipes = recipes.filter(recipe => recipe.id !== recipeId);
      setRecipes(updatedRecipes);
    } catch (error) {
      console.error("Error deleting recipe:", error);
      setError("Failed to delete recipe.");
    }
  };


  if (loading) {
    return <div>Loading recipes...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  return (
    <div className="App">
      <h1>Recipe Book</h1>
      <RecipeList 
          recipes={recipes} 
          onRecipeCreated={handleRecipeCreated}
          onRecipeUpdated={handleRecipeUpdated}
          onRecipeDeleted={handleRecipeDeleted}
      />
    </div>
  );
}

export default App;