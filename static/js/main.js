document.addEventListener('DOMContentLoaded', function() {
    const logoutBtn = document.getElementById('logoutBtn');
    const noteForm = document.getElementById('noteForm');
    const errorMessageEl = document.getElementById('errorMessage');

    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = '/';
        return;
    }

    logoutBtn.addEventListener('click', function() {
        localStorage.removeItem('token');
        window.location.href = '/';
    });

    noteForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const category = document.getElementById('category').value;
        const title = document.getElementById('title').value;
        const content = document.getElementById('content').value;

        try {
            const response = await fetch('/note', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    category: category,
                    title: title,
                    content: content
                })
            });

            const data = await response.json();

            if (response.ok) {
                alert('Note created successfully!');
                noteForm.reset();
            } else {
                errorMessageEl.textContent = data.error || 'Failed to create note';
            }
        } catch (err) {
            errorMessageEl.textContent = 'Error connecting to server';
            console.error('Note creation error:', err);
        }
    });
});