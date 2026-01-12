// app.js - 主应用逻辑

const token = localStorage.getItem('token');
const role = parseInt(localStorage.getItem('role') || '0');
const isAdmin = (role & 1) !== 0; // 1 = admin 权限位
const canStock = (role & 4) !== 0; // 4 = 库存管理权限位
const canApplyCoupon = (role & 8) !== 0; // 8 = 卡券申请权限位
const userName = localStorage.getItem('user_name') || '';
if (!token) window.location.href = '/static/html/login.html';

// 显示用户信息
document.getElementById('userName').textContent = userName;
document.getElementById('userAvatar').textContent = userName ? userName.charAt(0).toUpperCase() : '?';

// 用户管理相关
let currentPage = 1;
const pageSize = 10;
let totalUsers = 0;
let userList = [];
let roleList = []; // 系统角色列表

// 卡券管理相关
let couponPage = 1;
let totalCoupons = 0;
let couponList = [];
let couponTypeList = []; // 卡券类型列表
let currentPageName = 'users'; // 当前页面

// 我的卡券相关
let myCouponPage = 1;
let totalMyCoupons = 0;
let myCouponList = [];

// 请求封装
async function request(url, options = {}) {
    const currentToken = localStorage.getItem('token');
    if (!currentToken) {
        window.location.href = '/static/html/login.html';
        return;
    }
    const resp = await fetch(url, {
        ...options,
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + currentToken,
            ...(options.headers || {})
        }
    });
    const data = await resp.json();
    if (data.code === 401) {
        localStorage.removeItem('token');
        localStorage.removeItem('role');
        window.location.href = '/static/html/login.html';
    }
    return data;
}

// ========== 角色相关 ==========
async function loadRoles() {
    const data = await request('/api/v1/user/roles');
    if (data.code === 0) {
        roleList = data.data.list || [];
    }
}

function formatRoleTags(userRole) {
    const tags = roleList
        .filter(r => (userRole & r.Role) !== 0)
        .map(r => `<span class="perm-tag perm-role-${r.Role}">${r.Name}</span>`);
    return tags.length > 0 ? `<div class="permission-tags">${tags.join('')}</div>` : '-';
}

function renderRoleCheckboxes(userRole = 0) {
    const container = document.getElementById('roleGrid');
    container.innerHTML = roleList.map(r => `
        <label class="role-item">
            <input type="checkbox" name="roleCheck" value="${r.Role}" ${(userRole & r.Role) !== 0 ? 'checked' : ''}>
            <span class="role-label">${r.Name}</span>
        </label>
    `).join('');
}

// ========== 用户相关 ==========

async function loadUsers() {
    const data = await request(`/api/v1/user/list?page=${currentPage}&size=${pageSize}`);
    if (data.code !== 0) {
        toast(data.msg, 'error');
        return;
    }
    totalUsers = data.data.total;
    userList = data.data.list;
    const tbody = document.getElementById('userTable');
    tbody.innerHTML = userList.map(u => `
        <tr>
            <td data-label="ID">${u.id}</td>
            <td data-label="昵称">${u.name}</td>
            <td data-label="账号">${u.account}</td>
            <td data-label="权限">${formatRoleTags(u.role)}</td>
            <td data-label="创建时间">${formatTimestamp(u.created_at)}</td>
            <td class="actions">
                <button class="btn btn-primary btn-sm" onclick="showEditModal(${u.id})">编辑</button>
                <button class="btn btn-danger btn-sm" onclick="deleteUser(${u.id})">删除</button>
            </td>
        </tr>
    `).join('');

    document.getElementById('pageInfo').textContent = `第 ${currentPage} 页 / 共 ${Math.ceil(totalUsers/pageSize)} 页`;
    document.getElementById('prevBtn').disabled = currentPage <= 1;
    document.getElementById('nextBtn').disabled = currentPage >= Math.ceil(totalUsers/pageSize);
}

function changePage(delta) {
    currentPage += delta;
    loadUsers();
}

function showAddModal() {
    document.getElementById('modalTitle').textContent = '新增用户';
    document.getElementById('editId').value = '';
    document.getElementById('inputName').value = '';
    document.getElementById('inputAccount').value = '';
    document.getElementById('inputPassword').value = '';
    // 渲染权限复选框，默认勾选登录权限(2)
    renderRoleCheckboxes(2);
    document.getElementById('accountGroup').style.display = 'block';
    document.getElementById('roleGroup').style.display = 'block';
    document.getElementById('userModal').classList.add('show');
}

function showEditModal(id) {
    const user = userList.find(u => u.id === id);
    if (!user) return;
    document.getElementById('modalTitle').textContent = '编辑用户';
    document.getElementById('editId').value = id;
    document.getElementById('inputName').value = user.name;
    document.getElementById('inputAccount').value = user.account;
    document.getElementById('inputPassword').value = '';
    // 根据用户权限渲染复选框
    renderRoleCheckboxes(user.role);
    document.getElementById('accountGroup').style.display = 'block';
    document.getElementById('roleGroup').style.display = 'block';
    document.getElementById('userModal').classList.add('show');
}

function closeModal() {
    document.getElementById('userModal').classList.remove('show');
}

async function saveUser() {
    const id = document.getElementById('editId').value;
    const name = document.getElementById('inputName').value.trim();
    const account = document.getElementById('inputAccount').value.trim();
    const password = document.getElementById('inputPassword').value;

    // 收集所有选中的权限位
    let role = 0;
    document.querySelectorAll('input[name="roleCheck"]:checked').forEach(cb => {
        role |= parseInt(cb.value);
    });

    if (id) {
        const body = { id: parseInt(id) };
        if (name) body.name = name;
        if (account) body.account = account;
        if (password) body.password = password;
        body.role = role;

        showLoading();
        try {
            const data = await request('/api/v1/user/update', {
                method: 'PUT',
                body: JSON.stringify(body)
            });
            if (data.code === 0) {
                closeModal();
                toast('更新成功', 'success');
                await loadUsers();
            } else {
                toast(data.msg, 'error');
            }
        } finally {
            hideLoading();
        }
    } else {
        if (!name || !account || !password) {
            toast('请填写完整信息', 'warning');
            return;
        }
        showLoading();
        try {
            const data = await request('/api/v1/user/add', {
                method: 'POST',
                body: JSON.stringify({ name, account, password, role })
            });
            if (data.code === 0) {
                closeModal();
                toast('创建成功', 'success');
                await loadUsers();
            } else {
                toast(data.msg, 'error');
            }
        } finally {
            hideLoading();
        }
    }
}

async function deleteUser(id) {
    if (!confirm('确定要删除该用户吗？')) return;
    showLoading();
    try {
        const data = await request(`/api/v1/user/delete/${id}`, { method: 'DELETE' });
        if (data.code === 0) {
            toast('删除成功', 'success');
            await loadUsers();
        } else {
            toast(data.msg, 'error');
        }
    } finally {
        hideLoading();
    }
}

function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('role');
    localStorage.removeItem('user_name');
    window.location.href = '/static/html/login.html';
}

// ========== 用户下拉菜单 ==========
function toggleUserDropdown() {
    document.getElementById('userDropdown').classList.toggle('open');
}

// 点击其他地方关闭下拉菜单
document.addEventListener('click', function(e) {
    const dropdown = document.getElementById('userDropdown');
    if (!dropdown.contains(e.target)) {
        dropdown.classList.remove('open');
    }
});

// ========== 设置弹窗 ==========
async function showSettingsModal() {
    // 关闭下拉菜单
    document.getElementById('userDropdown').classList.remove('open');

    // 加载当前用户资料
    showLoading();
    try {
        const data = await request('/api/v1/user/profile');
        if (data.code === 0) {
            document.getElementById('settingsName').value = data.data.name || '';
            document.getElementById('settingsPassword').value = '';
            document.getElementById('settingsModal').classList.add('show');
        } else {
            toast(data.msg, 'error');
        }
    } finally {
        hideLoading();
    }
}

function closeSettingsModal() {
    document.getElementById('settingsModal').classList.remove('show');
}

async function saveSettings() {
    const name = document.getElementById('settingsName').value.trim();
    const password = document.getElementById('settingsPassword').value;

    const body = {};
    if (name) body.name = name;
    if (password) body.password = password;

    const btn = document.getElementById('saveSettingsBtn');
    setBtnLoading(btn, true);

    try {
        const data = await request('/api/v1/user/profile', {
            method: 'PUT',
            body: JSON.stringify(body)
        });

        if (data.code === 0) {
            closeSettingsModal();
            toast('保存成功', 'success');
            // 更新页面上显示的用户名
            if (name) {
                localStorage.setItem('user_name', name);
                document.getElementById('userName').textContent = name;
                document.getElementById('userAvatar').textContent = name.charAt(0).toUpperCase();
            }
        } else {
            toast(data.msg, 'error');
        }
    } finally {
        setBtnLoading(btn, false);
    }
}

// ========== 页面切换 ==========
function switchPage(pageName) {
    currentPageName = pageName;

    // 更新菜单激活状态
    document.querySelectorAll('.menu-item').forEach(item => {
        item.classList.remove('active');
        if (item.dataset.page === pageName) {
            item.classList.add('active');
        }
    });

    // 切换页面显示
    document.querySelectorAll('.content > .card').forEach(card => {
        card.classList.add('page-hidden');
    });
    const targetPage = document.getElementById(`page-${pageName}`);
    if (targetPage) {
        targetPage.classList.remove('page-hidden');
    }

    // 关闭移动端菜单
    closeMobileMenu();
}

// 绑定菜单点击事件
document.querySelectorAll('.menu-item[data-page]').forEach(item => {
    item.addEventListener('click', () => {
        const pageName = item.dataset.page;
        switchPage(pageName);

        // 首次切换到卡券页面时加载数据
        if (pageName === 'coupons' && couponList.length === 0) {
            loadCouponTypes().then(() => loadCoupons());
        }

        // 首次切换到我的卡券页面时加载数据
        if (pageName === 'my-coupons' && myCouponList.length === 0) {
            loadCouponTypes().then(() => {
                renderMyCouponTypeOptions();
                loadMyCoupons();
            });
        }
    });
});

// ========== 移动端菜单 ==========
function toggleMobileMenu() {
    const sidebar = document.getElementById('sidebar');
    const overlay = document.getElementById('mobileMenuOverlay');
    sidebar.classList.toggle('mobile-open');
    overlay.classList.toggle('show');
}

function closeMobileMenu() {
    const sidebar = document.getElementById('sidebar');
    const overlay = document.getElementById('mobileMenuOverlay');
    sidebar.classList.remove('mobile-open');
    overlay.classList.remove('show');
}

// ========== 卡券类型相关 ==========
async function loadCouponTypes() {
    const data = await request('/api/v1/coupon/types');
    if (data.code === 0) {
        couponTypeList = data.data.list || [];
        renderCouponTypeOptions();
    }
}

function renderCouponTypeOptions() {
    // 筛选下拉
    const filterSelect = document.getElementById('filterCouponType');
    filterSelect.innerHTML = '<option value="">全部</option>' +
        couponTypeList.map(t => `<option value="${t.type}">${t.name}</option>`).join('');

    // 添加/编辑卡券弹窗
    const inputSelect = document.getElementById('inputCouponType');
    inputSelect.innerHTML = couponTypeList.map(t =>
        `<option value="${t.type}">${t.name}</option>`
    ).join('');

    // 导入弹窗
    const importSelect = document.getElementById('importCouponType');
    importSelect.innerHTML = couponTypeList.map(t =>
        `<option value="${t.type}">${t.name}</option>`
    ).join('');
}

function getCouponTypeName(type) {
    const t = couponTypeList.find(ct => ct.type === type);
    return t ? t.name : '未知类型';
}

// ========== 卡券列表 ==========
async function loadCoupons() {
    const typeFilter = document.getElementById('filterCouponType').value;
    const takenFilter = document.getElementById('filterTakenStatus').value;

    let url = `/api/v1/coupon/list?page=${couponPage}&size=${pageSize}`;
    if (typeFilter) url += `&type=${typeFilter}`;
    if (takenFilter !== '') url += `&taken=${takenFilter}`;

    const data = await request(url);
    if (data.code !== 0) {
        toast(data.msg, 'error');
        return;
    }

    totalCoupons = data.data.total;
    couponList = data.data.list || [];

    renderCouponTable();
    renderCouponCards();
    updateCouponPagination();
}

function renderCouponTable() {
    const tbody = document.getElementById('couponTable');
    tbody.innerHTML = couponList.map(c => `
        <tr>
            <td data-label="ID">${c.id}</td>
            <td data-label="卡券码" class="coupon-code">${c.coupon}</td>
            <td data-label="类型"><span class="type-tag">${c.type_name}</span></td>
            <td data-label="状态">${c.is_taken ?
                '<span class="status-tag status-taken">已领取</span>' :
                '<span class="status-tag status-available">未领取</span>'}</td>
            <td data-label="领取人">${c.is_taken ? (c.taker_name || '-') : '-'}</td>
            <td data-label="创建时间">${formatTimestamp(c.created_at)}</td>
            <td class="actions">
                ${!c.is_taken ? `
                    <button class="btn btn-primary btn-sm" onclick="showEditCouponModal(${c.id})">编辑</button>
                    <button class="btn btn-danger btn-sm" onclick="deleteCoupon(${c.id})">删除</button>
                ` : '<span class="text-muted">-</span>'}
            </td>
        </tr>
    `).join('');
}

function renderCouponCards() {
    const container = document.getElementById('couponCardList');
    container.innerHTML = couponList.map(c => `
        <div class="coupon-card">
            <div class="coupon-card-header">
                <span class="coupon-card-id">#${c.id}</span>
                ${c.is_taken ?
                    '<span class="status-tag status-taken">已领取</span>' :
                    '<span class="status-tag status-available">未领取</span>'}
            </div>
            <div class="coupon-card-code">${c.coupon}</div>
            <div class="coupon-card-info">
                <span class="type-tag">${c.type_name}</span>
                <span class="coupon-card-time">${formatTimestamp(c.created_at)}</span>
            </div>
            ${c.is_taken && c.taker_name ? `
                <div class="coupon-card-taker">
                    <span class="taker-label">领取人：</span>
                    <span class="taker-name">${c.taker_name}</span>
                </div>
            ` : ''}
            ${!c.is_taken ? `
                <div class="coupon-card-actions">
                    <button class="btn btn-primary btn-sm" onclick="showEditCouponModal(${c.id})">编辑</button>
                    <button class="btn btn-danger btn-sm" onclick="deleteCoupon(${c.id})">删除</button>
                </div>
            ` : ''}
        </div>
    `).join('');
}

function updateCouponPagination() {
    const totalPages = Math.ceil(totalCoupons / pageSize);
    document.getElementById('couponPageInfo').textContent = `第 ${couponPage} 页 / 共 ${totalPages} 页`;
    document.getElementById('couponPrevBtn').disabled = couponPage <= 1;
    document.getElementById('couponNextBtn').disabled = couponPage >= totalPages;
}

function changeCouponPage(delta) {
    couponPage += delta;
    loadCoupons();
}

// ========== 添加/编辑卡券弹窗 ==========
function showAddCouponModal() {
    document.getElementById('couponModalTitle').textContent = '添加卡券';
    document.getElementById('couponEditId').value = '';
    document.getElementById('inputCouponCode').value = '';
    if (couponTypeList.length > 0) {
        document.getElementById('inputCouponType').value = couponTypeList[0].type;
    }
    document.getElementById('couponModal').classList.add('show');
}

function showEditCouponModal(id) {
    const coupon = couponList.find(c => c.id === id);
    if (!coupon) return;

    document.getElementById('couponModalTitle').textContent = '编辑卡券';
    document.getElementById('couponEditId').value = id;
    document.getElementById('inputCouponCode').value = coupon.coupon;
    document.getElementById('inputCouponType').value = coupon.type;
    document.getElementById('couponModal').classList.add('show');
}

function closeCouponModal() {
    document.getElementById('couponModal').classList.remove('show');
}

async function saveCoupon() {
    const id = document.getElementById('couponEditId').value;
    const coupon = document.getElementById('inputCouponCode').value.trim();
    const type = parseInt(document.getElementById('inputCouponType').value);

    if (!coupon) {
        toast('请输入卡券码', 'warning');
        return;
    }

    showLoading();
    try {
        if (id) {
            // 编辑
            const data = await request('/api/v1/coupon/update', {
                method: 'PUT',
                body: JSON.stringify({ id: parseInt(id), coupon, type })
            });
            if (data.code === 0) {
                closeCouponModal();
                toast('更新成功', 'success');
                await loadCoupons();
            } else {
                toast(data.msg, 'error');
            }
        } else {
            // 新增
            const data = await request('/api/v1/coupon/add', {
                method: 'POST',
                body: JSON.stringify({ coupon, type })
            });
            if (data.code === 0) {
                closeCouponModal();
                toast('添加成功', 'success');
                await loadCoupons();
            } else {
                toast(data.msg, 'error');
            }
        }
    } finally {
        hideLoading();
    }
}

async function deleteCoupon(id) {
    if (!confirm('确定要删除该卡券吗？')) return;

    showLoading();
    try {
        const data = await request(`/api/v1/coupon/delete/${id}`, { method: 'DELETE' });
        if (data.code === 0) {
            toast('删除成功', 'success');
            await loadCoupons();
        } else {
            toast(data.msg, 'error');
        }
    } finally {
        hideLoading();
    }
}

// ========== 批量导入弹窗 ==========
function showImportModal() {
    document.getElementById('importCoupons').value = '';
    if (couponTypeList.length > 0) {
        document.getElementById('importCouponType').value = couponTypeList[0].type;
    }
    document.getElementById('importModal').classList.add('show');
}

function closeImportModal() {
    document.getElementById('importModal').classList.remove('show');
}

async function importCoupons() {
    const coupons = document.getElementById('importCoupons').value.trim();
    const type = parseInt(document.getElementById('importCouponType').value);

    if (!coupons) {
        toast('请输入卡券码', 'warning');
        return;
    }

    const btn = document.getElementById('importBtn');
    setBtnLoading(btn, true);

    try {
        const data = await request('/api/v1/coupon/import', {
            method: 'POST',
            body: JSON.stringify({ coupons, type })
        });

        if (data.code === 0) {
            closeImportModal();
            showImportResult(data.data);
            await loadCoupons();
        } else {
            toast(data.msg, 'error');
        }
    } finally {
        setBtnLoading(btn, false);
    }
}

function showImportResult(result) {
    const resultDiv = document.getElementById('importResult');
    let html = `
        <div class="result-summary">
            <div class="result-item">
                <span class="result-label">总数:</span>
                <span class="result-value">${result.total}</span>
            </div>
            <div class="result-item success">
                <span class="result-label">成功:</span>
                <span class="result-value">${result.success}</span>
            </div>
            <div class="result-item ${result.failed > 0 ? 'warning' : ''}">
                <span class="result-label">重复:</span>
                <span class="result-value">${result.failed}</span>
            </div>
        </div>
    `;

    if (result.duplicates && result.duplicates.length > 0) {
        html += `
            <div class="result-duplicates">
                <div class="duplicates-title">重复的卡券:</div>
                <div class="duplicates-list">${result.duplicates.join(', ')}</div>
            </div>
        `;
    }

    resultDiv.innerHTML = html;
    document.getElementById('importResultModal').classList.add('show');
}

function closeImportResultModal() {
    document.getElementById('importResultModal').classList.remove('show');
}

// ========== 我的卡券 ==========
async function loadMyCoupons() {
    const typeFilter = document.getElementById('filterMyCouponType').value;

    let url = `/api/v1/my-coupon/list?page=${myCouponPage}&size=${pageSize}`;
    if (typeFilter) url += `&type=${typeFilter}`;

    const data = await request(url);
    if (data.code !== 0) {
        toast(data.msg, 'error');
        return;
    }

    totalMyCoupons = data.data.total;
    myCouponList = data.data.list || [];

    renderMyCouponTable();
    renderMyCouponCards();
    updateMyCouponPagination();
}

function renderMyCouponTable() {
    const tbody = document.getElementById('myCouponTable');
    tbody.innerHTML = myCouponList.map(c => `
        <tr>
            <td data-label="ID">${c.id}</td>
            <td data-label="类型"><span class="type-tag">${c.type_name}</span></td>
            <td data-label="领取时间">${formatTimestamp(c.taken_at)}</td>
            <td class="actions">
                <button class="btn btn-primary btn-sm" onclick="showMyCouponDetail(${c.id})">详情</button>
            </td>
        </tr>
    `).join('');
}

function renderMyCouponCards() {
    const container = document.getElementById('myCouponCardList');
    container.innerHTML = myCouponList.map(c => `
        <div class="my-coupon-card">
            <div class="my-coupon-card-header">
                <span class="my-coupon-card-id">#${c.id}</span>
                <span class="type-tag">${c.type_name}</span>
            </div>
            <div class="my-coupon-card-info">
                <span class="my-coupon-card-label">领取时间</span>
                <span class="my-coupon-card-time">${formatTimestamp(c.taken_at)}</span>
            </div>
            <div class="my-coupon-card-actions">
                <button class="btn btn-primary btn-sm" onclick="showMyCouponDetail(${c.id})">查看详情</button>
            </div>
        </div>
    `).join('');
}

function updateMyCouponPagination() {
    const totalPages = Math.ceil(totalMyCoupons / pageSize);
    document.getElementById('myCouponPageInfo').textContent = `第 ${myCouponPage} 页 / 共 ${totalPages || 1} 页`;
    document.getElementById('myCouponPrevBtn').disabled = myCouponPage <= 1;
    document.getElementById('myCouponNextBtn').disabled = myCouponPage >= totalPages;
}

function changeMyCouponPage(delta) {
    myCouponPage += delta;
    loadMyCoupons();
}

async function showMyCouponDetail(id) {
    showLoading();
    try {
        const data = await request(`/api/v1/my-coupon/detail/${id}`);
        if (data.code !== 0) {
            toast(data.msg, 'error');
            return;
        }

        const detail = data.data;
        document.getElementById('detailCouponType').textContent = detail.type_name;
        document.getElementById('detailCouponCode').textContent = detail.coupon;
        document.getElementById('detailTakenAt').textContent = formatTimestamp(detail.taken_at);
        document.getElementById('myCouponDetailModal').classList.add('show');
    } finally {
        hideLoading();
    }
}

function closeMyCouponDetailModal() {
    document.getElementById('myCouponDetailModal').classList.remove('show');
}

async function copyToClipboard(text) {
    try {
        await navigator.clipboard.writeText(text);
        toast('复制成功', 'success');
    } catch (err) {
        // 降级方案：使用传统方式复制
        const textArea = document.createElement('textarea');
        textArea.value = text;
        textArea.style.position = 'fixed';
        textArea.style.left = '-9999px';
        document.body.appendChild(textArea);
        textArea.select();
        try {
            document.execCommand('copy');
            toast('复制成功', 'success');
        } catch (e) {
            toast('复制失败，请手动复制', 'error');
        }
        document.body.removeChild(textArea);
    }
}

function copyCouponCode() {
    const code = document.getElementById('detailCouponCode').textContent;
    copyToClipboard(code);
}

function renderMyCouponTypeOptions() {
    const filterSelect = document.getElementById('filterMyCouponType');
    if (filterSelect) {
        filterSelect.innerHTML = '<option value="">全部</option>' +
            couponTypeList.map(t => `<option value="${t.type}">${t.name}</option>`).join('');
    }
}

// ========== 申领卡券 ==========
function showApplyCouponModal() {
    // 渲染卡券类型选项
    const select = document.getElementById('applyCouponType');
    select.innerHTML = couponTypeList.map(t =>
        `<option value="${t.type}">${t.name}</option>`
    ).join('');

    // 加载第一个类型的库存
    if (couponTypeList.length > 0) {
        loadCouponStock(couponTypeList[0].type);
    }

    document.getElementById('applyCouponModal').classList.add('show');
}

function closeApplyCouponModal() {
    document.getElementById('applyCouponModal').classList.remove('show');
}

async function onApplyCouponTypeChange() {
    const type = parseInt(document.getElementById('applyCouponType').value);
    await loadCouponStock(type);
}

async function loadCouponStock(type) {
    const stockValue = document.getElementById('stockValue');
    stockValue.textContent = '加载中...';

    try {
        const data = await request(`/api/v1/my-coupon/stock?type=${type}`);
        if (data.code === 0) {
            const stock = data.data.stock;
            stockValue.textContent = stock > 0 ? stock : '0 (无库存)';
            stockValue.className = 'stock-value' + (stock > 0 ? ' has-stock' : ' no-stock');
        } else {
            stockValue.textContent = '查询失败';
        }
    } catch (err) {
        stockValue.textContent = '查询失败';
    }
}

async function confirmApplyCoupon() {
    const type = parseInt(document.getElementById('applyCouponType').value);

    const btn = document.getElementById('btnConfirmApply');
    setBtnLoading(btn, true);

    try {
        const data = await request('/api/v1/my-coupon/take', {
            method: 'POST',
            body: JSON.stringify({ type })
        });

        if (data.code === 0) {
            closeApplyCouponModal();

            // 显示成功弹窗
            showResultModal({
                success: true,
                title: '领取成功',
                message: '卡券已成功领取，点击确定查看详情',
                onClose: () => {
                    // 关闭弹窗后显示卡券详情
                    showTakenCouponResult(data.data);
                }
            });

            // 刷新我的卡券列表
            await loadMyCoupons();
        } else {
            // 显示失败弹窗
            showResultModal({
                success: false,
                title: '领取失败',
                message: data.msg || '申领失败，请稍后重试'
            });
        }
    } finally {
        setBtnLoading(btn, false);
    }
}

function showTakenCouponResult(couponData) {
    // 复用详情弹窗显示领取结果
    document.getElementById('detailCouponType').textContent = couponData.type_name;
    document.getElementById('detailCouponCode').textContent = couponData.coupon;
    document.getElementById('detailTakenAt').textContent = '刚刚领取';
    document.getElementById('myCouponDetailModal').classList.add('show');
}

// 初始化
async function init() {
    // 根据权限显示菜单
    // 管理员 -> 用户管理
    if (isAdmin) {
        document.getElementById('menuUsers').style.display = 'block';
    }
    // 库存权限 -> 卡券管理
    if (canStock) {
        document.getElementById('menuCoupons').style.display = 'block';
    }
    // 所有登录用户 -> 我的卡券
    document.getElementById('menuMyCoupons').style.display = 'block';

    // 申领卡券权限 -> 显示申领按钮
    if (canApplyCoupon) {
        document.getElementById('btnApplyCoupon').style.display = 'block';
    }

    // 确定默认页面
    if (isAdmin) {
        await loadRoles();
        await loadUsers();
    } else if (canStock) {
        // 有库存权限，默认显示卡券管理页面
        switchPage('coupons');
        await loadCouponTypes();
        await loadCoupons();
    } else {
        // 普通用户默认显示我的卡券页面
        switchPage('my-coupons');
        await loadCouponTypes();
        renderMyCouponTypeOptions();
        await loadMyCoupons();
    }
}

init();