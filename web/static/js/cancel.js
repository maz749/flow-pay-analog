// API Base URL
const API_URL = '/api';

// State
let allServices = [];
let filteredServices = [];

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    setupCancelPage();
});

function setupCancelPage() {
    // Load all services
    loadServices();

    // Setup search
    const searchInput = document.getElementById('cancel-search');
    if (searchInput) {
        searchInput.addEventListener('input', handleSearch);
    }

    // Setup view all link
    const viewAllLink = document.getElementById('view-all-services');
    if (viewAllLink) {
        viewAllLink.addEventListener('click', (e) => {
            e.preventDefault();
            document.getElementById('cancel-search').value = '';
            loadServices();
        });
    }

    // Setup modal close
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

async function loadServices(query = '') {
    try {
        const url = query 
            ? `${API_URL}/cancellation-instructions?q=${encodeURIComponent(query)}&limit=100`
            : `${API_URL}/cancellation-instructions?limit=100`;
        
        const response = await fetch(url);
        const result = await response.json();
        
        // Handle both direct array and wrapped response
        const data = result.data || result;
        
        if (Array.isArray(data)) {
            allServices = data;
            filteredServices = data;
            renderServices();
        }
    } catch (error) {
        console.error('Failed to load services:', error);
        document.getElementById('services-list').innerHTML = 
            '<div class="error-message">Ошибка загрузки сервисов. Попробуйте позже.</div>';
    }
}

function handleSearch(e) {
    const query = e.target.value.trim();
    if (query.length >= 2) {
        loadServices(query);
    } else if (query.length === 0) {
        loadServices();
    }
}

function renderServices() {
    const container = document.getElementById('services-list');
    
    if (filteredServices.length === 0) {
        container.innerHTML = '<div class="no-results">Сервисы не найдены</div>';
        return;
    }

    container.innerHTML = filteredServices.map(service => `
        <div class="service-item" data-service-name="${service.service_name}">
            <div class="service-name">${service.service_name}</div>
            <div class="service-arrow">▼</div>
        </div>
    `).join('');

    // Add click handlers
    container.querySelectorAll('.service-item').forEach(item => {
        item.addEventListener('click', () => {
            const serviceName = item.dataset.serviceName;
            openServiceDetail(serviceName);
        });
    });
}

async function openServiceDetail(serviceName) {
    const modal = document.getElementById('service-detail-modal');
    const modalTitle = document.getElementById('modal-service-name');
    const modalInstructions = document.getElementById('modal-instructions');
    const modalNotes = document.getElementById('modal-notes');
    const modalCancelUrl = document.getElementById('modal-cancel-url');

    modalTitle.textContent = `Отмена подписки: ${serviceName}`;
    modalInstructions.innerHTML = '<div class="loading-spinner">Загрузка...</div>';
    modalNotes.innerHTML = '';
    modalCancelUrl.style.display = 'none';

    try {
        const response = await fetch(`${API_URL}/cancellation-instructions/${encodeURIComponent(serviceName)}`);
        const result = await response.json();
        
        // Handle both direct object and wrapped response
        const instruction = result.data || result;

        if (instruction && instruction.service_name) {
            
            // Render instructions
            modalInstructions.innerHTML = `
                <h3>Инструкция по отмене:</h3>
                <p>${instruction.instructions}</p>
            `;

            // Render notes if available
            if (instruction.notes) {
                modalNotes.innerHTML = `
                    <div class="notes-box">
                        <strong>Примечания:</strong>
                        <p>${instruction.notes}</p>
                    </div>
                `;
            }

            // Show cancellation URL if available
            if (instruction.cancellation_url) {
                modalCancelUrl.href = instruction.cancellation_url;
                modalCancelUrl.style.display = 'inline-block';
            }
        } else {
            modalInstructions.innerHTML = '<p>Инструкция не найдена</p>';
        }
    } catch (error) {
        console.error('Failed to load instruction:', error);
        modalInstructions.innerHTML = '<p>Ошибка загрузки инструкции</p>';
    }

    modal.classList.add('active');
}

