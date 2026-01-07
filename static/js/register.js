// register.js - 注册页面逻辑

document.getElementById('registerForm').addEventListener('submit', async function(e) {
    e.preventDefault();

    const btn = document.getElementById('btnRegister');
    const name = document.getElementById('name').value.trim();
    const account = document.getElementById('account').value.trim();
    const password = document.getElementById('password').value;

    if (!name || !account || !password) {
        toast('请填写完整信息', 'warning');
        return;
    }

    setBtnLoading(btn, true);

    try {
        const resp = await fetch('/api/v1/user/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, account, password })
        });

        const data = await resp.json();

        if (data.code === 0) {
            toast('注册成功！即将跳转登录...', 'success');
            setTimeout(() => {
                window.location.href = '/static/html/login.html';
            }, 1000);
        } else {
            toast(data.msg || '注册失败', 'error');
        }
    } catch (err) {
        toast('网络错误，请重试', 'error');
    } finally {
        setBtnLoading(btn, false);
    }
});
