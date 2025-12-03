// Teams and Members Management

// State
let currentTeams = [];
let currentMembers = [];
let currentInvitations = [];
let selectedTeam = null;

// Load teams screen
async function loadTeamsScreen() {
    try {
        const response = await apiRequest('/teams', 'GET');
        // Ensure response is an array
        currentTeams = Array.isArray(response) ? response : [];
        renderTeams();
    } catch (error) {
        console.error('Error loading teams:', error);
        currentTeams = [];
        renderTeams();
        alert('Ошибка при загрузке команд. Возможно, требуется применить миграции базы данных.');
    }
}

// Render teams list
function renderTeams() {
    const container = document.getElementById('teams-list');
    if (!container) return;

    if (currentTeams.length === 0) {
        container.innerHTML = `
            <div class="empty-state">
                <div class="empty-icon">👥</div>
                <h3>Нет команд</h3>
                <p>Создайте первую команду для совместного управления подписками</p>
                <button class="btn btn-primary" onclick="showCreateTeamModal()">
                    <span>➕</span> Создать команду
                </button>
            </div>
        `;
        return;
    }

    const teamsHTML = currentTeams.map(team => `
        <div class="team-card" onclick="showTeamDetails(${team.id})">
            <div class="team-header">
                <div class="team-icon">👥</div>
                <div class="team-info">
                    <h3>${team.name}</h3>
                    ${team.description ? `<p>${team.description}</p>` : ''}
                </div>
            </div>
            ${team.budget_amount ? `
                <div class="team-budget">
                    <span class="label">Бюджет:</span>
                    <span class="amount">${team.budget_amount} ${team.budget_currency}/${team.budget_period || 'month'}</span>
                </div>
            ` : ''}
            <div class="team-actions">
                <button class="btn btn-sm btn-secondary" onclick="event.stopPropagation(); editTeam(${team.id})">Редактировать</button>
                <button class="btn btn-sm btn-danger" onclick="event.stopPropagation(); deleteTeam(${team.id})">Удалить</button>
            </div>
        </div>
    `).join('');

    container.innerHTML = teamsHTML;
}

// Show create team modal
function showCreateTeamModal() {
    const modal = document.getElementById('team-modal');
    document.getElementById('team-modal-title').textContent = 'Создать команду';
    document.getElementById('team-form').reset();
    selectedTeam = null;
    modal.classList.add('active');
}

// Show edit team modal
async function editTeam(teamId) {
    try {
        const team = await apiRequest(`/teams/${teamId}`, 'GET');
        selectedTeam = team;

        const modal = document.getElementById('team-modal');
        document.getElementById('team-modal-title').textContent = 'Редактировать команду';

        document.getElementById('team-name').value = team.name;
        document.getElementById('team-description').value = team.description || '';
        document.getElementById('team-budget').value = team.budget_amount || '';
        document.getElementById('team-currency').value = team.budget_currency || 'RUB';
        document.getElementById('team-period').value = team.budget_period || 'monthly';

        modal.classList.add('active');
    } catch (error) {
        console.error('Error loading team:', error);
        alert('Ошибка при загрузке данных команды', 'error');
    }
}

// Save team
async function saveTeam(event) {
    event.preventDefault();

    const formData = {
        name: document.getElementById('team-name').value,
        description: document.getElementById('team-description').value || null,
        budget_amount: parseFloat(document.getElementById('team-budget').value) || null,
        budget_currency: document.getElementById('team-currency').value,
        budget_period: document.getElementById('team-period').value
    };

    try {
        if (selectedTeam) {
            // Update existing team
            await apiRequest(`/teams/${selectedTeam.id}`, 'PUT', formData);
            alert('Команда обновлена', 'success');
        } else {
            // Create new team
            await apiRequest('/teams', 'POST', formData);
            alert('Команда создана', 'success');
        }

        closeTeamModal();
        await loadTeamsScreen();
    } catch (error) {
        console.error('Error saving team:', error);
        alert(error.message || 'Ошибка при сохранении команды', 'error');
    }
}

// Delete team
async function deleteTeam(teamId) {
    if (!confirm('Вы уверены, что хотите удалить эту команду?')) {
        return;
    }

    try {
        await apiRequest(`/teams/${teamId}`, 'DELETE');
        alert('Команда удалена', 'success');
        await loadTeamsScreen();
    } catch (error) {
        console.error('Error deleting team:', error);
        alert('Ошибка при удалении команды', 'error');
    }
}

// Show team details with members
async function showTeamDetails(teamId) {
    try {
        const team = await apiRequest(`/teams/${teamId}`, 'GET');
        const members = await apiRequest(`/teams/${teamId}/members`, 'GET');

        selectedTeam = team;
        currentMembers = members;

        renderTeamDetails(team, members);

        const modal = document.getElementById('team-details-modal');
        modal.classList.add('active');
    } catch (error) {
        console.error('Error loading team details:', error);
        alert('Ошибка при загрузке данных команды', 'error');
    }
}

// Render team details
function renderTeamDetails(team, members) {
    const detailsContainer = document.getElementById('team-details-content');

    const membersHTML = members.length > 0 ? members.map(member => `
        <div class="member-item">
            <div class="member-info">
                <div class="member-icon">${member.role === 'admin' ? '👑' : '👤'}</div>
                <div>
                    <div class="member-name">${member.username}</div>
                    <div class="member-email">${member.email}</div>
                    <span class="member-role ${member.role}">${member.role === 'admin' ? 'Администратор' : 'Пользователь'}</span>
                </div>
            </div>
            ${team.owner_id !== member.id ? `
                <button class="btn btn-sm btn-danger" onclick="removeMember(${team.id}, ${member.id})">Удалить</button>
            ` : '<span class="owner-badge">Владелец</span>'}
        </div>
    `).join('') : '<p class="empty-message">Участников нет</p>';

    detailsContainer.innerHTML = `
        <div class="team-details-header">
            <h2>${team.name}</h2>
            ${team.description ? `<p>${team.description}</p>` : ''}
        </div>

        ${team.budget_amount ? `
            <div class="team-budget-info">
                <strong>Бюджет:</strong> ${team.budget_amount} ${team.budget_currency}/${team.budget_period || 'month'}
            </div>
        ` : ''}

        <div class="team-members-section">
            <div class="section-header">
                <h3>Участники (${members.length})</h3>
                <button class="btn btn-primary btn-sm" onclick="showInviteMemberModal(${team.id})">
                    <span>➕</span> Пригласить
                </button>
            </div>
            <div class="members-list">
                ${membersHTML}
            </div>
        </div>
    `;
}

// Show invite member modal
function showInviteMemberModal(teamId) {
    const modal = document.getElementById('invite-modal');
    document.getElementById('invite-team-id').value = teamId;
    document.getElementById('invite-form').reset();
    modal.classList.add('active');
}

// Send invitation
async function sendInvitation(event) {
    event.preventDefault();

    const teamId = parseInt(document.getElementById('invite-team-id').value);
    const formData = {
        email: document.getElementById('invite-email').value,
        role: document.getElementById('invite-role').value,
        team_id: teamId || null
    };

    try {
        const response = await apiRequest('/invitations', 'POST', formData);
        alert('Приглашение отправлено', 'success');
        console.log('Invitation URL:', response.url);
        closeInviteModal();
        await loadInvitations();
    } catch (error) {
        console.error('Error sending invitation:', error);
        alert(error.message || 'Ошибка при отправке приглашения', 'error');
    }
}

// Remove member from team
async function removeMember(teamId, memberId) {
    if (!confirm('Вы уверены, что хотите удалить этого участника из команды?')) {
        return;
    }

    try {
        await apiRequest(`/teams/${teamId}/members/${memberId}`, 'DELETE');
        alert('Участник удален из команды', 'success');
        await showTeamDetails(teamId);
    } catch (error) {
        console.error('Error removing member:', error);
        alert('Ошибка при удалении участника', 'error');
    }
}

// Load invitations
async function loadInvitations() {
    try {
        const invitations = await apiRequest('/invitations', 'GET');
        // Ensure response is an array
        currentInvitations = Array.isArray(invitations) ? invitations : [];
        renderInvitations();
    } catch (error) {
        console.error('Error loading invitations:', error);
        currentInvitations = [];
        renderInvitations();
        alert('Ошибка при загрузке приглашений. Возможно, требуется применить миграции базы данных.');
    }
}

// Render invitations
function renderInvitations() {
    const container = document.getElementById('invitations-list');
    if (!container) return;

    if (currentInvitations.length === 0) {
        container.innerHTML = '<p class="empty-message">Нет активных приглашений</p>';
        return;
    }

    const invitationsHTML = currentInvitations.map(inv => `
        <div class="invitation-item ${inv.status}">
            <div class="invitation-info">
                <div class="invitation-email">${inv.email}</div>
                <div class="invitation-details">
                    <span class="role ${inv.role}">${inv.role === 'admin' ? 'Администратор' : 'Пользователь'}</span>
                    ${inv.team_name ? `<span class="team">Команда: ${inv.team_name}</span>` : ''}
                    <span class="status">${inv.status === 'pending' ? 'Ожидает' : inv.status === 'accepted' ? 'Принято' : 'Истекло'}</span>
                </div>
                <div class="invitation-date">Отправлено: ${new Date(inv.created_at).toLocaleDateString('ru-RU')}</div>
            </div>
            ${inv.status === 'pending' ? `
                <button class="btn btn-sm btn-danger" onclick="deleteInvitation('${inv.token}')">Отменить</button>
            ` : ''}
        </div>
    `).join('');

    container.innerHTML = invitationsHTML;
}

// Delete invitation
async function deleteInvitation(token) {
    if (!confirm('Вы уверены, что хотите отменить это приглашение?')) {
        return;
    }

    try {
        await apiRequest(`/invitations/${token}`, 'DELETE');
        alert('Приглашение отменено', 'success');
        await loadInvitations();
    } catch (error) {
        console.error('Error deleting invitation:', error);
        alert('Ошибка при отмене приглашения', 'error');
    }
}

// Close modals
function closeTeamModal() {
    document.getElementById('team-modal').classList.remove('active');
}

function closeTeamDetailsModal() {
    document.getElementById('team-details-modal').classList.remove('active');
}

function closeInviteModal() {
    document.getElementById('invite-modal').classList.remove('active');
}

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    // Set up form listeners
    const teamForm = document.getElementById('team-form');
    if (teamForm) {
        teamForm.addEventListener('submit', saveTeam);
    }

    const inviteForm = document.getElementById('invite-form');
    if (inviteForm) {
        inviteForm.addEventListener('submit', sendInvitation);
    }
});
