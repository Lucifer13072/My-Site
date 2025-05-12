async function loadCaptcha() {
    const res = await fetch('/captcha/task');
    const task = await res.json();
    const grid = document.getElementById('captcha-grid');
    grid.innerHTML = '';
    task.images.forEach(src => {
        const img = document.createElement('img');
        img.src = '/assets/captcha/' + src;
        img.dataset.id = src;
        img.onclick = () => img.classList.toggle('selected');
        grid.appendChild(img);
    });

    document.getElementById('captcha-submit').onclick = async () => {
        const selected = {};
        document.querySelectorAll('#captcha-grid img').forEach(img => {
            selected[img.dataset.id] = img.classList.contains('selected');
        });
        const resp = await fetch('/captcha/submit', {
            method: 'POST',
            headers: {'Content-Type':'application/json'},
            body: JSON.stringify({
                session_id: task.session_id,
                task_type: task.task_type,
                answers: selected
            })
        });
        const result = await resp.json();
        document.getElementById('captcha-result').innerText =
            result.success ? 'Прошёл проверку 👍' : 'Не прошёл. Попробуй ещё раз.';
        if (result.success) {
            // скрыть капчу или разблокировать функционал
            document.getElementById('captcha-block').style.display = 'none';
        }
    };
}

// Запускаем на загрузке страницы
window.addEventListener('DOMContentLoaded', loadCaptcha);