function filterLocations() {
    const query = document.getElementById('location-search').value.toLowerCase();
    const options = document.querySelectorAll('.location-option');
    options.forEach(option => {
        const text = option.textContent.toLowerCase();
        if (text.includes(query)) {
            option.style.display = 'flex';
        } else {
            option.style.display = 'none';
        }
    });
}