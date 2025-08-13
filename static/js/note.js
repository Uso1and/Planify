document.addEventListener('DOMContentLoaded', function() {
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = '/';
        return;
    }

    // Получаем ID заметки из URL
    const pathParts = window.location.pathname.split('/');
    const noteId = pathParts[pathParts.length - 2]; // /note/123/view -> 123

    // Загружаем данные заметки
    loadNote(noteId);

    // Обработчики кнопок
    document.getElementById('backBtn').addEventListener('click', function() {
        window.location.href = '/main?token=' + token;
    });

    document.getElementById('editNoteBtn').addEventListener('click', function() {
        window.location.href = `/note/${noteId}/edit?token=${token}`;
    });

    document.getElementById('deleteNoteBtn').addEventListener('click', function() {
        if (confirm('Are you sure you want to delete this note?')) {
            deleteNote(noteId);
        }
    });

    async function loadNote(id) {
        try {
            const response = await fetch(`/note/${id}`, {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            if (!response.ok) {
                throw new Error('Failed to fetch note');
            }

            const note = await response.json();
            displayNote(note);
        } catch (err) {
            alert('Error loading note: ' + err.message);
            console.error(err);
        }
    }

    function displayNote(note) {
        document.getElementById('noteTitle').textContent = note.title || 'No title';
        document.getElementById('noteCategory').textContent = note.category || 'No category';
        document.getElementById('noteContent').textContent = note.content || 'No content';
        document.getElementById('noteCreatedAt').textContent = new Date(note.created_at).toLocaleString();
        document.getElementById('noteUpdatedAt').textContent = new Date(note.updated_at).toLocaleString();
    }

    async function deleteNote(id) {
        try {
            const response = await fetch(`/note/${id}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            if (!response.ok) {
                throw new Error('Failed to delete note');
            }

            alert('Note deleted successfully');
            window.location.href = '/main?token=' + token;
        } catch (err) {
            alert('Error deleting note: ' + err.message);
            console.error(err);
        }
    }
});