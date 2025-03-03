const lines = [
    'Инициализация системы...',
    'Подключение к базе данных...',
    'Загрузка библиотек...',
    'Установка соединения...',
    'Запуск интерфейса...',
    'Готово!\n',
];

let i = 0;
const consoleElem = document.getElementById('console');
const loaderContainer = document.getElementById('loaderContainer');
const progressCircle = document.querySelector('.progress-ring circle');

function typeLine() {
    if (i < lines.length) {
        consoleElem.innerHTML += lines[i] + "\n";
        i++;
        setTimeout(typeLine, 1000);
    } else {
        setTimeout(() => {
            consoleElem.style.display = 'none';
            loaderContainer.style.display = 'block';

            // Задержка перед запуском анимации
            setTimeout(() => {
                // Явный запуск анимации
                progressCircle.style.transition = 'stroke-dashoffset 1s linear';
                progressCircle.style.strokeDashoffset = "0"; // Заполняем круг
            }, 100);  // Маленькая задержка перед запуском анимации
        }, 500);
    }
}

progressCircle.addEventListener('transitionend', () => {
    // Скрыть анимацию и показать скрытый контент
    loaderContainer.style.display = 'none';
    hiddenContent.style.display = 'block';
});

typeLine();