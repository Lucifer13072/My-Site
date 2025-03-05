const matrix = document.querySelector('.matrix');
const word_list = [" ERROR ", " В РАЗРАБОТКЕ "];

function createMatrix() {
    // Очищаем предыдущие колонки
    matrix.innerHTML = '';
    const numColumns = Math.floor(window.innerWidth / 20); // Количество колонок

    for (let i = 0; i < numColumns; i++) {
        const column = document.createElement('div');
        column.classList.add('column');
        column.style.left = `${i * 20}px`; // Расстояние между колонками

        // Генерируем случайные слова из списка
        setInterval(() => {

            const randomWord = word_list[Math.floor(Math.random() * word_list.length)];
            column.textContent += randomWord; // Добавляем слово в столбец

            // Ограничиваем длину текста
            if (column.textContent.length > 20) {
                column.textContent = column.textContent.slice(randomWord.length); // Удаляем слово
            }
        }, 700); // Интервал для добавления слов

        // Анимация падения
        column.style.animationDuration = `${Math.random() * 3 + 2}s`; // Случайная скорость падения
        matrix.appendChild(column);
    }
}

// Инициализация матрицы
createMatrix();

// Обновление матрицы при изменении размера окна
window.addEventListener('resize', createMatrix);
