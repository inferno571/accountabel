/* === Accountabel AI Dashboard Logic === */

// Auth gate: redirect to home if not logged in
const userId = localStorage.getItem('accountabel_user_id');
if (!userId) {
    window.location.href = '/';
}
const userName = localStorage.getItem('accountabel_user_name') || 'User';
let dashboardData = null;

// ── Toast Notifications ──
function showToast(message, type = 'info') {
    const container = document.getElementById('toast-container');
    if (!container) return;
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    
    let icon = 'ℹ️';
    if (type === 'success') icon = '✅';
    if (type === 'error') icon = '⚠️';
    if (type === 'warning') icon = '⏳';
    
    toast.innerHTML = `
        <div class="toast-icon">${icon}</div>
        <div style="flex:1;">${message}</div>
        <button class="toast-dismiss" onclick="this.parentElement.remove()">&times;</button>
    `;
    
    container.appendChild(toast);
    
    // Auto remove after 5 seconds
    setTimeout(() => {
        if (toast.parentElement) {
            toast.classList.add('out');
            setTimeout(() => {
                if (toast.parentElement) toast.remove();
            }, 300);
        }
    }, 5000);
}

// ── Fetch & Render ──
async function fetchDashboard() {
    try {
        const res = await fetch('/api/recovery/dashboard', { headers: { 'X-User-Id': userId } });
        const data = await res.json();
        if (data.status === 'success') {
            dashboardData = data;
            renderAll(data);
        }
    } catch (e) { console.error('Dashboard fetch error:', e); }
}

function renderAll(d) {
    renderStreak(d.streak, d.profile);
    renderStats(d);
    renderHeatmap(d.check_ins, d.slips);
    renderUrgeLog(d.cravings);
    renderCravingChart(d.cravings);
    renderBadges(d.streak, d.check_ins, d.cravings);
    renderHabits(d.profile);
    renderToolkit(d.tools);
}

// ── Live Streak Timer ──
let streakInterval = null;
function renderStreak(streak, profile) {
    const current = (streak && streak.current_streak) || 0;
    const highest = (streak && streak.highest_streak) || 0;
    const lastCheckIn = (streak && streak.last_check_in) || 0;
    const profileCreateTime = (profile && profile.create_time) || 0;
    // Fallback: use localStorage signup time if profile create_time is missing
    const signupTime = profileCreateTime || parseInt(localStorage.getItem('accountabel_signup_time') || '0', 10);

    document.getElementById('highest-streak').textContent = highest;

    if (streakInterval) clearInterval(streakInterval);

    function updateTimer() {
        let baseSeconds;
        if (lastCheckIn > 0 && current > 0) {
            const lastDate = new Date(lastCheckIn * 1000);
            const streakStart = new Date(lastDate);
            streakStart.setDate(streakStart.getDate() - (current - 1));
            streakStart.setHours(0, 0, 0, 0);
            baseSeconds = Math.floor((Date.now() - streakStart.getTime()) / 1000);
        } else if (signupTime > 0) {
            baseSeconds = Math.floor((Date.now() / 1000) - signupTime);
        } else {
            baseSeconds = 0;
        }

        if (baseSeconds < 0) baseSeconds = 0;
        const days = Math.floor(baseSeconds / 86400);
        const hours = Math.floor((baseSeconds % 86400) / 3600);
        const minutes = Math.floor((baseSeconds % 3600) / 60);
        const seconds = baseSeconds % 60;

        const setVal = (id, v) => { const el = document.getElementById(id); if (el) el.textContent = String(v).padStart(2, '0'); };
        setVal('t-days', days);
        setVal('t-hours', hours);
        setVal('t-mins', minutes);
        setVal('t-secs', seconds);
    }
    updateTimer();
    streakInterval = setInterval(updateTimer, 1000);
}

// ── Stats Cards ──
function renderStats(d) {
    const checkins = d.check_ins || [];
    const cravings = d.cravings || [];
    const slips = d.slips || [];
    const streak = d.streak || {};

    document.getElementById('stat-streak').textContent = streak.current_streak || 0;

    // Total check-ins
    document.getElementById('stat-checkins').textContent = checkins.length;

    // Average craving
    if (cravings.length > 0) {
        const avg = (cravings.reduce((s, c) => s + c.intensity, 0) / cravings.length).toFixed(1);
        document.getElementById('stat-avg-craving').textContent = avg;
    } else {
        document.getElementById('stat-avg-craving').textContent = '--';
    }

    // Total slips
    document.getElementById('stat-slips').textContent = slips.length;
}

// ── Calendar Heatmap (GitHub-style) ──
function renderHeatmap(checkIns, slips) {
    const grid = document.getElementById('heatmap-grid');
    const monthsRow = document.getElementById('heatmap-months');
    if (!grid) return;
    grid.innerHTML = '';
    if (monthsRow) monthsRow.innerHTML = '';

    // Build date -> count map for last ~6 months
    const now = new Date();
    const startDate = new Date(now);
    startDate.setMonth(startDate.getMonth() - 5);
    startDate.setDate(1);
    // Align to nearest previous Sunday
    while (startDate.getDay() !== 0) startDate.setDate(startDate.getDate() - 1);

    const dayMap = {};
    const slipSet = new Set();

    (checkIns || []).forEach(ci => {
        const d = new Date(ci.timestamp * 1000);
        const key = d.toISOString().slice(0, 10);
        dayMap[key] = (dayMap[key] || 0) + 1;
    });

    (slips || []).forEach(s => {
        const d = new Date(s.slip_date * 1000);
        slipSet.add(d.toISOString().slice(0, 10));
    });

    // Generate cells
    const today = new Date();
    today.setHours(23, 59, 59, 999);
    const cursor = new Date(startDate);
    let lastMonth = -1;
    const months = [];

    while (cursor <= today) {
        const key = cursor.toISOString().slice(0, 10);
        const count = dayMap[key] || 0;
        const isSlip = slipSet.has(key);

        const cell = document.createElement('div');
        cell.className = 'heatmap-cell';
        cell.title = `${key}: ${count} check-in(s)${isSlip ? ' (slip)' : ''}`;

        if (isSlip) {
            cell.classList.add('slip');
        } else if (count >= 4) {
            cell.dataset.level = '4';
        } else if (count >= 3) {
            cell.dataset.level = '3';
        } else if (count >= 2) {
            cell.dataset.level = '2';
        } else if (count >= 1) {
            cell.dataset.level = '1';
        }

        grid.appendChild(cell);

        // Track months for labels
        const m = cursor.getMonth();
        if (m !== lastMonth && cursor.getDay() === 0) {
            months.push({ month: m, col: Math.floor((cursor - startDate) / (7 * 86400000)) });
            lastMonth = m;
        }

        cursor.setDate(cursor.getDate() + 1);
    }

    // Month labels
    if (monthsRow) {
        const monthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
        months.forEach(m => {
            const span = document.createElement('span');
            span.textContent = monthNames[m.month];
            monthsRow.appendChild(span);
        });
    }
}

// ── Urge Log List ──
function renderUrgeLog(cravings) {
    const list = document.getElementById('urge-list');
    if (!list) return;

    if (!cravings || cravings.length === 0) {
        list.innerHTML = '<div class="empty-state"><div class="icon">--</div><p>No urges logged yet. Use the bot or log one here.</p></div>';
        return;
    }

    list.innerHTML = cravings.slice(0, 20).map(c => {
        const level = c.intensity <= 3 ? 'low' : c.intensity <= 6 ? 'mid' : 'high';
        const action = c.action_taken === 'rode_it_out' ? 'resisted' : c.action_taken === 'used_coping_tool' ? 'coped' : 'slipped';
        const actionText = c.action_taken === 'rode_it_out' ? 'Resisted' : c.action_taken === 'used_coping_tool' ? 'Coped' : 'Slipped';
        const date = new Date(c.logged_at * 1000);
        const timeStr = date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' }) + ' ' + date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' });
        const ctx = c.trigger_context || 'No context';

        return `<div class="urge-item">
            <div class="urge-intensity ${level}">${c.intensity}</div>
            <div class="urge-details">
                <div class="context">${escHtml(ctx)}</div>
                <div class="meta">${timeStr}</div>
            </div>
            <span class="urge-action-badge ${action}">${actionText}</span>
        </div>`;
    }).join('');
}

function escHtml(t) {
    const d = document.createElement('div');
    d.textContent = t;
    return d.innerHTML;
}

// ── Craving Trend Chart (pure SVG) ──
function renderCravingChart(cravings) {
    const container = document.getElementById('craving-chart');
    if (!container) return;

    if (!cravings || cravings.length < 2) {
        container.innerHTML = '<div class="chart-empty">Need at least 2 data points</div>';
        return;
    }

    // Get last 30 entries, reversed to chronological
    const data = cravings.slice(0, 30).reverse();
    const w = container.clientWidth || 400;
    const h = 160;
    const padX = 30, padY = 20;
    const plotW = w - padX * 2;
    const plotH = h - padY * 2;

    const points = data.map((c, i) => ({
        x: padX + (i / (data.length - 1)) * plotW,
        y: padY + plotH - (c.intensity / 10) * plotH
    }));

    // Smooth path
    let pathD = `M ${points[0].x} ${points[0].y}`;
    for (let i = 1; i < points.length; i++) {
        const prev = points[i - 1];
        const curr = points[i];
        const cpx = (prev.x + curr.x) / 2;
        pathD += ` C ${cpx} ${prev.y}, ${cpx} ${curr.y}, ${curr.x} ${curr.y}`;
    }

    // Fill area
    const areaD = pathD + ` L ${points[points.length - 1].x} ${h - padY} L ${points[0].x} ${h - padY} Z`;

    // Y-axis labels
    const yLabels = [0, 5, 10].map(v => ({
        v,
        y: padY + plotH - (v / 10) * plotH
    }));

    container.innerHTML = `<svg width="${w}" height="${h}" viewBox="0 0 ${w} ${h}">
        <defs>
            <linearGradient id="areaGrad" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stop-color="var(--accent)" stop-opacity="0.3"/>
                <stop offset="100%" stop-color="var(--accent)" stop-opacity="0"/>
            </linearGradient>
        </defs>
        ${yLabels.map(l => `
            <line x1="${padX}" y1="${l.y}" x2="${w - padX}" y2="${l.y}" stroke="var(--border-subtle)" stroke-dasharray="4"/>
            <text x="${padX - 8}" y="${l.y + 4}" fill="var(--text-muted)" font-size="9" text-anchor="end">${l.v}</text>
        `).join('')}
        <path d="${areaD}" fill="url(#areaGrad)"/>
        <path d="${pathD}" fill="none" stroke="var(--accent)" stroke-width="2" stroke-linecap="round"/>
        ${points.map(p => `<circle cx="${p.x}" cy="${p.y}" r="3" fill="var(--accent)"/>`).join('')}
    </svg>`;
}

// ── Achievement Badges ──
function renderBadges(streak, checkIns, cravings) {
    const s = streak || {};
    const ciCount = (checkIns || []).length;
    const crCount = (cravings || []).length;
    const currentStreak = s.current_streak || 0;
    const highestStreak = s.highest_streak || 0;

    const badges = [
        { icon: '1', name: 'Day One', earned: highestStreak >= 1 },
        { icon: '7', name: 'One Week', earned: highestStreak >= 7 },
        { icon: '30', name: '30 Days', earned: highestStreak >= 30 },
        { icon: '90', name: '90 Days', earned: highestStreak >= 90 },
        { icon: 'C', name: 'First Check-in', earned: ciCount >= 1 },
        { icon: '10', name: '10 Check-ins', earned: ciCount >= 10 },
        { icon: 'U', name: 'Urge Logger', earned: crCount >= 1 },
        { icon: 'W', name: 'Warrior', earned: crCount >= 10 },
    ];

    const grid = document.getElementById('badges-grid');
    if (!grid) return;
    grid.innerHTML = badges.map(b => `
        <div class="badge-item ${b.earned ? 'earned' : ''}">
            <div class="badge-icon">${b.icon}</div>
            <div class="badge-name">${b.name}</div>
        </div>
    `).join('');
}

// ── Habits Tracked ──
function renderHabits(profile) {
    const container = document.getElementById('habits-list');
    if (!container) return;

    if (!profile || !profile.addictions) {
        container.innerHTML = '<span class="habit-chip"><span class="dot"></span>None set</span>';
        return;
    }

    const habits = profile.addictions.split(',').map(h => h.trim()).filter(Boolean);
    container.innerHTML = habits.map(h => `<span class="habit-chip"><span class="dot"></span>${escHtml(h)}</span>`).join('');
}

// ── Coping Toolkit ──
function renderToolkit(tools) {
    const container = document.getElementById('toolkit-list');
    if (!container) return;

    if (!tools || tools.length === 0) {
        container.innerHTML = '<div class="empty-state"><p>No tools added yet.</p></div>';
        return;
    }

    container.innerHTML = tools.map(t => `
        <div class="tool-item">
            <div class="tool-icon">*</div>
            <div class="tool-name">${escHtml(t.content)}</div>
        </div>
    `).join('');
}

// ── Modals ──
function openModal(id) { document.getElementById(id).classList.add('open'); }
function closeModal(id) { document.getElementById(id).classList.remove('open'); }
function closeAllModals() { document.querySelectorAll('.modal-overlay').forEach(m => m.classList.remove('open')); }

// ── Urge Range Display ──
function updateRangeDisplay() {
    const val = document.getElementById('urge-range').value;
    document.getElementById('urge-range-val').textContent = val;
    const label = val <= 3 ? 'Mild' : val <= 6 ? 'Moderate' : 'Intense';
    document.getElementById('urge-range-label').textContent = label;
}

// ── Submit Urge ──
async function submitUrge() {
    const intensity = parseInt(document.getElementById('urge-range').value);
    const context = document.getElementById('urge-context').value;
    const action = document.getElementById('urge-action').value;

    try {
        await fetch('/api/recovery/craving/log', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json', 'X-User-Id': userId },
            body: JSON.stringify({ intensity, trigger_context: context, action_taken: action, tags: '[]' })
        });
        closeAllModals();
        document.getElementById('urge-context').value = '';
        document.getElementById('urge-range').value = 5;
        updateRangeDisplay();
        fetchDashboard();
    } catch (e) { console.error(e); }
}

// ── Submit Slip ──
async function submitSlip() {
    const notes = document.getElementById('slip-notes').value;
    try {
        await fetch('/api/recovery/slip/log', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json', 'X-User-Id': userId },
            body: JSON.stringify({ notes })
        });
        closeAllModals();
        document.getElementById('slip-notes').value = '';
        fetchDashboard();
    } catch (e) { console.error(e); }
}

// ── Add Coping Tool ──
function addCopingTool() {
    openModal('tool-modal');
}

async function submitCopingTool() {
    const text = document.getElementById('tool-input').value.trim();
    if (!text) {
        showToast('Please enter a coping tool', 'warning');
        return;
    }
    
    try {
        await fetch('/api/recovery/tools/add', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json', 'X-User-Id': userId },
            body: JSON.stringify({ tool_type: 'custom_distraction', content: text })
        });
        document.getElementById('tool-input').value = '';
        closeAllModals();
        fetchDashboard();
        showToast('Tool added successfully!', 'success');
    } catch (e) { 
        console.error(e); 
        showToast('Failed to add tool', 'error');
    }
}

// ── Breathing Exercise ──
function startBreathing() {
    openModal('breathing-modal');
    runBreathingCycle();
}

function runBreathingCycle() {
    const phases = [
        { text: 'Breathe In', duration: 4000 },
        { text: 'Hold', duration: 7000 },
        { text: 'Breathe Out', duration: 8000 },
    ];
    let cycle = 0;
    const maxCycles = 3;
    const textEl = document.getElementById('breathing-text');
    const countEl = document.getElementById('breathing-count');
    const circleEl = document.getElementById('breathing-circle');
    if (!textEl) return;

    function runPhase(phaseIdx) {
        if (cycle >= maxCycles) {
            textEl.textContent = 'Well done';
            countEl.textContent = '';
            circleEl.style.transform = 'scale(1)';
            return;
        }
        const phase = phases[phaseIdx];
        textEl.textContent = phase.text;
        const totalSec = phase.duration / 1000;
        let sec = totalSec;

        if (phaseIdx === 0) circleEl.style.transform = 'scale(1.3)';
        else if (phaseIdx === 2) circleEl.style.transform = 'scale(0.8)';

        const iv = setInterval(() => {
            sec--;
            countEl.textContent = sec > 0 ? sec : '';
            if (sec <= 0) {
                clearInterval(iv);
                const next = phaseIdx + 1;
                if (next < phases.length) {
                    runPhase(next);
                } else {
                    cycle++;
                    runPhase(0);
                }
            }
        }, 1000);

        countEl.textContent = totalSec;
    }

    circleEl.style.transform = 'scale(1)';
    runPhase(0);
}

// ── Nav ──
function switchTab(tabName) {
    document.querySelectorAll('.nav-item').forEach(n => n.classList.remove('active'));
    document.querySelector(`.nav-item[data-tab="${tabName}"]`)?.classList.add('active');
}

function logout() {
    localStorage.removeItem('accountabel_user_id');
    localStorage.removeItem('accountabel_session_token');
    localStorage.removeItem('accountabel_user_name');
    localStorage.removeItem('accountabel_user_email');
    window.location.href = '/';
}

// ── Init ──
document.addEventListener('DOMContentLoaded', () => {
    // Set user name in sidebar
    const avatarEl = document.querySelector('.user-avatar');
    const nameEl = document.querySelector('.user-pill span');
    if (avatarEl && userName) avatarEl.textContent = userName.charAt(0).toUpperCase();
    if (nameEl) nameEl.textContent = userName;

    fetchDashboard();
    // Refresh every 60s
    setInterval(fetchDashboard, 60000);
});
