document.addEventListener('DOMContentLoaded', function() {
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = '/';
        return;
    }

    const pathParts = window.location.pathname.split('/');
    const noteId = pathParts[pathParts.length - 2];

    
    loadNote(noteId);

    document.getElementById('backBtn').addEventListener('click', function() {
        window.location.href = '/main?token=' + token;
    });

    document.getElementById('editNoteBtn').addEventListener('click', function() {
        toggleEditMode(true);
    });

    document.getElementById('saveNoteBtn').addEventListener('click', function() {
        updateNote(noteId);
    });

    document.getElementById('deleteNoteBtn').addEventListener('click', function() {
        if (confirm('Are you sure you want to delete this note? This action cannot be undone.')) {
            deleteNote(noteId);
        }
    });

    async function loadNote(id) {
        try {
            document.getElementById('noteContent').innerHTML = '<p>Loading note...</p>';
            
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
            document.getElementById('noteContent').innerHTML = `
                <div class="error">
                    Error loading note: ${err.message}
                </div>
            `;
            console.error(err);
        }
    }

    function displayNote(note) {
        document.getElementById('noteTitle').textContent = note.title || 'Untitled Note';
        document.getElementById('noteCategory').textContent = note.category || 'General';
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
            
            const titleInput = document.createElement('input');
            titleInput.type = 'text';
            titleInput.id = 'editTitle';
            titleInput.className = 'form-control';
            titleInput.value = title.textContent;
            title.replaceWith(titleInput);

            const categoryInput = document.createElement('input');
            categoryInput.type = 'text';
            categoryInput.id = 'editCategory';
            categoryInput.className = 'form-control';
            categoryInput.value = category.textContent;
            category.replaceWith(categoryInput);

            const contentTextarea = document.createElement('textarea');
            contentTextarea.id = 'editContent';
            contentTextarea.className = 'form-control';
            contentTextarea.rows = '10';
            contentTextarea.value = content.textContent;
            content.replaceWith(contentTextarea);

            editBtn.style.display = 'none';
            saveBtn.style.display = 'inline-block';
        } else {
            location.reload();
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

            showNotification('Note updated successfully', 'success');
            toggleEditMode(false);
            loadNote(id);
        } catch (err) {
            showNotification('Error updating note: ' + err.message, 'error');
            console.error(err);
        }
    }

    async function deleteNote(id) {
        try {
            const response = await fetch(`/notes/${id}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            if (!response.ok) {
                throw new Error('Failed to delete note');
            }

            showNotification('Note deleted successfully', 'success');
            setTimeout(() => {
                window.location.href = '/main?token=' + token;
            }, 1500);
        } catch (err) {
            showNotification('Error deleting note: ' + err.message, 'error');
            console.error(err);
        }
    }

    function showNotification(message, type = 'info') {
        const notification = document.createElement('div');
        notification.className = `notification ${type}`;
        notification.textContent = message;
        
        document.body.appendChild(notification);
        
        setTimeout(() => {
            notification.classList.add('fade-out');
            setTimeout(() => notification.remove(), 500);
        }, 3000);
    }
});