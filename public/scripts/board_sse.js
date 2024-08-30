function setupSSE() {
    if (typeof EventSource !== "undefined") {
        const source = new EventSource('/game/stream');

        source.onopen = function (event) {
            console.log('SSE connection opened', event);
        };

        source.onerror = function (event) {
            console.error('SSE connection error', event);
        };

        source.addEventListener('board-update', function (event) {
            console.log('Board update received');
            document.getElementById('board-container').innerHTML = event.data;
            htmx.process(document.getElementById('board-container'));
        });

        source.addEventListener('score-update', function (event) {
            console.log('Score update received');
            document.getElementById('score-container').innerHTML = event.data;
        });
    } else {
        console.error('SSE not supported');
    }
}

document.addEventListener('htmx:load', setupSSE);
setupSSE();  // Call it immediately as well