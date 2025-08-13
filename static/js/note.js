document.addEventListener('DOMContentLoaded', function() {
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = '/';
        return;
    }

    const pathParts = window.location.pathname.split('/');
    const noteId = pathParts[pathParts.length - 2];

    // Загружаем данные заметки
    loadNote(noteId);

    document.getElementById('backBtn').addEventListener('click', function() {
        window.location.href = '/main?token=' + token;
    });

    document.getElementById('editNoteBtn').addEventListener('click', function() {
        // Переключаем режим просмотра/редактирования
        toggleEditMode(true);
    });

    document.getElementById('saveNoteBtn').addEventListener('click', function() {
        updateNote(noteId);
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

    function toggleEditMode(editMode) {
        const title = document.getElementById('noteTitle');
        const category = document.getElementById('noteCategory');
        const content = document.getElementById('noteContent');
        const editBtn = document.getElementById('editNoteBtn');
        const saveBtn = document.getElementById('saveNoteBtn');

        if (editMode) {
            // Создаем input для редактирования
            const titleInput = document.createElement('input');
            titleInput.type = 'text';
            titleInput.id = 'editTitle';
            titleInput.value = title.textContent;
            title.replaceWith(titleInput);

            const categoryInput = document.createElement('input');
            categoryInput.type = 'text';
            categoryInput.id = 'editCategory';
            categoryInput.value = category.textContent;
            category.replaceWith(categoryInput);

            const contentTextarea = document.createElement('textarea');
            contentTextarea.id = 'editContent';
            contentTextarea.value = content.textContent;
            content.replaceWith(contentTextarea);

            editBtn.style.display = 'none';
            saveBtn.style.display = 'inline-block';
        } else {
            // Возвращаем обратно в режим просмотра
            location.reload(); // или обновляем данные через API
        }
    }

    async function updateNote(id) {
        try {
            const updatedNote = {
                category: document.getElementById('editCategory').value,
                title: document.getElementById('editTitle').value,
                content: document.getElementById('editContent').value
            };

            const response = await fetch(`/notes/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(updatedNote)
            });

            if (!response.ok) {
                throw new Error('Failed to update note');
            }

            alert('Note updated successfully');
            toggleEditMode(false);
            loadNote(id); // Перезагружаем заметку
        } catch (err) {
            alert('Error updating note: ' + err.message);
            console.error(err);
        }
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