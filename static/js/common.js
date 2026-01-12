// common.js - 共享 JavaScript 函数

// Toast 提示
function toast(message, type = 'info', duration = 3000) {
    const container = document.getElementById('toastContainer');
    const t = document.createElement('div');
    t.className = `toast ${type}`;
    t.innerHTML = `<span>${message}</span>`;
    container.appendChild(t);

    setTimeout(() => {
        t.classList.add('hide');
        setTimeout(() => t.remove(), 300);
    }, duration);
}

// 按钮加载状态
function setBtnLoading(btn, loading) {
    if (loading) {
        btn.classList.add('btn-loading');
        btn.disabled = true;
    } else {
        btn.classList.remove('btn-loading');
        btn.disabled = false;
    }
}

// Loading 遮罩控制
function showLoading() {
    const overlay = document.getElementById('loadingOverlay');
    if (overlay) {
        overlay.classList.add('show');
    }
}

function hideLoading() {
    const overlay = document.getElementById('loadingOverlay');
    if (overlay) {
        overlay.classList.remove('show');
    }
}

// 格式化13位时间戳为 yyyy-MM-dd HH:mm:ss
function formatTimestamp(timestamp) {
    if (!timestamp) return '-';
    try {
        const d = new Date(timestamp);
        const year = d.getFullYear();
        const month = String(d.getMonth() + 1).padStart(2, '0');
        const day = String(d.getDate()).padStart(2, '0');
        const hours = String(d.getHours()).padStart(2, '0');
        const minutes = String(d.getMinutes()).padStart(2, '0');
        const seconds = String(d.getSeconds()).padStart(2, '0');
        return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    } catch {
        return '-';
    }
}

// 格式化日期
function formatDate(isoStr) {
    if (!isoStr) return '-';
    try {
        const d = new Date(isoStr);
        return d.toLocaleDateString('zh-CN');
    } catch {
        return isoStr;
    }
}

// 格式化交易量/流动性（大数字显示为 K, M 等）
function formatVolume(vol) {
    if (!vol || vol === '-') return '-';
    const num = parseFloat(vol);
    if (isNaN(num)) return vol;
    if (num >= 1000000) return (num / 1000000).toFixed(2) + 'M';
    if (num >= 1000) return (num / 1000).toFixed(2) + 'K';
    return num.toFixed(2);
}

// 格式化价格（显示为百分比）
function formatPrice(price) {
    if (!price) return '-';
    const num = parseFloat(price);
    if (isNaN(num)) return price;
    return (num * 100).toFixed(1) + '%';
}

// 结果弹窗（带图片）
function showResultModal(options) {
    const { success, title, message, onClose } = options;

    // 创建弹窗元素
    let overlay = document.getElementById('resultModalOverlay');
    if (!overlay) {
        overlay = document.createElement('div');
        overlay.id = 'resultModalOverlay';
        overlay.className = 'result-modal-overlay';
        overlay.innerHTML = `
            <div class="result-modal">
                <div id="resultModalTitle" class="result-modal-title"></div>
                <img id="resultModalImage" class="result-modal-image" src="" alt="">
                <div id="resultModalMessage" class="result-modal-message"></div>
                <button id="resultModalBtn" class="result-modal-btn">确定</button>
            </div>
        `;
        document.body.appendChild(overlay);
    }

    // 设置内容
    const titleEl = document.getElementById('resultModalTitle');
    const imageEl = document.getElementById('resultModalImage');
    const messageEl = document.getElementById('resultModalMessage');
    const btnEl = document.getElementById('resultModalBtn');

    titleEl.textContent = title || (success ? '领取成功' : '领取失败');
    titleEl.className = 'result-modal-title ' + (success ? 'success' : 'error');

    // 设置图片
    imageEl.src = success ? '/static/image/ganbei.png' : '/static/image/iam_dead.png';

    messageEl.textContent = message || '';

    // 显示弹窗
    overlay.classList.add('show');

    // 绑定关闭事件
    const closeModal = () => {
        overlay.classList.remove('show');
        if (onClose) onClose();
    };

    btnEl.onclick = closeModal;
    overlay.onclick = (e) => {
        if (e.target === overlay) closeModal();
    };
}
