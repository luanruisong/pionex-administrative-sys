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
