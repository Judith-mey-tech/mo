// Адрес сервера
const API_URL = 'http://localhost:8080';

// Получаем элементы
const modal = document.getElementById('auth-modal');
const openBtn = document.getElementById('open-auth');
const closeBtn = document.getElementById('close-modal');
const form = document.getElementById('auth-form');
const submitBtn = document.getElementById('auth-submit');
const toggleBtn = document.getElementById('toggle-auth');
const title = document.getElementById('modal-title');
const msg = document.getElementById('auth-msg');

// Режим: true — вход, false — регистрация
let isLogin = true;

// === ФУНКЦИЯ ОТКРЫТИЯ МОДАЛЬНОГО ОКНА ===
function openAuthModal() {
    modal.style.display = 'flex'; // flex для центрирования
    msg.textContent = '';         // очищаем сообщения
    // Сбрасываем форму
    form.reset();
    // Убедимся, что email скрыт (если был в режиме регистрации)
    const emailField = document.getElementById('email');
    if (emailField) emailField.remove();
    // Возвращаем режим входа по умолчанию
    isLogin = true;
    updateModalUI();
}

// === ФУНКЦИЯ ЗАКРЫТИЯ ===
function closeAuthModal() {
    modal.style.display = 'none';
}

// === ОБНОВЛЕНИЕ UI В ЗАВИСИМОСТИ ОТ РЕЖИМА ===
function updateModalUI() {
    title.textContent = isLogin ? 'Войти' : 'Регистрация';
    submitBtn.textContent = isLogin ? 'Войти' : 'Зарегистрироваться';
    toggleBtn.textContent = isLogin
        ? 'Переключить на регистрацию'
        : 'Переключиться на вход';

    const emailField = document.getElementById('email');
    if (isLogin && emailField) {
        emailField.remove();
    } else if (!isLogin && !emailField) {
        const emailInput = document.createElement('input');
        emailInput.type = 'email';
        emailInput.id = 'email';
        emailInput.placeholder = 'Email';
        emailInput.required = true;
        emailInput.style.cssText = `
            width: 100%;
            padding: 0.75rem;
            margin: 0.75rem 0;
            border: 1px solid #ddd;
            border-radius: 6px;
            font-size: 1rem;
            box-sizing: border-box;
        `;
        document.getElementById('password').after(emailInput);
    }
}

// === ОБРАБОТЧИКИ СОБЫТИЙ ===

// 1. Открытие модалки по кнопке "Войти"
openBtn.addEventListener('click', openAuthModal);

// 2. Закрытие по крестику
closeBtn.addEventListener('click', closeAuthModal);

// 3. Закрытие при клике вне модалки
window.addEventListener('click', (e) => {
    if (e.target === modal) {
        closeAuthModal();
    }
});

// 4. Переключение режимов
toggleBtn.addEventListener('click', () => {
    isLogin = !isLogin;
    updateModalUI();
    msg.textContent = '';
});

// 5. Отправка формы
form.addEventListener('submit', async (e) => {
    e.preventDefault();

    const username = document.getElementById('username').value.trim();
    const password = document.getElementById('password').value;
    const email = isLogin ? null : document.getElementById('email')?.value.trim();

    if (!username || !password || (!isLogin && !email)) {
        msg.textContent = 'Заполните все поля!';
        msg.style.color = '#ff6b6b';
        return;
    }

    const endpoint = isLogin ? '/login' : '/register';
    const body = isLogin ? { username, password } : { username, password, email };

    try {
        const res = await fetch(API_URL + endpoint, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body)
        });
        const data = await res.json();

        if (data.token) {
            localStorage.setItem('token', data.token);
            closeAuthModal();
            document.getElementById('upload-section').style.display = 'block';
            msg.textContent = 'Успешно!';
            msg.style.color = '#6ee7b7';
        } else {
            msg.textContent = data.error || 'Ошибка сервера';
            msg.style.color = '#ff6b6b';
        }
    } catch (err) {
        msg.textContent = 'Нет связи с сервером';
        msg.style.color = '#ff6b6b';
        console.error('Ошибка:', err);
    }
});

// === Загрузка файла ===
document.getElementById('upload-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const token = localStorage.getItem('token');
    if (!token) {
        alert('Сначала войдите в систему!');
        return;
    }

    const file = document.getElementById('file-input').files[0];
    const desc = document.getElementById('desc').value;
    if (!file) {
        document.getElementById('upload-result').textContent = 'Выберите файл!';
        return;
    }

    const formData = new FormData();
    formData.append('file', file);
    if (desc) formData.append('description', desc);

    const resultEl = document.getElementById('upload-result');
    resultEl.textContent = 'Загрузка...';

    try {
        const res = await fetch(API_URL + '/upload', {
            method: 'POST',
            headers: { 'Authorization': 'Bearer ' + token },
            body: formData
        });
        const data = await res.json();

        if (data.file_id) {
            resultEl.innerHTML = `
                Файл загружен! ID: <strong>${data.file_id}</strong><br>
                Проверка подлинности...
            `;
            checkFile(data.file_id, resultEl);
        } else {
            resultEl.textContent = data.error || 'Ошибка загрузки';
            resultEl.style.color = '#ff6b6b';
        }
    } catch (err) {
        resultEl.textContent = 'Ошибка сети';
        resultEl.style.color = '#ff6b6b';
    }
});

// === Проверка на дипфейк ===
async function checkFile(id, el) {
    const token = localStorage.getItem('token');
    try {
        const res = await fetch(API_URL + `/check/${id}`, {
            headers: { 'Authorization': 'Bearer ' + token }
        });
        const data = await res.json();

        if (data.original) {
            el.innerHTML += `<br><strong>Оригинал подтверждён AuthenTrack</strong>`;
            el.style.color = '#6ee7b7';
        } else {
            el.innerHTML += `<br><strong>Обнаружен дипфейк!</strong>`;
            el.style.color = '#ff6b6b';
        }
    } catch (err) {
        el.innerHTML += `<br>Ошибка проверки`;
        el.style.color = '#ff6b6b';
    }
}