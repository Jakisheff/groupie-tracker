document.addEventListener('DOMContentLoaded', () => {
    const searchInput = document.getElementById('search-input');
    const suggestionsContainer = document.getElementById('suggestions-container');
    let debounceTimer;

    searchInput.addEventListener('input', (e) => {
        const query = e.target.value.trim();
        
        clearTimeout(debounceTimer);
        
        if (query.length === 0) {
            suggestionsContainer.style.display = 'none';
            suggestionsContainer.innerHTML = '';
            return;
        }

        debounceTimer = setTimeout(() => {
            fetchSuggestions(query);
        }, 300); // 300ms debounce
    });

    // Close suggestions when clicking outside
    document.addEventListener('click', (e) => {
        if (!searchInput.contains(e.target) && !suggestionsContainer.contains(e.target)) {
            suggestionsContainer.style.display = 'none';
        }
    });

    async function fetchSuggestions(query) {
        try {
            const response = await fetch(`/api/search?q=${encodeURIComponent(query)}`);
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            const suggestions = await response.json();
            renderSuggestions(suggestions);
        } catch (error) {
            console.error('Error fetching suggestions:', error);
        }
    }

    function renderSuggestions(suggestions) {
        suggestionsContainer.innerHTML = '';
        if (suggestions.length === 0) {
            suggestionsContainer.style.display = 'none';
            return;
        }

        suggestions.forEach(suggestion => {
            const div = document.createElement('div');
            div.classList.add('suggestion-item');
            
            // Format: "Queen - artist/band"
            const textSpan = document.createElement('span');
            textSpan.textContent = suggestion.text;
            textSpan.style.fontWeight = 'bold';

            const typeSpan = document.createElement('span');
            typeSpan.textContent = ` - ${suggestion.type}`;
            typeSpan.style.color = '#666';
            typeSpan.style.fontSize = '0.9em';

            div.appendChild(textSpan);
            div.appendChild(typeSpan);

            div.addEventListener('click', () => {
                window.location.href = `/artist/${suggestion.id}`;
            });

            suggestionsContainer.appendChild(div);
        });

        suggestionsContainer.style.display = 'block';
    }
});
