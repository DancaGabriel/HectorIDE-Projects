(function() {
    'use strict';

    document.addEventListener('DOMContentLoaded', () => {
        // --- API Configuration ---
        const API_BASE_URL = 'http://localhost:8000'; // Adjust if your backend runs on a different port
        const AUTH_TOKEN_KEY = 'taskflow_auth_token';

        // --- Helper Functions ---
        const setAuthToken = (token) => {
            localStorage.setItem(AUTH_TOKEN_KEY, token);
        };

        const getAuthToken = () => {
            return localStorage.getItem(AUTH_TOKEN_KEY);
        };

        const clearAuthToken = () => {
            localStorage.removeItem(AUTH_TOKEN_KEY);
        };

        const displayMessage = (message, type = 'info') => {
            const messageContainer = document.createElement('div');
            messageContainer.classList.add('message', type);
            messageContainer.textContent = message;

            // Append to the top of the main container
            const mainContainer = document.querySelector('.app-main');
            mainContainer.insertBefore(messageContainer, mainContainer.firstChild);

            // Automatically remove the message after a few seconds
            setTimeout(() => {
                messageContainer.remove();
            }, 5000); // Remove after 5 seconds
        };

        const handleApiError = (error) => {
            console.error('API Error:', error);
            displayMessage(`An error occurred: ${error.message || error}`, 'error');
        };

        const validateForm = (formId) => {
            const form = document.getElementById(formId);
            let isValid = true;

            // Reset error messages
            form.querySelectorAll('.error-message').forEach(el => el.textContent = '');

            form.querySelectorAll('[required]').forEach(input => {
                if (!input.value.trim()) {
                    const errorId = input.id + '-error';
                    const errorEl = document.getElementById(errorId);
                    if (errorEl) {
                        errorEl.textContent = 'This field is required.';
                    }
                    isValid = false;
                }

                // Email Validation
                if (input.type === 'email' && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(input.value)) {
                    const errorId = input.id + '-error';
                    const errorEl = document.getElementById(errorId);
                    if (errorEl) {
                        errorEl.textContent = 'Please enter a valid email address.';
                    }
                    isValid = false;
                }
            });

            return isValid;
        };

        // --- User Registration ---
        const registerUser = async (formData) => {
            try {
                const response = await fetch(`${API_BASE_URL}/users/register`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(Object.fromEntries(formData)),
                });

                if (!response.ok) {
                    const errorData = await response.json();
                    throw new Error(errorData.error || `Registration failed with status ${response.status}`);
                }

                const userData = await response.json();
                displayMessage(`Registration successful! Welcome, ${userData.username}!`, 'success');

                // Redirect to login or automatically log in
                window.location.hash = '#login';

            } catch (error) {
                handleApiError(error);
            }
        };

        const registrationForm = document.getElementById('registration-form');
        registrationForm.addEventListener('submit', (event) => {
            event.preventDefault();
            if (validateForm('registration-form')) {
                const formData = new FormData(registrationForm);
                registerUser(formData);
            }
        });

        // --- User Login ---
        const loginUser = async (formData) => {
            try {
                const response = await fetch(`${API_BASE_URL}/users/login`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(Object.fromEntries(formData)),
                });

                if (!response.ok) {
                    const errorData = await response.json();
                    throw new Error(errorData.error || `Login failed with status ${response.status}`);
                }

                const data = await response.json();
                setAuthToken(data.token);
                displayMessage('Login successful!', 'success');

                // Redirect to the posts section
                window.location.hash = '#posts'; // Or any other protected area
            } catch (error) {
                handleApiError(error);
            }
        };

        const loginForm = document.getElementById('login-form');
        loginForm.addEventListener('submit', (event) => {
            event.preventDefault();
            if (validateForm('login-form')) {
                const formData = new FormData(loginForm);
                loginUser(formData);
            }
        });

        // --- Post Creation ---
        const createPost = async (formData) => {
            try {
                const token = getAuthToken();
                if (!token) {
                    throw new Error('Not authenticated.');
                }

                const response = await fetch(`${API_BASE_URL}/posts`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`, // Include JWT token
                    },
                    body: JSON.stringify({
                        title: formData.get('post_title'),
                        content: formData.get('post_content'),
                        category_ids: Array.from(document.getElementById('post_categories').selectedOptions).map(option => parseInt(option.value)),
                        // author_id might come from the JWT payload on the server-side instead
                    }),
                });

                if (!response.ok) {
                    const errorData = await response.json();
                    throw new Error(errorData.error || `Post creation failed with status ${response.status}`);
                }

                const postData = await response.json();
                displayMessage('Post created successfully!', 'success');

                // Clear the form
                document.getElementById('create-post-form').reset();

                // Refresh the post list (or redirect to the new post)
                fetchPosts();

            } catch (error) {
                handleApiError(error);
            }
        };

        const createPostForm = document.getElementById('create-post-form');
        createPostForm.addEventListener('submit', (event) => {
            event.preventDefault();
            if (validateForm('create-post-form')) {
                const formData = new FormData(createPostForm);
                createPost(formData);
            }
        });

        // --- Fetch Posts ---
        const fetchPosts = async (page = 1, limit = 10) => {
            try {
                const token = getAuthToken();
                if (!token) {
                    throw new Error('Not authenticated.');
                }

                const response = await fetch(`${API_BASE_URL}/posts?page=${page}&limit=${limit}`, {
                    method: 'GET',
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'application/json', // Added to be explicit
                    }
                });

                if (!response.ok) {
                    const errorData = await response.json();
                    throw new Error(errorData.error || `Failed to fetch posts: ${response.status}`);
                }

                const posts = await response.json();
                renderPosts(posts);

            } catch (error) {
                handleApiError(error);
            }
        };

        const renderPosts = (posts) => {
          const postListBody = document.querySelector('#post-list tbody');
          postListBody.innerHTML = ''; // Clear existing posts

          if(posts.length === 0) {
            postListBody.innerHTML = '<tr><td colspan="5" class="text-center">No posts found.</td></tr>';
            return;
          }

          posts.forEach(post => {
              const row = document.createElement('tr');
              row.innerHTML = `
                  <td><a href="#post/${post.id}">${post.title}</a></td>
                  <td>Unknown Author</td>
                  <td>Unknown Date</td>
                  <td>Unknown Categories</td>
                  <td>
                      <button class="btn-secondary btn-small edit-post-button" data-id="${post.id}">Edit</button>
                      <button class="btn-action btn-small delete-post-button" data-id="${post.id}">Delete</button>
                  </td>
              `;
              postListBody.appendChild(row);
          });
        };

        // Initial fetch of posts
        fetchPosts();

        // --- Category Modal Functionality ---
        const createCategoryButton = document.getElementById('create-category-button');
        const categoryModal = document.getElementById('category-modal');
        const closeCategoryModalButton = document.getElementById('close-category-modal');
        const categoryForm = document.getElementById('category-form');
        const saveCategoryButton = document.getElementById('save-category-button');
        const deleteCategoryButton = document.getElementById('delete-category-button');
        const categoryNameInput = document.getElementById('category_name');
        const categoryIdInput = document.getElementById('category_id');

        createCategoryButton.addEventListener('click', () => {
            categoryIdInput.value = ''; // Clear ID for create
            categoryNameInput.value = ''; // Clear name field
            categoryModal.style.display = 'flex';
        });

        closeCategoryModalButton.addEventListener('click', () => {
            categoryModal.style.display = 'none';
        });

        window.addEventListener('click', (event) => {
            if (event.target === categoryModal) {
                categoryModal.style.display = 'none';
            }
        });

        const saveCategory = async (categoryData) => {
            try {
                const token = getAuthToken();
                if (!token) {
                    throw new Error('Not authenticated.');
                }

                const url = categoryData.id ? `${API_BASE_URL}/categories/${categoryData.id}` : `${API_BASE_URL}/categories`;
                const method = categoryData.id ? 'PUT' : 'POST';

                const response = await fetch(url, {
                    method: method,
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`,
                    },
                    body: JSON.stringify(categoryData),
                });

                if (!response.ok) {
                    const errorData = await response.json();
                    throw new Error(errorData.error || `Failed to save category: ${response.status}`);
                }

                displayMessage('Category saved successfully!', 'success');
                fetchCategories(); // Refresh category list
                categoryModal.style.display = 'none'; // Close the modal

            } catch (error) {
                handleApiError(error);
            }
        };

        categoryForm.addEventListener('submit', (event) => {
            event.preventDefault();
            if (validateForm('category-form')) {
                const categoryData = {
                    id: categoryIdInput.value || null,
                    name: categoryNameInput.value
                };
                saveCategory(categoryData);
            }
        });

        const deleteCategory = async (categoryId) => {
            try {
                const token = getAuthToken();
                if (!token) {
                    throw new Error('Not authenticated.');
                }

                const response = await fetch(`${API_BASE_URL}/categories/${categoryId}`, {
                    method: 'DELETE',
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    },
                });

                if (!response.ok) {
                    const errorData = await response.json();
                    throw new Error(errorData.error || `Failed to delete category: ${response.status}`);
                }

                displayMessage('Category deleted successfully!', 'success');
                fetchCategories(); // Refresh category list
                categoryModal.style.display = 'none';

            } catch (error) {
                handleApiError(error);
            }
        };

        deleteCategoryButton.addEventListener('click', () => {
            const categoryId = categoryIdInput.value;
            if (confirm('Are you sure you want to delete this category?')) {
                deleteCategory(categoryId);
            }
        });

        // --- Fetch Categories ---
        const fetchCategories = async () => {
            try {
                const token = getAuthToken();
                if (!token) {
                    throw new Error('Not authenticated.');
                }

                const response = await fetch(`${API_BASE_URL}/categories`, {
                    method: 'GET',
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    }
                });

                if (!response.ok) {
                    const errorData = await response.json();
                    throw new Error(errorData.error || `Failed to fetch categories: ${response.status}`);
                }

                const categories = await response.json();
                renderCategories(categories);

            } catch (error) {
                handleApiError(error);
            }
        };

        const renderCategories = (categories) => {
          const categoryListBody = document.querySelector('#category-list tbody');
          categoryListBody.innerHTML = ''; // Clear existing categories

          categories.forEach(category => {
              const row = document.createElement('tr');
              row.innerHTML = `
                  <td>${category.name}</td>
                  <td>
                      <button class="btn-secondary btn-small edit-category-button" data-id="${category.id}">Edit</button>
                      <button class="btn-action btn-small delete-category-button" data-id="${category.id}">Delete</button>
                  </td>
              `;
              categoryListBody.appendChild(row);
          });
        };

        // Edit Category Button Functionality (Event Delegation)
        document.querySelector('#category-list tbody').addEventListener('click', (event) => {
            if (event.target.classList.contains('edit-category-button')) {
                const categoryId = event.target.dataset.id;
                const row = event.target.closest('tr');
                const categoryName = row.querySelector('td:first-child').textContent;

                categoryIdInput.value = categoryId;
                categoryNameInput.value = categoryName;
                categoryModal.style.display = 'flex';
            }
        });

        // Initial fetch of categories
        fetchCategories();

        // Fetch categories on page load and populate the select
        const populateCategorySelect = async () => {
            try {
                const token = getAuthToken();
                if (!token) {
                    throw new Error('Not authenticated.');
                }

                const response = await fetch(`${API_BASE_URL}/categories`, {
                    method: 'GET',
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    }
                });

                if (!response.ok) {
                    const errorData = await response.json();
                    throw new Error(errorData.error || `Failed to fetch categories: ${response.status}`);
                }

                const categories = await response.json();
                const categorySelect = document.getElementById('post_categories');
                categorySelect.innerHTML = ''; // Clear existing options

                categories.forEach(category => {
                    const option = document.createElement('option');
                    option.value = category.id;
                    option.textContent = category.name;
                    categorySelect.appendChild(option);
                });

            } catch (error) {
                handleApiError(error);
            }
        };

        populateCategorySelect(); // Call on page load.
    });
})();