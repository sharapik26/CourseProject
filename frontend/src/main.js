// Sensory Navigator - Main JavaScript
// API configuration
const API_URL = 'http://localhost:8080/api';

// State management
let currentUser = null;
let accessToken = localStorage.getItem('accessToken');
let refreshToken = localStorage.getItem('refreshToken');

// DOM Elements
const authScreen = document.getElementById('auth-screen');
const mainScreen = document.getElementById('main-screen');

// Auth forms
const loginForm = document.getElementById('login-form');
const registerForm = document.getElementById('register-form');
const forgotForm = document.getElementById('forgot-form');

// Navigation
const navItems = document.querySelectorAll('.nav-item');
const pages = document.querySelectorAll('.page');

// Initialize application
document.addEventListener('DOMContentLoaded', () => {
  initializeApp();
  setupEventListeners();
});

// Check authentication and initialize
function initializeApp() {
  if (accessToken) {
    fetchUserProfile();
  } else {
    showAuthScreen();
  }
}

// Setup event listeners
function setupEventListeners() {
  // Auth form switches
  document.getElementById('show-register').addEventListener('click', (e) => {
    e.preventDefault();
    switchAuthForm('register');
  });
  
  document.getElementById('show-login').addEventListener('click', (e) => {
    e.preventDefault();
    switchAuthForm('login');
  });
  
  document.getElementById('show-forgot').addEventListener('click', (e) => {
    e.preventDefault();
    switchAuthForm('forgot');
  });
  
  document.getElementById('back-to-login').addEventListener('click', (e) => {
    e.preventDefault();
    switchAuthForm('login');
  });
  
  // Auth form submissions
  loginForm.addEventListener('submit', handleLogin);
  registerForm.addEventListener('submit', handleRegister);
  forgotForm.addEventListener('submit', handleForgotPassword);
  
  // Navigation
  navItems.forEach(item => {
    item.addEventListener('click', (e) => {
      e.preventDefault();
      const page = item.dataset.page;
      navigateToPage(page);
    });
  });
  
  // Logout
  document.getElementById('logout-btn').addEventListener('click', handleLogout);
  
  // Profile form
  document.getElementById('profile-form').addEventListener('submit', handleProfileUpdate);
  
  // Theme toggle
  document.getElementById('setting-theme').addEventListener('change', (e) => {
    document.documentElement.setAttribute('data-theme', e.target.checked ? 'dark' : 'light');
    localStorage.setItem('theme', e.target.checked ? 'dark' : 'light');
  });
  
  // Load saved theme
  const savedTheme = localStorage.getItem('theme') || 'light';
  document.documentElement.setAttribute('data-theme', savedTheme);
  document.getElementById('setting-theme').checked = savedTheme === 'dark';
}

// Switch between auth forms
function switchAuthForm(formName) {
  loginForm.classList.remove('active');
  registerForm.classList.remove('active');
  forgotForm.classList.remove('active');
  
  const formMap = {
    'login': loginForm,
    'register': registerForm,
    'forgot': forgotForm
  };
  
  formMap[formName].classList.add('active');
}

// Handle login
async function handleLogin(e) {
  e.preventDefault();
  
  const email = document.getElementById('login-email').value;
  const password = document.getElementById('login-password').value;
  
  try {
    const response = await fetch(`${API_URL}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password })
    });
    
    const data = await response.json();
    
    if (!response.ok) {
      throw new Error(data.error || 'Ошибка входа');
    }
    
    saveTokens(data.tokens);
    currentUser = data.user;
    showMainScreen();
    showToast('Добро пожаловать!', 'success');
  } catch (error) {
    showToast(error.message, 'error');
  }
}

// Handle registration
async function handleRegister(e) {
  e.preventDefault();
  
  const username = document.getElementById('register-username').value;
  const email = document.getElementById('register-email').value;
  const password = document.getElementById('register-password').value;
  
  try {
    const response = await fetch(`${API_URL}/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, email, password })
    });
    
    const data = await response.json();
    
    if (!response.ok) {
      throw new Error(data.error || 'Ошибка регистрации');
    }
    
    saveTokens(data.tokens);
    currentUser = data.user;
    showMainScreen();
    showToast('Регистрация успешна!', 'success');
  } catch (error) {
    showToast(error.message, 'error');
  }
}

// Handle forgot password
async function handleForgotPassword(e) {
  e.preventDefault();
  
  const email = document.getElementById('forgot-email').value;
  
  try {
    const response = await fetch(`${API_URL}/auth/forgot-password`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email })
    });
    
    showToast('Если email существует, вы получите ссылку для восстановления', 'success');
    switchAuthForm('login');
  } catch (error) {
    showToast('Произошла ошибка', 'error');
  }
}

// Handle logout
function handleLogout() {
  localStorage.removeItem('accessToken');
  localStorage.removeItem('refreshToken');
  accessToken = null;
  refreshToken = null;
  currentUser = null;
  showAuthScreen();
  showToast('Вы вышли из аккаунта', 'success');
}

// Handle profile update
async function handleProfileUpdate(e) {
  e.preventDefault();
  
  const username = document.getElementById('profile-username').value;
  const birthdate = document.getElementById('profile-birthdate').value;
  
  const updateData = {};
  if (username) updateData.username = username;
  if (birthdate) updateData.birth_date = birthdate;
  
  try {
    const response = await apiRequest('/users/me', 'PUT', updateData);
    
    if (response.error) {
      throw new Error(response.error);
    }
    
    currentUser = response;
    updateUserUI();
    showToast('Профиль обновлён', 'success');
  } catch (error) {
    showToast(error.message, 'error');
  }
}

// Fetch user profile
async function fetchUserProfile() {
  try {
    const response = await apiRequest('/users/me');
    
    if (response.error) {
      throw new Error(response.error);
    }
    
    currentUser = response;
    showMainScreen();
  } catch (error) {
    // Token might be expired
    handleLogout();
  }
}

// Navigate to page
function navigateToPage(pageName) {
  // Update nav
  navItems.forEach(item => {
    item.classList.toggle('active', item.dataset.page === pageName);
  });
  
  // Update pages
  pages.forEach(page => {
    page.classList.toggle('active', page.id === `${pageName}-page`);
  });
  
  // Load page data
  switch (pageName) {
    case 'reviews':
      loadMyReviews();
      break;
    case 'favorites':
      loadMyFavorites();
      break;
  }
}

// Load my reviews
async function loadMyReviews() {
  const container = document.getElementById('reviews-list');
  container.innerHTML = '<p style="text-align:center;color:var(--color-text-muted)">Загрузка...</p>';
  
  try {
    const response = await apiRequest('/users/me/reviews');
    
    if (!response.reviews || response.reviews.length === 0) {
      container.innerHTML = `
        <div class="empty-state">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
          </svg>
          <p>У вас пока нет отзывов</p>
        </div>
      `;
      return;
    }
    
    container.innerHTML = response.reviews.map(review => `
      <div class="review-card">
        <div class="review-header">
          <span class="review-place">${review.place_name || 'Место #' + review.place_id}</span>
          <span class="review-date">${formatDate(review.created_at)}</span>
        </div>
        <div class="review-ratings">
          ${review.sensory_rating ? `<span class="rating-badge">Сенсорность: ${review.sensory_rating}/5</span>` : ''}
          ${review.lighting_rating ? `<span class="rating-badge">Освещение: ${review.lighting_rating}/5</span>` : ''}
          ${review.sound_level_rating ? `<span class="rating-badge">Громкость: ${review.sound_level_rating}/5</span>` : ''}
          ${review.crowding_rating ? `<span class="rating-badge">Людность: ${review.crowding_rating}/5</span>` : ''}
          ${review.accessibility_rating ? `<span class="rating-badge">Доступность: ${review.accessibility_rating}/5</span>` : ''}
        </div>
        ${review.text ? `<p class="review-text">${review.text}</p>` : ''}
      </div>
    `).join('');
  } catch (error) {
    container.innerHTML = '<p style="text-align:center;color:var(--color-error)">Ошибка загрузки</p>';
  }
}

// Load my favorites
async function loadMyFavorites() {
  const container = document.getElementById('favorites-list');
  container.innerHTML = '<p style="text-align:center;color:var(--color-text-muted)">Загрузка...</p>';
  
  try {
    const response = await apiRequest('/users/me/favorites');
    
    if (!response.favorites || response.favorites.length === 0) {
      container.innerHTML = `
        <div class="empty-state">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
          </svg>
          <p>Список избранного пуст</p>
        </div>
      `;
      return;
    }
    
    container.innerHTML = response.favorites.map(fav => `
      <div class="favorite-card">
        <div class="review-header">
          <span class="review-place">${fav.place_name}</span>
          <span class="review-date">${formatDate(fav.created_at)}</span>
        </div>
        ${fav.address ? `<p class="review-text">${fav.address}</p>` : ''}
        ${fav.category ? `<span class="rating-badge">${fav.category}</span>` : ''}
      </div>
    `).join('');
  } catch (error) {
    container.innerHTML = '<p style="text-align:center;color:var(--color-error)">Ошибка загрузки</p>';
  }
}

// Show auth screen
function showAuthScreen() {
  authScreen.classList.remove('hidden');
  mainScreen.classList.add('hidden');
}

// Show main screen
function showMainScreen() {
  authScreen.classList.add('hidden');
  mainScreen.classList.remove('hidden');
  updateUserUI();
  navigateToPage('profile');
}

// Update user UI
function updateUserUI() {
  if (!currentUser) return;
  
  // Update sidebar
  document.getElementById('user-name').textContent = currentUser.username;
  document.getElementById('user-email').textContent = currentUser.email;
  
  const firstLetter = currentUser.username.charAt(0).toUpperCase();
  document.getElementById('avatar-letter').textContent = firstLetter;
  document.getElementById('profile-avatar-letter').textContent = firstLetter;
  
  // Update profile form
  document.getElementById('profile-username').value = currentUser.username;
  document.getElementById('profile-email').value = currentUser.email;
  if (currentUser.birth_date) {
    document.getElementById('profile-birthdate').value = currentUser.birth_date;
  }
}

// API request helper with token refresh
async function apiRequest(endpoint, method = 'GET', body = null) {
  const options = {
    method,
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${accessToken}`
    }
  };
  
  if (body) {
    options.body = JSON.stringify(body);
  }
  
  let response = await fetch(`${API_URL}${endpoint}`, options);
  
  // If unauthorized, try to refresh token
  if (response.status === 401 && refreshToken) {
    const refreshed = await refreshAccessToken();
    if (refreshed) {
      options.headers['Authorization'] = `Bearer ${accessToken}`;
      response = await fetch(`${API_URL}${endpoint}`, options);
    }
  }
  
  return response.json();
}

// Refresh access token
async function refreshAccessToken() {
  try {
    const response = await fetch(`${API_URL}/auth/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refreshToken })
    });
    
    if (!response.ok) {
      throw new Error('Token refresh failed');
    }
    
    const data = await response.json();
    saveTokens(data);
    return true;
  } catch (error) {
    handleLogout();
    return false;
  }
}

// Save tokens
function saveTokens(tokens) {
  accessToken = tokens.access_token;
  refreshToken = tokens.refresh_token;
  localStorage.setItem('accessToken', accessToken);
  localStorage.setItem('refreshToken', refreshToken);
}

// Format date
function formatDate(dateString) {
  const date = new Date(dateString);
  return date.toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric'
  });
}

// Show toast notification
function showToast(message, type = 'success') {
  const container = document.getElementById('toast-container');
  const toast = document.createElement('div');
  toast.className = `toast ${type}`;
  toast.textContent = message;
  
  container.appendChild(toast);
  
  setTimeout(() => {
    toast.style.animation = 'slideIn 0.3s ease reverse';
    setTimeout(() => toast.remove(), 300);
  }, 3000);
}

// Export for Tauri
window.app = {
  showToast,
  navigateToPage,
  handleLogout
};

