document.addEventListener('DOMContentLoaded', function() {
    const logoutBtn = document.getElementById('logoutBtn');
    const noteForm = document.getElementById('noteForm');
    const errorMessageEl = document.getElementById('errorMessage');
    const notesContainer = document.getElementById('notesContainer');

    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = '/';
        return;
    }

    // Загружаем заметки при загрузке страницы
    loadNotes();

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
               
                loadNotes();
            } else {
                errorMessageEl.textContent = data.error || 'Failed to create note';
            }
        } catch (err) {
            errorMessageEl.textContent = 'Error connecting to server';
            console.error('Note creation error:', err);
        }
    });

    // Функция для загрузки заметок
    async function loadNotes() {
        try {
            const response = await fetch('/note', {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            if (!response.ok) {
                throw new Error('Failed to fetch notes');
            }

            const notes = await response.json();
            renderNotes(notes);
        } catch (err) {
            errorMessageEl.textContent = 'Error loading notes';
            console.error('Error loading notes:', err);
        }
    }

    // Функция для отображения заметок
    function renderNotes(notes) {
        if (!notesContainer) return;
        
        notesContainer.innerHTML = ''; 

        if (!notes || notes.length === 0) {
            notesContainer.innerHTML = '<p>No notes yet.</p>';
            return;
        }

        notes.forEach(note => {
            const noteElement = document.createElement('div');
            noteElement.className = 'note';
            noteElement.innerHTML = `
                <h3>${note.title || 'No title'}</h3>
                <p><strong>Category:</strong> ${note.category || 'No category'}</p>
                <p>${note.content || 'No content'}</p>
                <p><small>Created: ${new Date(note.created_at).toLocaleString()}</small></p>
            `;
            notesContainer.appendChild(noteElement);
        });
    }
});