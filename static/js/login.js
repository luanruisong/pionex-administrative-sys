// login.js - 登录页面逻辑

document.getElementById('loginForm').addEventListener('submit', async function(e) {
    e.preventDefault();

    const btn = document.getElementById('btnLogin');
    const account = document.getElementById('account').value.trim();
    const password = document.getElementById('password').value;

    if (!account || !password) {
        toast('请输入账号和密码', 'warning');
        return;
    }

    setBtnLoading(btn, true);

    try {
        const resp = await fetch('/api/v1/user/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ account, password })
        });

        const data = await resp.json();

        if (data.code === 0) {
            toast('登录成功！', 'success');
            localStorage.setItem('token', data.data.token);
            localStorage.setItem('role', String(data.data.role || 0));
            localStorage.setItem('user_name', data.data.name || '');
            setTimeout(() => {
                window.location.href = '/static/html/main.html';
            }, 800);
        } else {
            toast(data.msg || '登录失败', 'error');
        }
    } catch (err) {
        toast('网络错误，请重试', 'error');
    } finally {
        setBtnLoading(btn, false);
    }
});
