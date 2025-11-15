// API Base URL
const API_URL = '/api';

// State
let currentUser = null;
let currentSubscription = null;
let subscriptions = [];

// Init
document.addEventListener('DOMContentLoaded', () => {
    initApp();
});

function initApp() {
    const token = localStorage.getItem('token');

    if (token) {
        loadApp();
    } else {
        showScreen('auth-screen');
    }

    setupEventListeners();
}

function setupEventListeners() {
    // Auth tabs
    document.querySelectorAll('.tab-btn').forEach(btn => {
        btn.addEventListener('click', (e) => {
            const tab = e.target.dataset.tab;
            switchTab(tab);
        });
    });

    // Auth forms
    document.getElementById('login-form').addEventListener('submit', handleLogin);
    document.getElementById('register-form').addEventListener('submit', handleRegister);

    // App buttons
    document.getElementById('logout-btn')?.addEventListener('click', handleLogout);
    document.getElementById('add-subscription-btn')?.addEventListener('click', () => openSubscriptionModal());
    document.getElementById('settings-btn')?.addEventListener('click', openSettingsModal);

    // Subscription form
    document.getElementById('subscription-form').addEventListener('submit', handleSaveSubscription);
    document.getElementById('sub-period').addEventListener('change', (e) => {
        const customDaysGroup = document.getElementById('custom-days-group');
        customDaysGroup.style.display = e.target.value === 'custom' ? 'block' : 'none';
    });

    // Telegram settings form
    document.getElementById('telegram-settings-form').addEventListener('submit', handleSaveTelegramSettings);

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

function switchTab(tab) {
    document.querySelectorAll('.tab-btn').forEach(btn => btn.classList.remove('active'));
    document.querySelectorAll('.auth-form').forEach(form => form.classList.remove('active'));

    document.querySelector(`[data-tab="${tab}"]`).classList.add('active');
    document.getElementById(`${tab}-form`).classList.add('active');
}

function showScreen(screenId) {
    document.querySelectorAll('.screen').forEach(screen => screen.classList.remove('active'));
    document.getElementById(screenId).classList.add('active');
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
    showScreen('auth-screen');
}

// App loading
async function loadApp() {
    console.log('loadApp() called');
    try {
        console.log('Switching to app-screen...');
        showScreen('app-screen');
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
            document.getElementById('total-subscriptions').textContent = stats.total_subscriptions;
            document.getElementById('monthly-total').textContent = `${stats.monthly_total.toFixed(2)} ₽`;
            document.getElementById('yearly-total').textContent = `${stats.yearly_total.toFixed(2)} ₽`;

            if (stats.next_payment) {
                const date = new Date(stats.next_payment.next_billing_date);
                const dateStr = date.toLocaleDateString('ru-RU');
                document.getElementById('next-payment').textContent = `${stats.next_payment.name} (${dateStr})`;
            } else {
                document.getElementById('next-payment').textContent = '-';
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
        }
    } catch (error) {
        console.error('Failed to load subscriptions:', error);
    }
}

function renderSubscriptions() {
    const list = document.getElementById('subscriptions-list');
    const emptyState = document.getElementById('empty-state');

    if (subscriptions.length === 0) {
        list.innerHTML = '';
        emptyState.classList.add('active');
        return;
    }

    emptyState.classList.remove('active');

    list.innerHTML = subscriptions.map(sub => {
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
                <div class="subscription-icon" style="${sub.color ? `background: ${sub.color}` : ''}">
                    ${icon}
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
                <div>
                    <div class="subscription-amount">
                        ${sub.amount.toFixed(2)} ${currencySymbols[sub.currency] || sub.currency}
                    </div>
                    <div class="subscription-next-date">
                        Списание: ${dateStr} ${daysUntil >= 0 ? `(через ${daysUntil} дн.)` : '(просрочено)'}
                    </div>
                    <div class="subscription-actions" style="margin-top: 12px;">
                        <button class="btn btn-secondary" onclick="editSubscription(${sub.id})">✏️</button>
                        <button class="btn btn-danger" onclick="deleteSubscription(${sub.id})">🗑️</button>
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

    form.reset();

    if (subscription) {
        title.textContent = 'Редактировать подписку';
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
        title.textContent = 'Добавить подписку';
        const today = new Date().toISOString().split('T')[0];
        document.getElementById('sub-start-date').value = today;
    }

    modal.classList.add('active');
}

async function handleSaveSubscription(e) {
    e.preventDefault();

    const name = document.getElementById('sub-name').value;
    const amount = parseFloat(document.getElementById('sub-amount').value);
    const currency = document.getElementById('sub-currency').value;
    const description = document.getElementById('sub-description').value || null;
    const billing_period = document.getElementById('sub-period').value;
    const start_date = document.getElementById('sub-start-date').value;
    const category = document.getElementById('sub-category').value || null;

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
            await loadStats();
            await loadSubscriptions();
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
            await loadStats();
            await loadSubscriptions();
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
