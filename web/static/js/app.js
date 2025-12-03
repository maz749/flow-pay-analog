// API Base URL
const API_URL = '/api';

// State
let currentUser = null;
let currentSubscription = null;
let subscriptions = [];
let selectedSubscriptionTemplate = null;
let currentCategoryFilter = 'all';
let userSubscriptionsFilter = 'all';
let userSubscriptionsSearchQuery = '';

// Charts instances
let monthlyExpensesChart = null;
let categoryChart = null;
let topSubscriptionsChart = null;
let upcomingPaymentsChart = null;

// Inline charts instances
let inlineMonthlyChart = null;
let inlineCategoryChart = null;

// Subscription Templates Database
const subscriptionTemplates = [
    // Популярные сервисы (8 сервисов)
    { name: 'Apple Music', logo: '🎼', category: 'music', price: 169, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'ChatGPT Plus', logo: '🤖', category: 'software', price: 20, currency: 'USD', period: 'monthly', popular: true },
    { name: 'GitHub Copilot', logo: '💻', category: 'software', price: 10, currency: 'USD', period: 'monthly', popular: true },
    { name: 'Kinopoisk', logo: '🎥', category: 'streaming', price: 399, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Netflix', logo: '🎬', category: 'streaming', price: 649, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Spotify', logo: '🎵', category: 'music', price: 169, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Yandex 360', logo: '☁️', category: 'cloud', price: 199, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'YouTube Premium', logo: '📺', category: 'streaming', price: 399, currency: 'RUB', period: 'monthly', popular: true },

    // Маркетплейсы - Аналитика и статистика (46 сервисов)
    { name: 'MPSTATS', logo: '📊', category: 'marketplace', price: 5990, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Маяк', logo: '🔍', category: 'marketplace', price: 15990, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'MarketGuru', logo: '🎯', category: 'marketplace', price: 9990, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'EGGHEADS', logo: '🥚', category: 'marketplace', price: 25000, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Wildbox', logo: '📦', category: 'marketplace', price: 3500, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Stat4Market', logo: '📈', category: 'marketplace', price: 2500, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'SellerStats', logo: '📉', category: 'marketplace', price: 2490, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'woysa.club', logo: '🎪', category: 'marketplace', price: 10990, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'XWAY', logo: '🛤️', category: 'marketplace', price: 7900, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Sellmonitor', logo: '📱', category: 'marketplace', price: 19000, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'WBStat.PRO', logo: '💼', category: 'marketplace', price: 2990, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'SellerFox', logo: '🦊', category: 'marketplace', price: 1950, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Topseller', logo: '🏆', category: 'marketplace', price: 3500, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Pi-Data', logo: '🥧', category: 'marketplace', price: 2900, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'SelSup', logo: '🛒', category: 'marketplace', price: 4000, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Модульселлер', logo: '🧩', category: 'marketplace', price: 2490, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Salist', logo: '📋', category: 'marketplace', price: 4900, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Тинькофф Селлер', logo: '💳', category: 'marketplace', price: 10000, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'ADAPTER', logo: '🔌', category: 'marketplace', price: 5000, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'MPSPACE', logo: '🚀', category: 'marketplace', price: 6900, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Huckster', logo: '🎰', category: 'marketplace', price: 10900, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'MP SURF', logo: '🏄', category: 'marketplace', price: 890, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'inSales', logo: '💰', category: 'marketplace', price: 1490, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'SalesFinder', logo: '🔎', category: 'marketplace', price: 3500, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'MP Manager', logo: '⚙️', category: 'marketplace', price: 8970, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Smart Seller', logo: '🧠', category: 'marketplace', price: 4900, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Indeepa', logo: '🌊', category: 'marketplace', price: 18990, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Точка маркетплейсы', logo: '🎯', category: 'marketplace', price: 800, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'JVO', logo: '📡', category: 'marketplace', price: 30000, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'MarketProvider', logo: '🏪', category: 'marketplace', price: 5000, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Ракета', logo: '🚀', category: 'marketplace', price: 10000, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'SoykaSoft', logo: '🐦', category: 'marketplace', price: 14990, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Mpfinassist', logo: '💵', category: 'marketplace', price: 2990, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'WB FIN', logo: '💸', category: 'marketplace', price: 3000, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Таблички', logo: '📊', category: 'marketplace', price: 7000, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Нейромаркет', logo: '🧠', category: 'marketplace', price: 1690, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Sellego', logo: '🎲', category: 'marketplace', price: 3990, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'ОТВЕТО', logo: '💬', category: 'marketplace', price: 7990, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Отвечумба', logo: '🤖', category: 'marketplace', price: 1400, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'ТОРГСТАТ', logo: '📊', category: 'marketplace', price: 3990, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'MPfact', logo: '📌', category: 'marketplace', price: 2850, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'TrueStats', logo: '✅', category: 'marketplace', price: 2999, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'MP Profit', logo: '💹', category: 'marketplace', price: 2399, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'RASK', logo: '🎨', category: 'marketplace', price: 7500, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'TACTICs', logo: '♟️', category: 'marketplace', price: 7000, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Спикс', logo: '🗣️', category: 'marketplace', price: 4900, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Джем', logo: '🍓', category: 'marketplace', price: 18990, currency: 'RUB', period: 'monthly', popular: true },

    // CRM системы (3 сервиса)
    { name: 'AMOCRM', logo: '📈', category: 'crm', price: 1199, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Мегаплан', logo: '📋', category: 'crm', price: 450, currency: 'RUB', period: 'monthly', popular: true },
    { name: 'Битрикс24', logo: '📞', category: 'crm', price: 2490, currency: 'RUB', period: 'monthly', popular: true },

    // Остальные сервисы
    { name: 'Adobe Creative Cloud', logo: '🎨', category: 'software', price: 54.99, currency: 'USD', period: 'monthly', popular: false },
    { name: 'Amazon Prime', logo: '📦', category: 'marketplace', price: 14.99, currency: 'USD', period: 'monthly', popular: false },
    { name: 'Asana', logo: '✅', category: 'software', price: 10.99, currency: 'USD', period: 'monthly', popular: false },
    { name: 'Canva Pro', logo: '🖼️', category: 'software', price: 12.99, currency: 'USD', period: 'monthly', popular: false },
    { name: 'Coursera Plus', logo: '🎓', category: 'education', price: 59, currency: 'USD', period: 'monthly', popular: false },
    { name: 'Deezer', logo: '🎶', category: 'music', price: 169, currency: 'RUB', period: 'monthly', popular: false },
    { name: 'Discord Nitro', logo: '💬', category: 'gaming', price: 9.99, currency: 'USD', period: 'monthly', popular: false },
    { name: 'Dropbox', logo: '📦', category: 'cloud', price: 9.99, currency: 'USD', period: 'monthly', popular: false },
    { name: 'Duolingo Plus', logo: '🦉', category: 'education', price: 6.99, currency: 'USD', period: 'monthly', popular: false },
    { name: 'Evernote', logo: '📓', category: 'software', price: 7.99, currency: 'USD', period: 'monthly', popular: false },
    { name: 'Figma', logo: '🎯', category: 'software', price: 12, currency: 'USD', period: 'monthly', popular: false },
    { name: 'GeForce NOW', logo: '🎮', category: 'gaming', price: 999, currency: 'RUB', period: 'monthly', popular: false },
    { name: 'Google One', logo: '💾', category: 'cloud', price: 139, currency: 'RUB', period: 'monthly', popular: false },
    { name: 'Grammarly', logo: '✍️', category: 'software', price: 12, currency: 'USD', period: 'monthly', popular: false },
    { name: 'HBO Max', logo: '🎭', category: 'streaming', price: 9.99, currency: 'USD', period: 'monthly', popular: false },
    { name: 'Headspace', logo: '🧘', category: 'education', price: 12.99, currency: 'USD', period: 'monthly', popular: false },
    { name: 'HubSpot', logo: '🚀', category: 'crm', price: 50, currency: 'USD', period: 'monthly', popular: false },
    { name: 'iCloud+', logo: '☁️', category: 'cloud', price: 149, currency: 'RUB', period: 'monthly', popular: false },
    { name: 'ivi', logo: '🎞️', category: 'streaming', price: 399, currency: 'RUB', period: 'monthly', popular: false },
    { name: 'LinkedIn Premium', logo: '💼', category: 'software', price: 29.99, currency: 'USD', period: 'monthly', popular: false },
    { name: 'Мегафон ТВ', logo: '📱', category: 'streaming', price: 299, currency: 'RUB', period: 'monthly', popular: false },
    { name: 'Microsoft 365', logo: '📊', category: 'software', price: 7, currency: 'USD', period: 'monthly', popular: false },
    { name: 'МТС Premium', logo: '📡', category: 'streaming', price: 399, currency: 'RUB', period: 'monthly', popular: false },
    { name: 'Netflix Basic', logo: '🎬', category: 'streaming', price: 449, currency: 'RUB', period: 'monthly', popular: false },
    { name: 'Nintendo Switch Online', logo: '🕹️', category: 'gaming', price: 299, currency: 'RUB', period: 'monthly', popular: false },
    { name: 'Notion', logo: '📝', category: 'software', price: 8, currency: 'USD', period: 'monthly', popular: false },
    { name: 'Okko', logo: '📹', category: 'streaming', price: 599, currency: 'RUB', period: 'monthly', popular: false },
    { name: 'Ozon Premium', logo: '🛒', category: 'marketplace', price: 199, currency: 'RUB', period: 'monthly', popular: false },
    { name: 'PlayStation Plus', logo: '🎮', category: 'gaming', price: 599, currency: 'RUB', period: 'monthly', popular: false },
    { name: 'Salesforce', logo: '🌐', category: 'crm', price: 25, currency: 'USD', period: 'monthly', popular: false },
    { name: 'Setka', logo: '📰', category: 'software', price: 15, currency: 'USD', period: 'monthly', popular: false },
    { name: 'Skillbox', logo: '📚', category: 'education', price: 3990, currency: 'RUB', period: 'monthly', popular: false },
    { name: 'Slack', logo: '💬', category: 'software', price: 6.67, currency: 'USD', period: 'monthly', popular: false },
    { name: 'START', logo: '▶️', category: 'streaming', price: 349, currency: 'RUB', period: 'monthly', popular: false },
    { name: 'Telegram Premium', logo: '✈️', category: 'software', price: 599, currency: 'RUB', period: 'monthly', popular: false },
    { name: 'Tidal', logo: '🎵', category: 'music', price: 9.99, currency: 'USD', period: 'monthly', popular: false },
    { name: 'Tinkoff Банк Pro', logo: '💳', category: 'marketplace', price: 299, currency: 'RUB', period: 'monthly', popular: false },
    { name: 'Trello', logo: '📋', category: 'software', price: 5, currency: 'USD', period: 'monthly', popular: false },
    { name: 'VK Музыка', logo: '🎧', category: 'music', price: 199, currency: 'RUB', period: 'monthly', popular: false },
    { name: 'Wildberries Premium', logo: '🛍️', category: 'marketplace', price: 199, currency: 'RUB', period: 'monthly', popular: false },
    { name: 'Xbox Game Pass', logo: '🎯', category: 'gaming', price: 499, currency: 'RUB', period: 'monthly', popular: false },
    { name: 'Yandex Music', logo: '🎼', category: 'music', price: 199, currency: 'RUB', period: 'monthly', popular: false },
    { name: 'Yandex Plus', logo: '🟡', category: 'marketplace', price: 299, currency: 'RUB', period: 'monthly', popular: false },
    { name: 'Zoom Pro', logo: '🎥', category: 'software', price: 14.99, currency: 'USD', period: 'monthly', popular: false },
    { name: 'Альфа-Банк Premium', logo: '🏦', category: 'marketplace', price: 199, currency: 'RUB', period: 'monthly', popular: false },
    { name: 'Литрес Библиотека', logo: '📚', category: 'education', price: 399, currency: 'RUB', period: 'monthly', popular: false },
    { name: 'Ростелеком Ключ', logo: '🔑', category: 'streaming', price: 249, currency: 'RUB', period: 'monthly', popular: false }
];

// Init
document.addEventListener('DOMContentLoaded', () => {
    initApp();
});

function initApp() {
    const token = localStorage.getItem('token');

    if (token) {
        loadApp();
    } else {
        showScreen('landing-screen');
    }

    setupEventListeners();
}

function setupEventListeners() {
    // Landing page buttons
    document.getElementById('landing-login-btn')?.addEventListener('click', openLoginModal);
    document.getElementById('landing-register-btn')?.addEventListener('click', openRegisterModal);

    // Modal switchers
    document.getElementById('switch-to-register')?.addEventListener('click', (e) => {
        e.preventDefault();
        closeModal('login-modal');
        openRegisterModal();
    });

    document.getElementById('switch-to-login')?.addEventListener('click', (e) => {
        e.preventDefault();
        closeModal('register-modal');
        openLoginModal();
    });

    // Auth forms
    document.getElementById('login-form')?.addEventListener('submit', handleLogin);
    document.getElementById('register-form')?.addEventListener('submit', handleRegister);

    // App buttons
    document.getElementById('logout-btn')?.addEventListener('click', handleLogout);
    document.getElementById('add-subscription-btn')?.addEventListener('click', () => openSubscriptionModal());
    document.getElementById('settings-btn')?.addEventListener('click', openSettingsModal);
    document.getElementById('stats-page-btn')?.addEventListener('click', openStatsPage);
    document.getElementById('theme-toggle-btn')?.addEventListener('click', toggleTheme);
    document.getElementById('back-to-app-btn')?.addEventListener('click', () => showScreen('app-screen'));
    document.getElementById('logout-btn-stats')?.addEventListener('click', handleLogout);

    // Subscription form
    document.getElementById('subscription-form')?.addEventListener('submit', handleSaveSubscription);
    const subPeriodEl = document.getElementById('sub-period');
    if (subPeriodEl) {
        subPeriodEl.addEventListener('change', (e) => {
            const customDaysGroup = document.getElementById('custom-days-group');
            if (customDaysGroup) {
                customDaysGroup.style.display = e.target.value === 'custom' ? 'block' : 'none';
            }
        });
    }

    // Subscription selection mode buttons
    document.querySelectorAll('.mode-btn').forEach(btn => {
        btn.addEventListener('click', handleModeSwitch);
    });

    // Search functionality
    document.getElementById('subscription-search')?.addEventListener('input', handleSubscriptionSearch);

    // Category filter buttons
    document.querySelectorAll('.filter-btn').forEach(btn => {
        btn.addEventListener('click', handleCategoryFilter);
    });

    // User subscriptions filters
    document.getElementById('user-subscription-search')?.addEventListener('input', handleUserSubscriptionSearch);
    document.querySelectorAll('.user-filter-btn').forEach(btn => {
        btn.addEventListener('click', handleUserCategoryFilter);
    });

    // Telegram settings form
    document.getElementById('telegram-settings-form')?.addEventListener('submit', handleSaveTelegramSettings);

    // Modal close buttons
    document.querySelectorAll('.modal-close').forEach(btn => {
        btn.addEventListener('click', (e) => {
            e.target.closest('.modal').classList.remove('active');
        });
    });

    // Close modal on backdrop click
    document.querySelectorAll('.modal').forEach(modal => {
        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                modal.classList.remove('active');
            }
        });
    });
}

function showScreen(screenId) {
    document.querySelectorAll('.screen').forEach(screen => screen.classList.remove('active'));
    const screen = document.getElementById(screenId);
    if (screen) {
        screen.classList.add('active');
    }
}

function openModal(modalId) {
    document.getElementById(modalId).classList.add('active');
}

function closeModal(modalId) {
    document.getElementById(modalId).classList.remove('active');
}

function openLoginModal() {
    openModal('login-modal');
    document.getElementById('login-form').reset();
}

function openRegisterModal() {
    openModal('register-modal');
    document.getElementById('register-form').reset();
}

// Auth handlers
async function handleLogin(e) {
    e.preventDefault();

    const email = document.getElementById('login-email').value;
    const password = document.getElementById('login-password').value;

    console.log('Attempting login with email:', email);

    try {
        const response = await fetch(`${API_URL}/auth/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });

        console.log('Response status:', response.status);

        const data = await response.json();
        console.log('Response data:', data);

        if (data.success) {
            console.log('Login successful, saving token...');
            localStorage.setItem('token', data.data.token);
            currentUser = data.data.user;
            closeModal('login-modal');
            console.log('Loading app...');
            loadApp();
        } else {
            console.error('Login failed:', data.error);
            alert(data.error || 'Ошибка входа');
        }
    } catch (error) {
        console.error('Login error:', error);
        alert('Ошибка подключения к серверу: ' + error.message);
    }
}

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

        if (data.success) {
            localStorage.setItem('token', data.data.token);
            currentUser = data.data.user;
            closeModal('register-modal');
            loadApp();
        } else {
            alert(data.error || 'Ошибка регистрации');
        }
    } catch (error) {
        console.error('Register error:', error);
        alert('Ошибка подключения к серверу');
    }
}

function handleLogout() {
    localStorage.removeItem('token');
    currentUser = null;
    subscriptions = [];
    showScreen('landing-screen');
}

// App loading
async function loadApp() {
    console.log('loadApp() called');
    try {
        console.log('Switching to app-screen...');
        showScreen('app-screen');

        // Show user name
        if (currentUser) {
            document.getElementById('user-name').textContent = currentUser.username;
        }

        console.log('Loading stats...');
        await loadStats();
        console.log('Loading subscriptions...');
        await loadSubscriptions();
        console.log('App loaded successfully!');
    } catch (error) {
        console.error('Error loading app:', error);
    }
}

async function loadStats() {
    try {
        const response = await apiRequest('/subscriptions/stats');

        if (response.success) {
            const stats = response.data;
            const totalSubs = document.getElementById('total-subscriptions');
            const monthlyTotal = document.getElementById('monthly-total');
            const yearlyTotal = document.getElementById('yearly-total');
            const nextPayment = document.getElementById('next-payment');

            if (totalSubs) totalSubs.textContent = stats.total_subscriptions;
            if (monthlyTotal) monthlyTotal.textContent = `${stats.monthly_total.toFixed(2)} ₽`;
            if (yearlyTotal) yearlyTotal.textContent = `${stats.yearly_total.toFixed(2)} ₽`;

            if (nextPayment) {
                if (stats.next_payment) {
                    const date = new Date(stats.next_payment.next_billing_date);
                    const dateStr = date.toLocaleDateString('ru-RU');
                    nextPayment.textContent = `${stats.next_payment.name} (${dateStr})`;
                } else {
                    nextPayment.textContent = '-';
                }
            }
        }
    } catch (error) {
        console.error('Failed to load stats:', error);
    }
}

async function loadSubscriptions() {
    try {
        const response = await apiRequest('/subscriptions');

        if (response.success) {
            subscriptions = response.data || [];
            renderSubscriptions();
            // Create inline charts
            createInlineCharts();
        }
    } catch (error) {
        console.error('Failed to load subscriptions:', error);
    }
}

function renderSubscriptions() {
    const list = document.getElementById('subscriptions-list');
    const emptyState = document.getElementById('empty-state');
    const filtersContainer = document.getElementById('user-subscriptions-filters');

    // If elements don't exist (e.g., on teams page), skip rendering
    if (!list || !emptyState) return;

    // Show/hide filters based on whether there are subscriptions
    if (filtersContainer) {
        filtersContainer.style.display = subscriptions.length > 0 ? 'block' : 'none';
    }

    // Apply filters
    let filtered = [...subscriptions];

    // Filter by category
    if (userSubscriptionsFilter !== 'all') {
        filtered = filtered.filter(sub => sub.category === userSubscriptionsFilter);
    }

    // Filter by search query
    if (userSubscriptionsSearchQuery) {
        filtered = filtered.filter(sub =>
            sub.name.toLowerCase().includes(userSubscriptionsSearchQuery)
        );
    }

    if (filtered.length === 0) {
        list.innerHTML = '';
        emptyState.classList.add('active');

        // Update empty state message based on filters
        const emptyIcon = emptyState.querySelector('.empty-icon');
        const emptyTitle = emptyState.querySelector('h3');
        const emptyText = emptyState.querySelector('p');

        if (subscriptions.length > 0 && (userSubscriptionsFilter !== 'all' || userSubscriptionsSearchQuery)) {
            emptyIcon.textContent = '🔍';
            emptyTitle.textContent = 'Не найдено';
            emptyText.textContent = 'Попробуйте изменить фильтры или поисковый запрос';
        } else {
            emptyIcon.textContent = '📭';
            emptyTitle.textContent = 'Нет подписок';
            emptyText.textContent = 'Добавьте свою первую подписку, чтобы начать отслеживание платежей';
        }
        return;
    }

    emptyState.classList.remove('active');

    list.innerHTML = filtered.map(sub => {
        const nextDate = new Date(sub.next_billing_date);
        const dateStr = nextDate.toLocaleDateString('ru-RU');
        const daysUntil = Math.ceil((nextDate - new Date()) / (1000 * 60 * 60 * 24));

        const periodLabels = {
            monthly: 'Ежемесячно',
            yearly: 'Ежегодно',
            weekly: 'Еженедельно',
            custom: `Каждые ${sub.custom_period_days} дней`
        };

        const currencySymbols = {
            RUB: '₽',
            USD: '$',
            EUR: '€'
        };

        const icon = getSubscriptionIcon(sub.category);
        const isActive = sub.is_active;

        return `
            <div class="subscription-card">
                <div class="subscription-card-header">
                    <div class="subscription-icon" style="${sub.color ? `background: ${sub.color}` : ''}">
                        ${icon}
                    </div>
                    <div class="subscription-actions-top">
                        <button class="btn btn-secondary btn-sm" onclick="editSubscription(${sub.id})">Изменить</button>
                        <button class="btn btn-warning btn-sm" onclick="cancelSubscriptionModal(${sub.id}, '${sub.name.replace(/'/g, "\\'")}')">Отменить</button>
                        <button class="btn btn-danger btn-sm" onclick="deleteSubscription(${sub.id})">Удалить</button>
                    </div>
                </div>
                <div class="subscription-info">
                    <div class="subscription-name">
                        ${sub.name}
                        ${!isActive ? '<span class="badge badge-danger">Неактивна</span>' : ''}
                    </div>
                    <div class="subscription-meta">
                        <span>${periodLabels[sub.billing_period] || sub.billing_period}</span>
                        ${sub.category ? `<span>• ${getCategoryLabel(sub.category)}</span>` : ''}
                        ${sub.description ? `<span>• ${sub.description}</span>` : ''}
                    </div>
                </div>
                <div class="subscription-footer">
                    <div class="subscription-amount">
                        ${sub.amount.toFixed(2)} ${currencySymbols[sub.currency] || sub.currency}
                    </div>
                    <div class="subscription-next-date">
                        Списание: ${dateStr} ${daysUntil >= 0 ? `(через ${daysUntil} дн.)` : '(просрочено)'}
                    </div>
                </div>
            </div>
        `;
    }).join('');
}

function getSubscriptionIcon(category) {
    const icons = {
        streaming: '📺',
        music: '🎵',
        software: '💻',
        gaming: '🎮',
        education: '📚',
        cloud: '☁️',
        marketplace: '🛒',
        crm: '📊',
        other: '📦'
    };
    return icons[category] || '💳';
}

function getCategoryLabel(category) {
    const labels = {
        streaming: 'Стриминг',
        music: 'Музыка',
        software: 'Софт',
        gaming: 'Игры',
        education: 'Образование',
        cloud: 'Облако',
        marketplace: 'Маркетплейсы',
        crm: 'CRM',
        other: 'Другое'
    };
    return labels[category] || category;
}

// Subscription modal
function openSubscriptionModal(subscription = null) {
    currentSubscription = subscription;
    const modal = document.getElementById('subscription-modal');
    const form = document.getElementById('subscription-form');
    const title = document.getElementById('modal-title');
    const modeSelector = document.getElementById('mode-selector');

    form.reset();

    // Clear readonly/disabled state from previous selections
    clearFieldRestrictions();

    if (subscription) {
        // Editing existing subscription - show only manual mode
        title.textContent = 'Редактировать подписку';
        modeSelector.style.display = 'none';

        // Show only manual mode
        document.querySelectorAll('.selection-mode').forEach(section => {
            section.classList.remove('active');
        });
        document.getElementById('manual-mode').classList.add('active');

        // Fill form with existing data
        document.getElementById('sub-name').value = subscription.name;
        document.getElementById('sub-amount').value = subscription.amount;
        document.getElementById('sub-currency').value = subscription.currency;
        document.getElementById('sub-description').value = subscription.description || '';
        document.getElementById('sub-period').value = subscription.billing_period;
        document.getElementById('sub-start-date').value = subscription.start_date.split('T')[0];
        document.getElementById('sub-category').value = subscription.category || '';

        if (subscription.billing_period === 'custom') {
            document.getElementById('custom-days-group').style.display = 'block';
            document.getElementById('sub-custom-days').value = subscription.custom_period_days;
        }

        // Set notification checkboxes
        if (subscription.notifications) {
            subscription.notifications.forEach(notif => {
                const checkbox = document.querySelector(`input[name="notify"][value="${notif.notify_days_before}"]`);
                if (checkbox) checkbox.checked = notif.is_enabled;
            });
        }
    } else {
        // Adding new subscription - show mode selector and list
        title.textContent = 'Добавить подписку';
        modeSelector.style.display = 'flex';

        // Reset to "from list" mode
        document.querySelectorAll('.mode-btn').forEach(btn => {
            btn.classList.remove('active');
        });
        document.querySelector('.mode-btn[data-mode="list"]').classList.add('active');

        document.querySelectorAll('.selection-mode').forEach(section => {
            section.classList.remove('active');
        });
        document.getElementById('from-list-mode').classList.add('active');

        // Render subscription templates
        renderSubscriptionTemplates();

        // Set default date
        const today = new Date().toISOString().split('T')[0];
        document.getElementById('sub-start-date').value = today;
    }

    modal.classList.add('active');
}

async function handleSaveSubscription(e) {
    e.preventDefault();

    // Temporarily enable category field to get its value (disabled fields don't submit)
    const categoryField = document.getElementById('sub-category');
    const wasDisabled = categoryField.disabled;
    if (wasDisabled) {
        categoryField.disabled = false;
    }

    const name = document.getElementById('sub-name').value;
    const amount = parseFloat(document.getElementById('sub-amount').value);
    const currency = document.getElementById('sub-currency').value;
    const description = document.getElementById('sub-description').value || null;
    const billing_period = document.getElementById('sub-period').value;
    const start_date = document.getElementById('sub-start-date').value;
    const category = categoryField.value || null;

    const custom_period_days = billing_period === 'custom'
        ? parseInt(document.getElementById('sub-custom-days').value)
        : null;

    const notify_days_before = Array.from(document.querySelectorAll('input[name="notify"]:checked'))
        .map(cb => parseInt(cb.value));

    const data = {
        name,
        amount,
        currency,
        description,
        billing_period,
        start_date,
        category,
        custom_period_days,
        notify_days_before
    };

    try {
        let response;
        if (currentSubscription) {
            response = await apiRequest(`/subscriptions/${currentSubscription.id}`, 'PUT', data);
        } else {
            response = await apiRequest('/subscriptions', 'POST', data);
        }

        if (response.success) {
            document.getElementById('subscription-modal').classList.remove('active');
            // ВАЖНО: Обновляем и статистику и подписки
            await loadStats();
            await loadSubscriptions();
            alert('Подписка успешно сохранена!');
        } else {
            alert(response.error || 'Ошибка сохранения');
        }
    } catch (error) {
        console.error('Save subscription error:', error);
        alert('Ошибка сохранения подписки');
    }
}

async function editSubscription(id) {
    const subscription = subscriptions.find(s => s.id === id);
    if (subscription) {
        openSubscriptionModal(subscription);
    }
}

async function deleteSubscription(id) {
    if (!confirm('Удалить эту подписку?')) return;

    try {
        const response = await apiRequest(`/subscriptions/${id}`, 'DELETE');

        if (response.success) {
            // ВАЖНО: Обновляем и статистику и подписки
            await loadStats();
            await loadSubscriptions();
            alert('Подписка удалена');
        } else {
            alert(response.error || 'Ошибка удаления');
        }
    } catch (error) {
        console.error('Delete subscription error:', error);
        alert('Ошибка удаления подписки');
    }
}

// Settings modal
async function openSettingsModal() {
    const modal = document.getElementById('settings-modal');
    modal.classList.add('active');

    try {
        const response = await apiRequest('/telegram/settings');

        if (response.success) {
            const settings = response.data;
            document.getElementById('telegram-chat-id').value = settings.chat_id || '';
            document.getElementById('telegram-enabled').checked = settings.is_enabled;
        }
    } catch (error) {
        console.error('Failed to load telegram settings:', error);
    }
}

async function handleSaveTelegramSettings(e) {
    e.preventDefault();

    const chatIdValue = document.getElementById('telegram-chat-id').value.trim();
    const chat_id = chatIdValue ? parseInt(chatIdValue) : null;
    const is_enabled = document.getElementById('telegram-enabled').checked;

    try {
        const response = await apiRequest('/telegram/settings', 'PUT', {
            chat_id,
            is_enabled
        });

        if (response.success) {
            alert('Настройки сохранены');
            document.getElementById('settings-modal').classList.remove('active');
        } else {
            alert(response.error || 'Ошибка сохранения');
        }
    } catch (error) {
        console.error('Save telegram settings error:', error);
        alert('Ошибка сохранения настроек');
    }
}

// API helper
async function apiRequest(endpoint, method = 'GET', body = null) {
    const token = localStorage.getItem('token');

    const options = {
        method,
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        }
    };

    if (body) {
        options.body = JSON.stringify(body);
    }

    const response = await fetch(`${API_URL}${endpoint}`, options);

    if (response.status === 401) {
        handleLogout();
        throw new Error('Unauthorized');
    }

    return await response.json();
}

// ============================================
// STATISTICS PAGE FUNCTIONS
// ============================================

async function openStatsPage() {
    showScreen('stats-screen');
    await loadStatsData();
}

async function loadStatsData() {
    try {
        // Load subscriptions
        const response = await apiRequest('/subscriptions');

        if (response.success) {
            const subs = response.data;

            // Load stats summary
            const statsResponse = await apiRequest('/subscriptions/stats');
            if (statsResponse.success) {
                updateStatsSummary(statsResponse.data);
            }

            // Create charts
            createMonthlyExpensesChart(subs);
            createCategoryChart(subs);
            createTopSubscriptionsChart(subs);
            createUpcomingPaymentsChart(subs);

            // Fill table
            fillStatsTable(subs);
        }
    } catch (error) {
        console.error('Failed to load stats data:', error);
    }
}

function updateStatsSummary(stats) {
    document.getElementById('stats-monthly-total').textContent = `${stats.total_monthly.toFixed(2)} ₽`;
    document.getElementById('stats-yearly-total').textContent = `${stats.total_yearly.toFixed(2)} ₽`;
    document.getElementById('stats-total-count').textContent = stats.total_subscriptions;
}

function createMonthlyExpensesChart(subs) {
    const ctx = document.getElementById('monthlyExpensesChart');

    // Destroy previous chart if exists
    if (monthlyExpensesChart) {
        monthlyExpensesChart.destroy();
    }

    // Calculate expenses for next 6 months
    const months = [];
    const expenses = [];
    const today = new Date();

    for (let i = 0; i < 6; i++) {
        const month = new Date(today.getFullYear(), today.getMonth() + i, 1);
        const monthName = month.toLocaleDateString('ru-RU', { month: 'short', year: 'numeric' });
        months.push(monthName);

        // Calculate expenses for this month
        let monthExpense = 0;
        subs.forEach(sub => {
            const monthlyAmount = calculateMonthlyAmount(sub);
            monthExpense += monthlyAmount;
        });
        expenses.push(monthExpense.toFixed(2));
    }

    monthlyExpensesChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: months,
            datasets: [{
                label: 'Расходы (₽)',
                data: expenses,
                borderColor: 'rgb(99, 102, 241)',
                backgroundColor: 'rgba(99, 102, 241, 0.1)',
                tension: 0.4,
                fill: true
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    display: true,
                    position: 'top'
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: {
                        callback: function(value) {
                            return value + ' ₽';
                        }
                    }
                }
            }
        }
    });
}

function createCategoryChart(subs) {
    const ctx = document.getElementById('categoryChart');

    // Destroy previous chart if exists
    if (categoryChart) {
        categoryChart.destroy();
    }

    // Group by category
    const categoryExpenses = {};
    subs.forEach(sub => {
        const category = sub.category || 'other';
        const monthlyAmount = calculateMonthlyAmount(sub);

        if (!categoryExpenses[category]) {
            categoryExpenses[category] = 0;
        }
        categoryExpenses[category] += monthlyAmount;
    });

    const categories = Object.keys(categoryExpenses);
    const amounts = Object.values(categoryExpenses);

    // Category names in Russian
    const categoryNames = {
        'streaming': 'Стриминг',
        'music': 'Музыка',
        'software': 'Софт',
        'gaming': 'Игры',
        'education': 'Образование',
        'cloud': 'Облако',
        'marketplace': 'Маркетплейсы',
        'crm': 'CRM',
        'other': 'Другое'
    };

    const colors = [
        'rgb(239, 68, 68)',
        'rgb(139, 92, 246)',
        'rgb(59, 130, 246)',
        'rgb(34, 197, 94)',
        'rgb(245, 158, 11)',
        'rgb(14, 165, 233)',
        'rgb(107, 114, 128)'
    ];

    categoryChart = new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: categories.map(c => categoryNames[c] || c),
            datasets: [{
                data: amounts,
                backgroundColor: colors.slice(0, categories.length),
                borderWidth: 2,
                borderColor: '#fff'
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    position: 'bottom'
                },
                tooltip: {
                    callbacks: {
                        label: function(context) {
                            return context.label + ': ' + context.parsed.toFixed(2) + ' ₽/мес';
                        }
                    }
                }
            }
        }
    });
}

function createTopSubscriptionsChart(subs) {
    const ctx = document.getElementById('topSubscriptionsChart');

    // Destroy previous chart if exists
    if (topSubscriptionsChart) {
        topSubscriptionsChart.destroy();
    }

    // Sort by monthly amount and take top 5
    const sorted = [...subs].sort((a, b) => {
        return calculateMonthlyAmount(b) - calculateMonthlyAmount(a);
    }).slice(0, 5);

    const names = sorted.map(s => s.name);
    const amounts = sorted.map(s => calculateMonthlyAmount(s));

    topSubscriptionsChart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: names,
            datasets: [{
                label: 'Расходы в месяц (₽)',
                data: amounts,
                backgroundColor: 'rgba(139, 92, 246, 0.8)',
                borderColor: 'rgb(139, 92, 246)',
                borderWidth: 1
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    display: false
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: {
                        callback: function(value) {
                            return value + ' ₽';
                        }
                    }
                }
            }
        }
    });
}

function createUpcomingPaymentsChart(subs) {
    const ctx = document.getElementById('upcomingPaymentsChart');

    // Destroy previous chart if exists
    if (upcomingPaymentsChart) {
        upcomingPaymentsChart.destroy();
    }

    // Get upcoming payments for next 30 days
    const today = new Date();
    const next30Days = new Date(today.getTime() + 30 * 24 * 60 * 60 * 1000);

    const upcomingPayments = [];
    subs.forEach(sub => {
        const nextDate = new Date(sub.next_billing_date);
        if (nextDate >= today && nextDate <= next30Days) {
            upcomingPayments.push({
                name: sub.name,
                date: nextDate,
                amount: sub.amount
            });
        }
    });

    // Sort by date
    upcomingPayments.sort((a, b) => a.date - b.date);

    // Take first 10
    const payments = upcomingPayments.slice(0, 10);
    const labels = payments.map(p => {
        return p.name + ' (' + p.date.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' }) + ')';
    });
    const amounts = payments.map(p => p.amount);

    upcomingPaymentsChart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [{
                label: 'Сумма платежа (₽)',
                data: amounts,
                backgroundColor: 'rgba(99, 102, 241, 0.8)',
                borderColor: 'rgb(99, 102, 241)',
                borderWidth: 1
            }]
        },
        options: {
            indexAxis: 'y',
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    display: false
                }
            },
            scales: {
                x: {
                    beginAtZero: true,
                    ticks: {
                        callback: function(value) {
                            return value + ' ₽';
                        }
                    }
                }
            }
        }
    });
}

// Inline Charts for Main Dashboard
function createInlineCharts() {
    if (subscriptions.length === 0) return;

    createInlineMonthlyChart();
    createInlineCategoryChart();
}

function createInlineMonthlyChart() {
    const ctx = document.getElementById('inlineMonthlyChart');
    if (!ctx) return;

    // Destroy previous chart if exists
    if (inlineMonthlyChart) {
        inlineMonthlyChart.destroy();
    }

    // Calculate for 3 months past, current, and 3 months future (7 total)
    const months = [];
    const pastExpenses = [];
    const futureExpenses = [];
    const today = new Date();
    const currentMonth = today.getMonth();
    const currentYear = today.getFullYear();

    // Helper function to calculate actual payments in a specific month
    function calculatePaymentsInMonth(targetYear, targetMonth) {
        let total = 0;
        const monthStart = new Date(targetYear, targetMonth, 1);
        const monthEnd = new Date(targetYear, targetMonth + 1, 0, 23, 59, 59);

        subscriptions.forEach(sub => {
            // Parse the next billing date
            let nextBillingDate = new Date(sub.next_billing_date);

            // Determine billing period in days
            let periodDays = 30; // default
            if (sub.billing_period === 'monthly') {
                periodDays = 30;
            } else if (sub.billing_period === 'yearly') {
                periodDays = 365;
            } else if (sub.billing_period === 'weekly') {
                periodDays = 7;
            } else if (sub.billing_period === 'custom' && sub.custom_period_days) {
                periodDays = sub.custom_period_days;
            }

            // Count how many payments occur in this month
            // We need to check both backwards (for past months) and forwards
            let paymentCount = 0;
            const maxIterations = 100; // safety limit

            // Start from a date well before the target month
            let checkDate = new Date(nextBillingDate);
            // Go back until we're before the month start
            while (checkDate > monthStart && paymentCount < maxIterations) {
                checkDate = new Date(checkDate.getTime() - periodDays * 24 * 60 * 60 * 1000);
            }

            // Now go forward and count payments in this month
            let iterations = 0;
            while (checkDate <= monthEnd && iterations < maxIterations) {
                if (checkDate >= monthStart && checkDate <= monthEnd) {
                    paymentCount++;
                }
                checkDate = new Date(checkDate.getTime() + periodDays * 24 * 60 * 60 * 1000);
                iterations++;
            }

            total += sub.amount * paymentCount;
        });

        return total;
    }

    // Create 7 months: -3, -2, -1, 0 (current), +1, +2, +3
    for (let i = -3; i <= 3; i++) {
        const date = new Date(currentYear, currentMonth + i, 1);
        const monthName = date.toLocaleDateString('ru-RU', { month: 'short', year: 'numeric' });
        months.push(monthName);

        // Calculate actual expenses for this specific month
        const monthExpense = calculatePaymentsInMonth(date.getFullYear(), date.getMonth());

        // Past and current months have actual data, future months have projected data
        if (i <= 0) {
            pastExpenses.push(monthExpense);
            futureExpenses.push(null);
        } else {
            pastExpenses.push(null);
            futureExpenses.push(monthExpense);
        }
    }

    inlineMonthlyChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: months,
            datasets: [
                {
                    label: 'Прошлые расходы',
                    data: pastExpenses,
                    borderColor: 'rgb(99, 102, 241)',
                    backgroundColor: 'rgba(99, 102, 241, 0.1)',
                    borderWidth: 3,
                    fill: true,
                    tension: 0.4,
                    pointRadius: 5,
                    pointBackgroundColor: 'rgb(99, 102, 241)'
                },
                {
                    label: 'Прогноз',
                    data: futureExpenses,
                    borderColor: 'rgb(139, 92, 246)',
                    backgroundColor: 'rgba(139, 92, 246, 0.05)',
                    borderWidth: 3,
                    borderDash: [8, 4],
                    fill: true,
                    tension: 0.4,
                    pointRadius: 5,
                    pointBackgroundColor: 'rgb(139, 92, 246)',
                    pointBorderColor: 'rgb(139, 92, 246)',
                    pointStyle: 'circle'
                }
            ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    display: true,
                    position: 'top',
                    labels: {
                        padding: 10,
                        font: {
                            size: 11
                        },
                        usePointStyle: true
                    }
                },
                tooltip: {
                    callbacks: {
                        label: function(context) {
                            return context.dataset.label + ': ' + context.parsed.y.toFixed(2) + ' ₽';
                        }
                    }
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: {
                        callback: function(value) {
                            return value + ' ₽';
                        }
                    }
                }
            }
        }
    });
}

function createInlineCategoryChart() {
    const ctx = document.getElementById('inlineCategoryChart');
    if (!ctx) return;

    // Destroy previous chart if exists
    if (inlineCategoryChart) {
        inlineCategoryChart.destroy();
    }

    // Calculate expenses by category
    const categoryExpenses = {};
    subscriptions.forEach(sub => {
        const category = sub.category || 'other';
        const monthlyAmount = calculateMonthlyAmount(sub);
        categoryExpenses[category] = (categoryExpenses[category] || 0) + monthlyAmount;
    });

    const categoryNames = {
        'streaming': 'Стриминг',
        'music': 'Музыка',
        'software': 'Софт',
        'gaming': 'Игры',
        'education': 'Образование',
        'cloud': 'Облако',
        'marketplace': 'Маркетплейсы',
        'crm': 'CRM',
        'other': 'Другое'
    };

    const labels = Object.keys(categoryExpenses).map(cat => categoryNames[cat] || cat);
    const data = Object.values(categoryExpenses);

    const colors = [
        'rgba(99, 102, 241, 0.8)',
        'rgba(139, 92, 246, 0.8)',
        'rgba(236, 72, 153, 0.8)',
        'rgba(251, 146, 60, 0.8)',
        'rgba(34, 197, 94, 0.8)',
        'rgba(59, 130, 246, 0.8)',
        'rgba(249, 115, 22, 0.8)',
        'rgba(168, 85, 247, 0.8)',
        'rgba(107, 114, 128, 0.8)'
    ];

    inlineCategoryChart = new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: labels,
            datasets: [{
                data: data,
                backgroundColor: colors,
                borderWidth: 2,
                borderColor: '#fff'
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    position: 'bottom',
                    labels: {
                        padding: 15,
                        font: {
                            size: 12
                        }
                    }
                }
            }
        }
    });
}

function fillStatsTable(subs) {
    const tbody = document.getElementById('stats-table-body');
    tbody.innerHTML = '';

    if (subs.length === 0) {
        tbody.innerHTML = '<tr><td colspan="7" style="text-align: center;">Нет подписок</td></tr>';
        return;
    }

    // Sort by monthly amount descending
    const sorted = [...subs].sort((a, b) => {
        return calculateMonthlyAmount(b) - calculateMonthlyAmount(a);
    });

    sorted.forEach(sub => {
        const monthlyAmount = calculateMonthlyAmount(sub);
        const yearlyAmount = monthlyAmount * 12;

        const category = sub.category || 'other';
        const categoryNames = {
            'streaming': 'Стриминг',
            'music': 'Музыка',
            'software': 'Софт',
            'gaming': 'Игры',
            'education': 'Образование',
            'cloud': 'Облако',
            'marketplace': 'Маркетплейсы',
            'crm': 'CRM',
            'other': 'Другое'
        };

        const periodNames = {
            'monthly': 'Ежемесячно',
            'yearly': 'Ежегодно',
            'weekly': 'Еженедельно',
            'custom': 'Другой'
        };

        const row = document.createElement('tr');
        row.innerHTML = `
            <td><strong>${sub.name}</strong></td>
            <td><span class="category-badge category-${category}">${categoryNames[category]}</span></td>
            <td>${sub.amount.toFixed(2)} ${getCurrencySymbol(sub.currency)}</td>
            <td>${periodNames[sub.billing_period] || sub.billing_period}</td>
            <td><strong>${monthlyAmount.toFixed(2)} ₽</strong></td>
            <td>${yearlyAmount.toFixed(2)} ₽</td>
            <td>${new Date(sub.next_billing_date).toLocaleDateString('ru-RU')}</td>
        `;
        tbody.appendChild(row);
    });
}

function calculateMonthlyAmount(sub) {
    let amount = sub.amount;

    // Convert to RUB if needed (simplified - assume 1:1 for demo)
    if (sub.currency !== 'RUB') {
        // In real app, you'd use exchange rates
        if (sub.currency === 'USD') amount *= 90;
        if (sub.currency === 'EUR') amount *= 100;
    }

    // Convert to monthly
    switch (sub.billing_period) {
        case 'weekly':
            return amount * 4.33; // Average weeks per month
        case 'monthly':
            return amount;
        case 'yearly':
            return amount / 12;
        case 'custom':
            if (sub.custom_period_days) {
                return (amount / sub.custom_period_days) * 30;
            }
            return amount;
        default:
            return amount;
    }
}

function getCurrencySymbol(currency) {
    const symbols = {
        'RUB': '₽',
        'USD': '$',
        'EUR': '€'
    };
    return symbols[currency] || currency;
}

// ============================================
// SUBSCRIPTION SELECTION FUNCTIONS
// ============================================

// Clear readonly/disabled restrictions from form fields
function clearFieldRestrictions() {
    const nameField = document.getElementById('sub-name');
    const categoryField = document.getElementById('sub-category');

    // Remove readonly and disabled attributes
    nameField.removeAttribute('readonly');
    categoryField.removeAttribute('disabled');

    // Reset styles
    nameField.style.backgroundColor = '';
    nameField.style.cursor = '';
    categoryField.style.backgroundColor = '';
    categoryField.style.cursor = '';
}

// Handle mode switching (from list / manual)
function handleModeSwitch(e) {
    const mode = e.target.dataset.mode;

    // Update mode buttons
    document.querySelectorAll('.mode-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    e.target.classList.add('active');

    // Update selection modes
    document.querySelectorAll('.selection-mode').forEach(section => {
        section.classList.remove('active');
    });

    if (mode === 'list') {
        document.getElementById('from-list-mode').classList.add('active');
        renderSubscriptionTemplates();
    } else {
        document.getElementById('manual-mode').classList.add('active');
        // Allow editing when switching to manual mode
        clearFieldRestrictions();
    }
}

// Render subscription templates
function renderSubscriptionTemplates() {
    // Reset filter and search
    currentCategoryFilter = 'all';
    const searchInput = document.getElementById('subscription-search');
    if (searchInput) searchInput.value = '';

    // Reset filter buttons
    document.querySelectorAll('.filter-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    const allBtn = document.querySelector('.filter-btn[data-category="all"]');
    if (allBtn) allBtn.classList.add('active');

    renderPopularSubscriptions();
    renderAllSubscriptions();
}

// Helper function to render subscription logo
function renderSubscriptionLogo(template) {
    if (template.logoUrl) {
        return `<img src="${template.logoUrl}" alt="${template.name}" onerror="this.style.display='none'; this.nextElementSibling.style.display='flex';">
                <div class="subscription-logo-fallback" style="display:none;">${template.logo}</div>`;
    }
    return template.logo;
}

// Render popular subscriptions (all popular items)
function renderPopularSubscriptions() {
    const container = document.getElementById('popular-subscriptions');
    const popular = subscriptionTemplates.filter(t => t.popular);

    container.innerHTML = popular.map(template => `
        <div class="subscription-card-item" data-subscription='${JSON.stringify(template)}'>
            <div class="subscription-logo">${renderSubscriptionLogo(template)}</div>
            <div class="subscription-card-name">${template.name}</div>
            <div class="subscription-card-price">от ${template.price} ${getCurrencySymbol(template.currency)}</div>
        </div>
    `).join('');

    // Add click handlers
    container.querySelectorAll('.subscription-card-item').forEach(card => {
        card.addEventListener('click', handleSubscriptionSelect);
    });
}

// Render all subscriptions list
function renderAllSubscriptions(searchFilter = '') {
    const container = document.getElementById('all-subscriptions-list');

    // Filter by search text
    let templates = subscriptionTemplates;
    if (searchFilter) {
        const lowerFilter = searchFilter.toLowerCase();
        templates = templates.filter(t =>
            t.name.toLowerCase().includes(lowerFilter)
        );
    }

    // Filter by category
    if (currentCategoryFilter !== 'all') {
        templates = templates.filter(t => t.category === currentCategoryFilter);
    }

    // Sort alphabetically
    templates = [...templates].sort((a, b) => a.name.localeCompare(b.name));

    // Show message if no results
    if (templates.length === 0) {
        container.innerHTML = '<div style="text-align: center; padding: 40px; color: var(--text-secondary);">Подписки не найдены</div>';
        return;
    }

    container.innerHTML = templates.map(template => `
        <div class="subscription-list-item" data-subscription='${JSON.stringify(template)}'>
            <div class="subscription-list-logo">${template.logo}</div>
            <div class="subscription-list-info">
                <div class="subscription-list-name">${template.name}</div>
                <div class="subscription-list-price">от ${template.price} ${getCurrencySymbol(template.currency)}/мес</div>
            </div>
        </div>
    `).join('');

    // Add click handlers
    container.querySelectorAll('.subscription-list-item').forEach(item => {
        item.addEventListener('click', handleSubscriptionSelect);
    });
}

// Handle subscription search
function handleSubscriptionSearch(e) {
    const query = e.target.value;
    renderAllSubscriptions(query);
}

// Handle category filter
function handleCategoryFilter(e) {
    const category = e.target.dataset.category;

    // Update active state
    document.querySelectorAll('.filter-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    e.target.classList.add('active');

    // Set current filter
    currentCategoryFilter = category;

    // Re-render with current search query
    const searchQuery = document.getElementById('subscription-search')?.value || '';
    renderAllSubscriptions(searchQuery);
}

// Handle subscription selection from template
function handleSubscriptionSelect(e) {
    const card = e.currentTarget;
    const templateData = JSON.parse(card.dataset.subscription);

    // Remove previous selection
    document.querySelectorAll('.subscription-card-item, .subscription-list-item').forEach(item => {
        item.classList.remove('selected');
    });

    // Mark as selected
    card.classList.add('selected');

    // Store selected template
    selectedSubscriptionTemplate = templateData;

    // Fill form and switch to manual mode
    fillFormWithTemplate(templateData);

    // Switch to manual mode
    document.querySelectorAll('.mode-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    document.querySelector('.mode-btn[data-mode="manual"]').classList.add('active');

    document.querySelectorAll('.selection-mode').forEach(section => {
        section.classList.remove('active');
    });
    document.getElementById('manual-mode').classList.add('active');
}

// Fill form with template data
function fillFormWithTemplate(template) {
    const nameField = document.getElementById('sub-name');
    const categoryField = document.getElementById('sub-category');

    // Fill all fields
    nameField.value = template.name;
    document.getElementById('sub-amount').value = template.price;
    document.getElementById('sub-currency').value = template.currency;
    document.getElementById('sub-period').value = template.period;
    categoryField.value = template.category;

    // Make name and category readonly (can't be edited)
    nameField.setAttribute('readonly', 'readonly');
    categoryField.setAttribute('disabled', 'disabled');
    nameField.style.backgroundColor = '#f5f5f5';
    nameField.style.cursor = 'not-allowed';
    categoryField.style.backgroundColor = '#f5f5f5';
    categoryField.style.cursor = 'not-allowed';

    // Set start date to today
    const today = new Date().toISOString().split('T')[0];
    document.getElementById('sub-start-date').value = today;
}

// Theme management
function initTheme() {
    const savedTheme = localStorage.getItem('theme') || 'light';
    document.documentElement.setAttribute('data-theme', savedTheme);
    updateThemeIcon(savedTheme);
}

function toggleTheme() {
    const currentTheme = document.documentElement.getAttribute('data-theme') || 'light';
    const newTheme = currentTheme === 'light' ? 'dark' : 'light';

    document.documentElement.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);
    updateThemeIcon(newTheme);
}

function updateThemeIcon(theme) {
    const themeBtn = document.getElementById('theme-toggle-btn');
    if (themeBtn) {
        themeBtn.textContent = theme === 'light' ? '🌙' : '☀️';
        themeBtn.title = theme === 'light' ? 'Темная тема' : 'Светлая тема';
    }
}

// User subscriptions filtering
function handleUserSubscriptionSearch(e) {
    userSubscriptionsSearchQuery = e.target.value.toLowerCase();
    renderSubscriptions();
}

function handleUserCategoryFilter(e) {
    const category = e.target.dataset.category;

    // Update active state
    document.querySelectorAll('.user-filter-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    e.target.classList.add('active');

    userSubscriptionsFilter = category;
    renderSubscriptions();
}

// Cancellation modal functions
async function cancelSubscriptionModal(subId, subName) {
    try {
        // Get cancellation instruction for this service
        const response = await fetch(`/api/cancellation/instructions/${encodeURIComponent(subName)}`, {
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
        });

        const modal = document.getElementById('cancellation-modal');
        const modalTitle = document.getElementById('cancellation-modal-title');
        const modalInstructions = document.getElementById('cancellation-modal-instructions');
        const modalNote = document.getElementById('cancellation-modal-note');
        const modalLink = document.getElementById('cancellation-modal-link');
        const confirmBtn = document.getElementById('confirm-cancellation-btn');

        modalTitle.textContent = `Отмена подписки: ${subName}`;

        if (response.ok) {
            const data = await response.json();
            const instruction = data.data || data;

            // Show instructions
            modalInstructions.innerHTML = `
                <h4>Инструкция по отмене:</h4>
                <p>${instruction.instructions}</p>
            `;

            // Show note if available
            if (instruction.note) {
                modalNote.innerHTML = `
                    <div class="info-box">
                        <strong>ℹ️ Примечание:</strong><br>
                        ${instruction.note}
                    </div>
                `;
                modalNote.style.display = 'block';
            } else {
                modalNote.style.display = 'none';
            }

            // Show link if available
            if (instruction.url) {
                modalLink.innerHTML = `
                    <a href="${instruction.url}" target="_blank" class="btn btn-primary btn-large">
                        Перейти к отмене ↗
                    </a>
                `;
                modalLink.style.display = 'block';
            } else {
                modalLink.style.display = 'none';
            }

            // Update confirm button
            confirmBtn.onclick = () => confirmCancellation(subId);
        } else {
            // No instructions found - show generic message
            modalInstructions.innerHTML = `
                <h4>Инструкция по отмене:</h4>
                <p>К сожалению, у нас пока нет инструкций для этого сервиса.
                   Обычно отмену подписки можно найти в настройках вашего аккаунта
                   или в разделе "Подписки" / "Биллинг".</p>
            `;
            modalNote.innerHTML = `
                <div class="info-box">
                    <strong>ℹ️ Примечание:</strong><br>
                    Для подписок, оформленных через App Store или Google Play,
                    отмена производится только в соответствующем магазине приложений.
                </div>
            `;
            modalNote.style.display = 'block';
            modalLink.style.display = 'none';

            // Update confirm button
            confirmBtn.onclick = () => confirmCancellation(subId);
        }

        // Show modal
        modal.classList.add('active');
        document.body.style.overflow = 'hidden';

    } catch (error) {
        console.error('Error loading cancellation instructions:', error);
        showNotification('Ошибка при загрузке инструкций', 'error');
    }
}

async function confirmCancellation(subId) {
    try {
        const response = await apiRequest(`/subscriptions/${subId}/cancel`, 'POST');

        if (response.success) {
            showNotification('Подписка помечена как отмененная', 'success');
            closeCancellationModal();
            await loadSubscriptions();
            await loadStats();
        } else {
            showNotification(response.error || 'Не удалось отменить подписку', 'error');
        }
    } catch (error) {
        console.error('Error cancelling subscription:', error);
        showNotification('Ошибка при отмене подписки', 'error');
    }
}

function closeCancellationModal() {
    const modal = document.getElementById('cancellation-modal');
    modal.classList.remove('active');
    document.body.style.overflow = '';
}

// Initialize theme on app load
document.addEventListener('DOMContentLoaded', () => {
    initTheme();
});
