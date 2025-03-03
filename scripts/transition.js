document.addEventListener('DOMContentLoaded', function () {
    const button = document.getElementById('transitionButton');

    // Переменные для настройки скорости и размеров пикселей
    const pixelSize = 2; // Размер пикселей (ширина и высота)
    const pixelSpeed = 5 * 1000; // Время падения пикселей (в миллисекундах)
    const numberOfPixels = 5000; // Количество пикселей (увеличено в разы)
    const pixelInterval = 1; // Интервал появления пикселей (мс)
    const time = 5 * 1000;

    // Функция для анимации распада страницы
    button.addEventListener('click', (e) => {
        e.preventDefault(); // Запрещаем стандартное поведение кнопки

        // Скроем контент
        document.getElementById('console').style.display = 'none';
        document.getElementById('loaderContainer').style.display = 'none';

        // Добавляем черный элемент, который будет заполнять экран снизу вверх
        let blackCover = document.createElement('div');
        blackCover.classList.add('black-cover');
        document.body.appendChild(blackCover);

        // Функция для создания пикселей
        function createPixel() {
            const pixel = document.createElement('div');
            pixel.classList.add('pixel');

            // Устанавливаем размеры пикселя
            pixel.style.width = `${pixelSize}px`;
            pixel.style.height = `${pixelSize}px`;

            // Позиция пикселя будет в верхней части черного блока
            const topPosition = window.innerHeight - blackCover.offsetHeight + 'px'; // Позиция пикселя в верхней части черного блока
            const leftPosition = Math.random() * window.innerWidth + 'px'; // Позиция пикселя по ширине
            pixel.style.top = topPosition;
            pixel.style.left = leftPosition;

            // Добавляем пиксель на страницу
            document.body.appendChild(pixel);

            // Анимация для пикселя (с корректировкой скорости)
            pixel.style.animationDuration = `${pixelSpeed / 1000}s`; // Задаем длительность анимации для падения
        }

        // Создаем пиксели с количеством numberOfPixels
        let pixelCount = 0;
        const pixelCreationInterval = setInterval(() => {
            createPixel();
            createPixel();
            createPixel();
            createPixel();
            createPixel();
            createPixel();
            pixelCount++;

            // Когда все пиксели созданы, останавливаем интервал
            if (pixelCount >= numberOfPixels) {
                clearInterval(pixelCreationInterval);
            }
        }, pixelInterval);

        // Останавливаем создание пикселей, когда черный блок полностью закроет экран
        setTimeout(() => {
            clearInterval(pixelCreationInterval);
        }, time);

        // После завершения анимации, редирект на страницу
        setTimeout(() => {
            window.location.href = button.href;
        }, time); // Убираем пиксели и черный экран через 2 секунды после завершения анимации
    });
});
