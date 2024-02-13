function setExpression() {
    const expression = document.getElementById('expression').value;

    // Отправка выражения на сервер и получение айди
    fetch('http://localhost:8080/ex/set', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ expression }),
    })
        .then(response => response.json())
        .then(data => {
            // Отображение полученного айди
            document.getElementById('result').innerText = `Generated ID: ${data}`;
        })
        .catch(error => {
            console.error('Error:', error);
            document.getElementById('result').innerText = 'Error occurred while setting expression';
        });
}
